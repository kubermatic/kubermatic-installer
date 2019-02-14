package yamled

import (
	"fmt"
	"io"

	yaml "gopkg.in/yaml.v2"
)

type node struct {
	root *yaml.MapSlice
}

func Load(r io.Reader) (*node, error) {
	var data yaml.MapSlice
	if err := yaml.NewDecoder(r).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode input YAML: %v", err)
	}

	return &node{
		root: &data,
	}, nil
}

func (n *node) MarshalYAML() (interface{}, error) {
	return n.root, nil
}

func (n *node) Get(path ...interface{}) (interface{}, bool) {
	result := interface{}(*n.root)

	for _, step := range path {
		stepFound := false

		if sstep, ok := step.(string); ok {
			// step is string => try descending down a map
			node, ok := result.(yaml.MapSlice)
			if !ok {
				return nil, false
			}

			for _, item := range node {
				if item.Key.(string) == sstep {
					stepFound = true
					result = item.Value
					break
				}
			}
		} else if istep, ok := step.(int); ok {
			// step is int => try getting Nth element of list
			node, ok := result.([]interface{})
			if !ok {
				return nil, false
			}

			if istep < 0 || istep >= len(node) {
				return nil, false
			}

			stepFound = true
			result = node[istep]
		}

		if !stepFound {
			return nil, false
		}
	}

	return result, true
}

func (n *node) GetString(path ...interface{}) (string, bool) {
	val, exists := n.Get(path...)
	if !exists {
		return "", exists
	}

	asserted, ok := val.(string)

	return asserted, ok
}

func (n *node) GetInt(path ...interface{}) (int, bool) {
	val, exists := n.Get(path...)
	if !exists {
		return 0, exists
	}

	asserted, ok := val.(int)

	return asserted, ok
}

func (n *node) GetBool(path ...interface{}) (bool, bool) {
	val, exists := n.Get(path...)
	if !exists {
		return false, exists
	}

	asserted, ok := val.(bool)

	return asserted, ok
}

func (n *node) GetArray(path ...interface{}) ([]interface{}, bool) {
	val, exists := n.Get(path...)
	if !exists {
		return nil, exists
	}

	asserted, ok := val.([]interface{})

	return asserted, ok
}

func (n *node) Set(path []interface{}, newValue interface{}) bool {
	// we always need a key or array position to work with
	if len(path) == 0 {
		return false
	}

	return n.setInternal(path, newValue)
}

func (n *node) setInternal(path []interface{}, newValue interface{}) bool {
	// when we have reached the root level,
	// replace our root element with the new data structure
	if len(path) == 0 {
		return n.setRoot(newValue)
	}

	leafKey := path[len(path)-1]
	parentPath := path[0 : len(path)-1]
	target := interface{}(n.root)

	// check if the parent element exists;
	// create parent if missing

	if len(parentPath) > 0 {
		var exists bool

		target, exists = n.Get(parentPath...)
		if !exists {
			if _, ok := leafKey.(int); ok {
				// this slice can be empty for now because we will extend it later
				if !n.setInternal(parentPath, []interface{}{}) {
					return false
				}
			} else if _, ok := leafKey.(string); ok {
				if !n.setInternal(parentPath, map[string]interface{}{}) {
					return false
				}
			} else {
				return false
			}

			target, _ = n.Get(parentPath)
		}
	}

	// Now we know that the parent element exists.

	if pos, ok := leafKey.(int); ok {
		// check if we are really in an array
		if array, ok := target.([]interface{}); ok {
			for i := len(array); i <= pos; i++ {
				array = append(array, nil)
			}

			array[pos] = newValue

			return n.setInternal(parentPath, array)
		}
	} else if key, ok := leafKey.(string); ok {
		// check if we are really in a map
		if m, ok := target.(map[string]interface{}); ok {
			m[key] = newValue
			return n.setInternal(parentPath, m)
		}

		if m, ok := target.(*yaml.MapSlice); ok {
			target = *m
		}

		if m, ok := target.(yaml.MapSlice); ok {
			return n.setInternal(parentPath, setValueInMapSlice(m, key, newValue))
		}
	}

	return false
}

func (n *node) setRoot(newValue interface{}) bool {
	if asserted, ok := newValue.(yaml.MapSlice); ok {
		n.root = &asserted
		return true
	}

	if asserted, ok := newValue.(*yaml.MapSlice); ok {
		n.root = asserted
		return true
	}

	if asserted, ok := newValue.(map[string]interface{}); ok {
		n.root = makeMapSlice(asserted)
		return true
	}

	// attempted to set something that's not a map
	return false
}

func (n *node) Append(path []interface{}, newValue interface{}) bool {
	// we require maps at the root level, so the path cannot be empty
	if len(path) == 0 {
		return false
	}

	node, ok := n.Get(path...)
	if !ok {
		return n.Set(path, []interface{}{newValue})
	}

	array, ok := node.([]interface{})
	if !ok {
		return false
	}

	return n.Set(path, append(array, newValue))
}

func (n *node) Remove(path []interface{}) bool {
	// nuke everything
	if len(path) == 0 {
		return n.setRoot(yaml.MapSlice{})
	}

	leafKey := path[len(path)-1]
	parentPath := path[0 : len(path)-1]

	parent, exists := n.Get(parentPath...)
	if !exists {
		return true
	}

	if pos, ok := leafKey.(int); ok {
		if array, ok := parent.([]interface{}); ok {
			return n.setInternal(parentPath, removeArrayItem(array, pos))
		}
	} else if key, ok := leafKey.(string); ok {
		// check if we are really in a map
		if m, ok := parent.(map[string]interface{}); ok {
			delete(m, key)
			return n.setInternal(parentPath, m)
		}

		if m, ok := parent.(*yaml.MapSlice); ok {
			parent = *m
		}

		if m, ok := parent.(yaml.MapSlice); ok {
			return n.setInternal(parentPath, removeKeyFromMapSlice(m, key))
		}
	}

	return false
}
