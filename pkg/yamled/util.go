package yamled

import yaml "gopkg.in/yaml.v2"

func makeMapSlice(m map[string]interface{}) *yaml.MapSlice {
	result := make(yaml.MapSlice, 0)

	for k, v := range m {
		result = append(result, yaml.MapItem{
			Key:   k,
			Value: v,
		})
	}

	return &result
}

func setValueInMapSlice(m yaml.MapSlice, key interface{}, value interface{}) yaml.MapSlice {
	for idx, item := range m {
		if item.Key == key {
			m[idx].Value = value
			return m
		}
	}

	return append(m, yaml.MapItem{
		Key:   key,
		Value: value,
	})
}

func removeArrayItem(array []interface{}, pos int) []interface{} {
	if pos >= 0 && pos < len(array) {
		array = append(array[:pos], array[pos+1:]...)
	}

	return array
}

func removeKeyFromMapSlice(m yaml.MapSlice, key interface{}) yaml.MapSlice {
	for idx, item := range m {
		if item.Key == key {
			return append(m[:idx], m[idx+1:]...)
		}
	}

	return m
}
