package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

var errNotFound = fmt.Errorf("not found")

func main() {
	var (
		input       string
		passthrough bool
		isMaster    bool
	)
	flag.StringVar(&input, "input", "", "The file to read from. Leave empty to use STDIN.")
	flag.BoolVar(&passthrough, "passthrough", false, "Do not convert. Just decode and re-encode.")
	flag.BoolVar(&isMaster, "is-master", true, "The values are for a master cluster and not an ordinary seed.")
	flag.Parse()

	logrus.SetOutput(os.Stderr)

	var source io.ReadCloser
	if len(input) > 0 {
		file, err := os.Open(input)
		if err != nil {
			panic(err)
		}
		source = file
	} else {
		source = os.Stdin
	}
	defer source.Close() // nolint: errcheck

	// We use a yaml.MapSlice because it preserves the order of the item during
	// decode-encode. This results in the output being identical to input except
	// comments and empty lines being stripped.
	var values yaml.MapSlice
	if err := yaml.NewDecoder(source).Decode(&values); err != nil {
		panic(err)
	}

	if !passthrough {
		if err := convert_27_to_28(&values, isMaster); err != nil {
			panic(err)
		}
	}

	if err := yaml.NewEncoder(os.Stdout).Encode(values); err != nil {
		panic(err)
	}
}

func getEntry(v *yaml.MapSlice, name string) *yaml.MapItem {
	for i := range *v {
		if (*v)[i].Key.(string) == name {
			return &(*v)[i]
		}
	}
	return nil
}

func getPath(v *yaml.MapSlice, entryPath []string) (*yaml.MapItem, error) {
	if len(entryPath) == 0 {
		return &yaml.MapItem{Key: nil, Value: v}, nil
	}

	// descend down the hierarchy
	for {
		e := getEntry(v, entryPath[0])
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

func removeEntry(v *yaml.MapSlice, entryPath []string) (*yaml.MapItem, error) {
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
	e := getEntry(v, entryPath[0])
	if e == nil {
		return nil, fmt.Errorf("'%q' not found", entryPath[0])
	}

	slice := e.Value.(yaml.MapSlice)
	entry, err := removeEntry(&slice, entryPath[1:])
	if err != nil {
		return nil, err
	}

	e.Value = slice
	return entry, nil
}

func modifyEntry(v *yaml.MapSlice, entryPath []string, value interface{}) error {
	if len(entryPath) == 0 {
		return fmt.Errorf("path cannot be empty")
	}

	// descend down the hierarchy
	for len(entryPath) > 1 {
		e := getEntry(v, entryPath[0])
		if e == nil {
			return fmt.Errorf("'%q' not found", entryPath[0])
		}
		slice := e.Value.(yaml.MapSlice)

		v = &slice
		entryPath = entryPath[1:]
	}

	// overwrite the leaf
	e := getEntry(v, entryPath[0])
	if e == nil {
		return fmt.Errorf("'%q' not found", entryPath[0])
	}
	e.Value = value

	return nil
}

func mergeSection(v *yaml.MapSlice, newSectionString string) error {
	var newSection yaml.MapSlice
	if err := yaml.Unmarshal([]byte(newSectionString), &newSection); err != nil {
		return fmt.Errorf("unmarshalling the section to merge: %s", err)
	}

	return mergeTree(v, newSection)
}

func mergeTree(v *yaml.MapSlice, newSection yaml.MapSlice) error {
	for _, section := range newSection {
		existingEntry := getEntry(v, section.Key.(string))

		// the section doesn't exist in current tree
		if existingEntry == nil {
			*v = append(*v, section)
			continue
		}

		// the section exists and is a MapSlice - we merge
		existingSlice, oldIsSlice := (existingEntry.Value).(yaml.MapSlice)
		newSlice, newIsSlice := (section.Value).(yaml.MapSlice)
		if oldIsSlice && newIsSlice {
			if err := mergeTree(&existingSlice, newSlice); err != nil {
				return err
			}
			existingEntry.Value = existingSlice
			continue
		}

		return fmt.Errorf("Cannot merge cleanly - data collides")
	}

	return nil
}
