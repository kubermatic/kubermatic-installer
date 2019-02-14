package yamled

import (
	"strings"
	"testing"

	yaml "gopkg.in/yaml.v2"
)

const testdoc = `
rootIntKey: 12
rootStringKey: a string
rootNullKey: null
rootBoolKey: true
rootArrayKey:
- first
- second
- last
rootMapKey:
  subKey: another string
  subKey2: foobar
  subKey3:
  - first
  - second: value
    key: another value
    number: 123
  - third:
    - a
    - b
    - c
    - 123
    - null
    - 99
`

func getTestDocument(t *testing.T) *node {
	input := strings.TrimSpace(testdoc)

	node, err := Load(strings.NewReader(input))
	if err != nil {
		t.Errorf("failed to parse test document: %v", err)
	}

	return node
}

func TestGetRootStringKey(t *testing.T) {
	doc := getTestDocument(t)
	expected := "a string"

	val, ok := doc.GetString("rootStringKey")
	if !ok {
		t.Error("should have been able to retrieve root-level string, but failed")
	}

	if val != expected {
		t.Errorf("found string '%s', but expected '%s'", val, expected)
	}
}

func TestGetRootIntKey(t *testing.T) {
	doc := getTestDocument(t)
	expected := 12

	val, ok := doc.GetInt("rootIntKey")
	if !ok {
		t.Error("should have been able to retrieve root-level int, but failed")
	}

	if val != expected {
		t.Errorf("found int %d, but expected %d", val, expected)
	}
}

func TestGetRootBoolKey(t *testing.T) {
	doc := getTestDocument(t)
	expected := true

	val, ok := doc.GetBool("rootBoolKey")
	if !ok {
		t.Error("should have been able to retrieve root-level bool, but failed")
	}

	if val != expected {
		t.Errorf("found int %v, but expected %v", val, expected)
	}
}

func TestGetRootNullKey(t *testing.T) {
	doc := getTestDocument(t)

	val, ok := doc.Get("rootNullKey")
	if !ok {
		t.Error("should have been able to retrieve root-level null, but failed")
	}

	if val != nil {
		t.Errorf("found %#v, but expected nil", val)
	}
}

func TestGetSubStringKey(t *testing.T) {
	doc := getTestDocument(t)
	expected := "another string"

	val, ok := doc.GetString("rootMapKey", "subKey")
	if !ok {
		t.Error("should have been able to retrieve sub-level string, but failed")
	}

	if val != expected {
		t.Errorf("found '%s', but expected '%s'", val, expected)
	}
}

func TestGetArrayItemExists(t *testing.T) {
	doc := getTestDocument(t)
	expected := "first"

	val, ok := doc.GetString("rootArrayKey", 0)
	if !ok {
		t.Error("should have been able to retrieve root-level array item, but failed")
	}

	if val != expected {
		t.Errorf("found '%s', but expected '%s'", val, expected)
	}
}

func TestGetArrayItemOutOfRange1(t *testing.T) {
	doc := getTestDocument(t)

	_, ok := doc.GetString("rootArrayKey", -1)
	if ok {
		t.Error("should NOT have been able to retrieve array item -1")
	}
}

func TestGetArrayItemOutOfRange2(t *testing.T) {
	doc := getTestDocument(t)

	_, ok := doc.GetString("rootArrayKey", 3)
	if ok {
		t.Error("should NOT have been able to retrieve array item 3")
	}
}

func TestGetSubArrayItemExists(t *testing.T) {
	doc := getTestDocument(t)
	expected := "b"

	val, ok := doc.GetString("rootMapKey", "subKey3", 2, "third", 1)
	if !ok {
		t.Error("should have been able to retrieve sub-sub-level array item, but failed")
	}

	if val != expected {
		t.Errorf("found '%s', but expected '%s'", val, expected)
	}
}

func TestSetExistingRootKey(t *testing.T) {
	doc := getTestDocument(t)
	key := "rootNullKey"
	expected := "new value"

	ok := doc.Set([]interface{}{key}, expected)
	if !ok {
		t.Error("should have been able to set a new root level key")
	}

	val, ok := doc.Get(key)
	if !ok {
		t.Error("should have been able to retrieve the new root level key")
	}

	if val != expected {
		t.Errorf("found '%s', but expected '%s'", val, expected)
	}
}

