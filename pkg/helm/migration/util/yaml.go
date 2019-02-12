package util

import (
	"fmt"

	yaml "gopkg.in/yaml.v2"
)

func GetEntry(v *yaml.MapSlice, name string) *yaml.MapItem {
	for i := range *v {
		if (*v)[i].Key.(string) == name {
			return &(*v)[i]
		}
	}

	return nil
}

func GetPath(v *yaml.MapSlice, entryPath []string) (*yaml.MapItem, error) {
	if len(entryPath) == 0 {
		return &yaml.MapItem{Key: nil, Value: v}, nil
	}

	// descend down the hierarchy
	for {
		e := GetEntry(v, entryPath[0])
		if e == nil {
			return nil, fmt.Errorf("'%q' not found", entryPath[0])
		}

		if len(entryPath) == 1 {
			return e, nil
		}

		slice := e.Value.(yaml.MapSlice)
		v = &slice
		entryPath = entryPath[1:]
	}
}

func RemoveEntry(v *yaml.MapSlice, entryPath []string) (*yaml.MapItem, error) {
	if len(entryPath) == 0 {
		return nil, fmt.Errorf("path cannot be empty")
	}

	if len(entryPath) == 1 {
		index := -1
		for i := range *v {
			if (*v)[i].Key.(string) == entryPath[0] {
				index = i
				break
			}
		}

		if index == -1 {
			return nil, fmt.Errorf("'%q' not found", entryPath[0])
		}

		e := (*v)[index]
		*v = append((*v)[:index], (*v)[index+1:]...)
		return &e, nil
	}

	// descend down the hierarchy
	e := GetEntry(v, entryPath[0])
	if e == nil {
		return nil, fmt.Errorf("'%q' not found", entryPath[0])
	}

	slice := e.Value.(yaml.MapSlice)
	entry, err := RemoveEntry(&slice, entryPath[1:])
	if err != nil {
		return nil, err
	}

	e.Value = slice
	return entry, nil
}

func ModifyEntry(v *yaml.MapSlice, entryPath []string, value interface{}) error {
	if len(entryPath) == 0 {
		return fmt.Errorf("path cannot be empty")
	}

	// descend down the hierarchy
	for len(entryPath) > 1 {
		e := GetEntry(v, entryPath[0])
		if e == nil {
			return fmt.Errorf("'%q' not found", entryPath[0])
		}
		slice := e.Value.(yaml.MapSlice)

		v = &slice
		entryPath = entryPath[1:]
	}

	// overwrite the leaf
	e := GetEntry(v, entryPath[0])
	if e == nil {
		return fmt.Errorf("'%q' not found", entryPath[0])
	}
	e.Value = value

	return nil
}

func MergeSection(v *yaml.MapSlice, newSectionString string) error {
	var newSection yaml.MapSlice
	if err := yaml.Unmarshal([]byte(newSectionString), &newSection); err != nil {
		return fmt.Errorf("unmarshalling the section to merge: %s", err)
	}

	return MergeTree(v, newSection)
}

func MergeTree(v *yaml.MapSlice, newSection yaml.MapSlice) error {
	for _, section := range newSection {
		existingEntry := GetEntry(v, section.Key.(string))

		// the section doesn't exist in current tree
		if existingEntry == nil {
			*v = append(*v, section)
			continue
		}

		// the section exists and is a MapSlice - we merge
		existingSlice, oldIsSlice := (existingEntry.Value).(yaml.MapSlice)
		newSlice, newIsSlice := (section.Value).(yaml.MapSlice)
		if oldIsSlice && newIsSlice {
			if err := MergeTree(&existingSlice, newSlice); err != nil {
				return err
			}
			existingEntry.Value = existingSlice
			continue
		}

		return fmt.Errorf("Cannot merge cleanly - data collides")
	}

	return nil
}
