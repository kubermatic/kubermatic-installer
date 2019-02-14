package yamled

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/pmezard/go-difflib/difflib"
	yaml "gopkg.in/yaml.v2"
)

func getTestcaseYAML(t *testing.T, filename string) string {
	if filename == "" {
		filename = "document.yaml"
	}

	content, err := ioutil.ReadFile("testcases/" + filename)
	if err != nil {
		t.Errorf("could not load document %s: %v", filename, err)
	}

	return strings.TrimSpace(string(content))
}

func loadTestcase(t *testing.T, name string) (*node, string) {
	content := getTestcaseYAML(t, name)
	parts := strings.Split(content, "###")

	doc, err := Load(strings.NewReader(strings.TrimSpace(parts[0])))
	if err != nil {
		t.Errorf("failed to parse test document %s: %v", name, err)
	}

	output := ""
	if len(parts) >= 2 {
		output = strings.TrimSpace(parts[1])
	}

	return doc, output
}

func assertEqualYAML(t *testing.T, actual *node, expected string) {
	out, _ := yaml.Marshal(actual)

	diff := difflib.UnifiedDiff{
		A:        difflib.SplitLines(expected),
		B:        difflib.SplitLines(strings.TrimSpace(string(out))),
		FromFile: "Expected",
		ToFile:   "Actual",
		Context:  3,
	}
	diffStr, err := difflib.GetUnifiedDiffString(diff)
	if err != nil {
		t.Fatal(err)
	}

	if diffStr != "" {
		t.Errorf("got diff between expected and actual result: \n%s\n", diffStr)
	}
}

func TestGetRootStringKey(t *testing.T) {
	doc, _ := loadTestcase(t, "")
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
	doc, _ := loadTestcase(t, "")
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
	doc, _ := loadTestcase(t, "")
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
	doc, _ := loadTestcase(t, "")

	val, ok := doc.Get("rootNullKey")
	if !ok {
		t.Error("should have been able to retrieve root-level null, but failed")
	}

	if val != nil {
		t.Errorf("found %#v, but expected nil", val)
	}
}

func TestGetSubStringKey(t *testing.T) {
	doc, _ := loadTestcase(t, "")
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
	doc, _ := loadTestcase(t, "")
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
	doc, _ := loadTestcase(t, "")

	_, ok := doc.GetString("rootArrayKey", -1)
	if ok {
		t.Error("should NOT have been able to retrieve array item -1")
	}
}

func TestGetArrayItemOutOfRange2(t *testing.T) {
	doc, _ := loadTestcase(t, "")

	_, ok := doc.GetString("rootArrayKey", 3)
	if ok {
		t.Error("should NOT have been able to retrieve array item 3")
	}
}

func TestGetSubArrayItemExists(t *testing.T) {
	doc, _ := loadTestcase(t, "")
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
	doc, expected := loadTestcase(t, "set-existing-root-key.yaml")

	ok := doc.Set([]interface{}{"rootNullKey"}, "new value")
	if !ok {
		t.Error("should have been able to set a new root level key")
	}

	assertEqualYAML(t, doc, expected)
}

func TestSetNewRootKey(t *testing.T) {
	doc, expected := loadTestcase(t, "set-new-root-key.yaml")

	ok := doc.Set([]interface{}{"newKey"}, "new value")
	if !ok {
		t.Error("should have been able to set a new root level key")
	}

	assertEqualYAML(t, doc, expected)
}

func TestSetNewSubKey(t *testing.T) {
	doc, expected := loadTestcase(t, "set-new-sub-key.yaml")

	ok := doc.Set([]interface{}{"root", "newKey"}, "foo")
	if !ok {
		t.Error("should have been able to set a new sub level key")
	}

	assertEqualYAML(t, doc, expected)
}

func TestSetExistingArrayItem(t *testing.T) {
	doc, expected := loadTestcase(t, "set-existing-array-item.yaml")

	ok := doc.Set([]interface{}{"items", 1}, "new b")
	if !ok {
		t.Error("should have been able to overwrite array item")
	}

	assertEqualYAML(t, doc, expected)
}

func TestSetNewArrayItem(t *testing.T) {
	doc, expected := loadTestcase(t, "set-new-array-item.yaml")

	ok := doc.Set([]interface{}{"items", 3}, "d")
	if !ok {
		t.Error("should have been able to add new array item")
	}

	assertEqualYAML(t, doc, expected)
}

func TestSetNewFarArrayItem(t *testing.T) {
	doc, expected := loadTestcase(t, "set-new-far-array-item.yaml")

	ok := doc.Set([]interface{}{"items", 6}, "d")
	if !ok {
		t.Error("should have been able to add new array item")
	}

	assertEqualYAML(t, doc, expected)
}

func TestAppendToExistingArray(t *testing.T) {
	doc, expected := loadTestcase(t, "append-existing-array.yaml")

	ok := doc.Append([]interface{}{"items"}, "d")
	if !ok {
		t.Error("should have been able to append array item")
	}

	assertEqualYAML(t, doc, expected)
}

func TestAppendToNewArray(t *testing.T) {
	doc, expected := loadTestcase(t, "append-new-array.yaml")

	ok := doc.Append([]interface{}{"newItems"}, "d")
	if !ok {
		t.Error("should have been able to append array item")
	}

	assertEqualYAML(t, doc, expected)
}

func TestRemoveNonexistingRootKey(t *testing.T) {
	doc, expected := loadTestcase(t, "remove-nonexisting-root-key.yaml")

	ok := doc.Remove([]interface{}{"idontexist"})
	if !ok {
		t.Error("removing a non-existing key should be a no-op")
	}

	assertEqualYAML(t, doc, expected)
}

func TestRemoveExistingRootKey(t *testing.T) {
	doc, expected := loadTestcase(t, "remove-existing-root-key.yaml")

	ok := doc.Remove([]interface{}{"foo"})
	if !ok {
		t.Error("removing an existing key should be a successful")
	}

	assertEqualYAML(t, doc, expected)
}

func TestRemoveExistingSubKey(t *testing.T) {
	doc, expected := loadTestcase(t, "remove-existing-sub-key.yaml")

	ok := doc.Remove([]interface{}{"xyz", "foo"})
	if !ok {
		t.Error("removing an existing key should be a successful")
	}

	assertEqualYAML(t, doc, expected)
}

func TestRemoveNonexistingArrayElement(t *testing.T) {
	doc, expected := loadTestcase(t, "remove-nonexisting-array-item.yaml")

	ok := doc.Remove([]interface{}{"items", 5})
	if !ok {
		t.Error("removing a non-existing key should be a no-op")
	}

	assertEqualYAML(t, doc, expected)
}

func TestRemoveExistingArrayElement(t *testing.T) {
	doc, expected := loadTestcase(t, "remove-existing-array-item.yaml")

	ok := doc.Remove([]interface{}{"items", 1})
	if !ok {
		t.Error("removing an existing key should be a no-op")
	}

	assertEqualYAML(t, doc, expected)
}

func TestMarshalling(t *testing.T) {
	doc, expected := loadTestcase(t, "")

	assertEqualYAML(t, doc, expected)
}