func TestSetNewRootKey(t *testing.T) {
	doc := getTestDocument(t)
	key := "newRootKey"
	expected := "new value"

	ok := doc.Set([]interface{}{key}, expected)
	if !ok {
		t.Error("should have been able to set a new root level key")
	}

	val, ok := doc.Get(key)
	if !ok {
		t.Error("should have been able to retrieve the new root level key")
	}

	if val != expected {
		t.Errorf("found '%s', but expected '%s'", val, expected)
	}
}

func TestSetNewSubKey(t *testing.T) {
	doc := getTestDocument(t)
	keys := []interface{}{"rootMapKey", "newRootKey"}
	expected := "new value"

	ok := doc.Set(keys, expected)
	if !ok {
		t.Error("should have been able to set a new sub level key")
	}

	val, ok := doc.Get(keys...)
	if !ok {
		t.Error("should have been able to retrieve the new sub level key")
	}

	if val != expected {
		t.Errorf("found '%s', but expected '%s'", val, expected)
	}
}

func TestSetExistingArrayItem(t *testing.T) {
	doc := getTestDocument(t)
	keys := []interface{}{"rootArrayKey", 0}
	expected := "new value"

	ok := doc.Set(keys, expected)
	if !ok {
		t.Error("should have been able to overwrite array item")
	}

	val, ok := doc.Get(keys...)
	if !ok {
		t.Error("should have been able to retrieve the overwritten array item")
	}

	if val != expected {
		t.Errorf("found '%s', but expected '%s'", val, expected)
	}

	val, _ = doc.Get(keys[0])
	array, ok := val.([]interface{})
	if !ok {
		t.Error("should have been able to retrieve the array")
	}

	if len(array) != 3 {
		t.Errorf("array should have exactly 3 elements, but has %d", len(array))
	}
}

func TestSetNewArrayItem(t *testing.T) {
	doc := getTestDocument(t)
	keys := []interface{}{"rootArrayKey", 3}
	expected := "new value"

	ok := doc.Set(keys, expected)
	if !ok {
		t.Error("should have been able to add new array item")
	}

	val, ok := doc.Get(keys...)
	if !ok {
		t.Error("should have been able to retrieve the added array item")
	}

	if val != expected {
		t.Errorf("found '%s', but expected '%s'", val, expected)
	}

	val, _ = doc.Get(keys[0])
	array, ok := val.([]interface{})
	if !ok {
		t.Error("should have been able to retrieve the array")
	}

	if len(array) != 4 {
		t.Errorf("array should have exactly 4 elements, but has %d", len(array))
	}
}

func TestSetNewFarArrayItem(t *testing.T) {
	doc := getTestDocument(t)
	keys := []interface{}{"rootArrayKey", 9}
	expected := "new value"

	ok := doc.Set(keys, expected)
	if !ok {
		t.Error("should have been able to add new array item")
	}

	val, ok := doc.Get(keys...)
	if !ok {
		t.Error("should have been able to retrieve the added array item")
	}

	if val != expected {
		t.Errorf("found '%s', but expected '%s'", val, expected)
	}

	val, _ = doc.Get(keys[0])
	array, ok := val.([]interface{})
	if !ok {
		t.Error("should have been able to retrieve the array")
	}

	if len(array) != 10 {
		t.Errorf("array should have exactly 10 elements, but has %d", len(array))
	}
}

func TestAppendToExistingArray(t *testing.T) {
	doc := getTestDocument(t)
	keys := []interface{}{"rootArrayKey"}
	expected := "new value"

	ok := doc.Append(keys, expected)
	if !ok {
		t.Error("should have been able to append array item")
	}

	val, ok := doc.Get(keys...)
	if !ok {
		t.Error("should have been able to retrieve the array itself")
	}

	array, ok := val.([]interface{})
	if !ok {
		t.Error("the array should still be a slice, but is not")
	}

	if len(array) != 4 {
		t.Errorf("after appending the array should have 4 items, but has %d", len(array))
	}

	if array[len(array)-1] != expected {
		t.Errorf("last array item is '%s', but expected '%s'", val, expected)
	}
}

func TestAppendToNewArray(t *testing.T) {
	doc := getTestDocument(t)
	keys := []interface{}{"myNewArray"}
	expected := "new value"

	ok := doc.Append(keys, expected)
	if !ok {
		t.Error("should have been able to append array item")
	}

	val, ok := doc.Get(keys...)
	if !ok {
		t.Error("should have been able to retrieve the array itself")
	}

	array, ok := val.([]interface{})
	if !ok {
		t.Error("the array should still be a slice, but is not")
	}

	if len(array) != 1 {
		t.Errorf("after appending the array should have 1 item, but has %d", len(array))
	}

	if array[len(array)-1] != expected {
		t.Errorf("last array item is '%s', but expected '%s'", val, expected)
	}
}

func TestRemoveNonexistingRootKey(t *testing.T) {
	doc := getTestDocument(t)
	keys := []interface{}{"idontexist"}

	ok := doc.Remove(keys)
	if !ok {
		t.Error("removing a non-existing key should be a no-op")
	}

	marshalled, err := yaml.Marshal(doc)
	if err != nil {
		t.Errorf("failed to marshal document as YAML: %v", err)
	}

	if strings.TrimSpace(testdoc) != strings.TrimSpace(string(marshalled)) {
		t.Errorf("failed to create identical output when marshalling; expected '%s' but got '%s'", testdoc, marshalled)
	}
}

func TestRemoveExistingRootKey(t *testing.T) {
	doc := getTestDocument(t)
	keys := []interface{}{"rootBoolKey"}

	ok := doc.Remove(keys)
	if !ok {
		t.Error("removing an existing key should be a successful")
	}

	_, ok = doc.Get(keys...)
	if ok {
		t.Error("should not have been able to retrieve a removed key")
	}
}

func TestRemoveExistingSubKey(t *testing.T) {
	doc := getTestDocument(t)
	keys := []interface{}{"rootMapKey", "subKey"}

	ok := doc.Remove(keys)
	if !ok {
		t.Error("removing an existing key should be a successful")
	}

	_, ok = doc.Get(keys...)
	if ok {
		t.Error("should not have been able to retrieve a removed key")
	}

	parent, ok := doc.Get(keys[0])
	if !ok {
		t.Error("should have been able to retrieve the parent element of a removed key")
	}

	m, _ := parent.(yaml.MapSlice)
	if len(m) != 2 {
		t.Errorf("map should only have two elements left, but has %d", len(m))
	}
}

func TestRemoveNonexistingArrayElement(t *testing.T) {
	doc := getTestDocument(t)
	keys := []interface{}{"rootArrayKey", 5}

	ok := doc.Remove(keys)
	if !ok {
		t.Error("removing a non-existing key should be a no-op")
	}
}

func TestRemoveExistingArrayElement(t *testing.T) {
	doc := getTestDocument(t)
	keys := []interface{}{"rootArrayKey", 0}

	ok := doc.Remove(keys)
	if !ok {
		t.Error("removing an existing key should be a no-op")
	}

	parent, ok := doc.Get(keys[0])
	if !ok {
		t.Error("should have been able to retrieve the parent element of a removed key")
	}

	array, _ := parent.([]interface{})
	if len(array) != 2 {
		t.Errorf("array should only have two elements left, but has %d", len(array))
	}

	if s, _ := array[0].(string); s != "second" {
		t.Errorf("first element of remaining array should be 'second' but is '%s'", s)
	}

	if s, _ := array[1].(string); s != "last" {
		t.Errorf("last element of remaining array should be 'last' but is '%s'", s)
	}
}

func TestMarshalling(t *testing.T) {
	doc := getTestDocument(t)

	marshalled, err := yaml.Marshal(doc)
	if err != nil {
		t.Errorf("failed to marshal document as YAML: %v", err)
	}

	if strings.TrimSpace(testdoc) != strings.TrimSpace(string(marshalled)) {
		t.Errorf("failed to create identical output when marshalling; expected '%s' but got '%s'", testdoc, marshalled)
	}
}
