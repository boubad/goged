package infoged

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestNewTreeFromFile(t *testing.T) {
	fixture := NewTestFixture()
	testGedcomFile := fixture.TestGedcomFilePath
	otree, err := NewTreeFromFile(testGedcomFile)
	if err != nil {
		t.Fatalf("Error parsing GEDCOM file: %v", err)
	}
	if otree == nil {
		t.Fatalf("Failed to parse GEDCOM file: %s", testGedcomFile)
	}
	if !otree.IsValid() {
		t.Fatalf("Document is nil in the parsed tree")
	}
	dataString := otree.String()
	if dataString == "" {
		t.Fatalf("Failed to parse GEDCOM file: %s", testGedcomFile)
	}
} // TestNewTreeFromFile
func TestNewFromSting(t *testing.T) {
	fixture := NewTestFixture()
	testGedcomFile := fixture.TestGedcomFilePath
	data, err := os.ReadFile(testGedcomFile)
	if err != nil {
		t.Fatalf("Error parsing GEDCOM file: %v", err)
	}
	otree, err := NewTreeFromString(string(data))
	if err != nil {
		t.Fatalf("Error parsing GEDCOM file: %v", err)
	}
	if otree == nil {
		t.Fatalf("Failed to parse GEDCOM file: %s", testGedcomFile)
	}
	if !otree.IsValid() {
		t.Fatalf("Document is nil in the parsed tree")
	}
} // TestNewFromSting
func TestDescribe(t *testing.T) {
	fixture := NewTestFixture()
	testGedcomFile := fixture.TestGedcomFilePath
	otree, err := NewTreeFromFile(testGedcomFile)
	if err != nil {
		t.Fatalf("Error parsing GEDCOM file: %v", err)
	}
	if otree == nil {
		t.Fatalf("Failed to parse GEDCOM file: %s", testGedcomFile)
	}
	if !otree.IsValid() {
		t.Fatalf("Document is nil in the parsed tree")
	}
	oMap, err := otree.Describe()
	if err != nil {
		t.Fatalf("Error describing tree: %v", err)
	}
	if len(oMap) < 1 {
		t.Fatalf("Failed to parse GEDCOM file: %s", testGedcomFile)
	}
	fmt.Printf("Tree description: %v", oMap)

} // TestDescribe
func TestWriteToFile(t *testing.T) {
	fixture := NewTestFixture()
	otree := fixture.GetTree()
	if otree == nil {
		t.Fatalf("Failed to parse GEDCOM file: %s", fixture.TestGedcomFilePath)
	}
	/*
		err := otree.WriteToFile(fixture.TestOutputFilePath)
		if err != nil {
			t.Fatalf("Error writing to file: %v", err)
		}
	*/
} // TestWriteToFile
func TestFixNodesUniqueness(t *testing.T) {
	fixture := NewTestFixture()
	otree := fixture.GetTree()
	if otree == nil {
		t.Fatalf("Failed to parse GEDCOM file: %s", fixture.TestGedcomFilePath)
	}
	notuniques, err := otree.CheckNodesUniqueness()
	if err != nil {
		t.Fatalf("Error fixing nodes uniqueness: %v", err)
	}
	if len(notuniques) > 0 {
		for _, node := range notuniques {
			fmt.Printf("Not unique node: %s\n", node)
		}
	}
} // TestFixNodesUniqueness
func TestCleanNodes(t *testing.T) {
	fixture := NewTestFixture()
	otree := fixture.GetTree()
	if otree == nil {
		t.Fatalf("Failed to parse GEDCOM file: %s", fixture.TestGedcomFilePath)
	}
	err := otree.CleanNodes()
	if err != nil {
		t.Fatalf("Error filtering changes: %v", err)
	}
	err = otree.WriteToFile(fixture.TestOutputFilePath2)
	if err != nil {
		t.Fatalf("Error writing to file: %v", err)
	}
} // TestCleanNodes
func TestFixReferences(t *testing.T) {
	fixture := NewTestFixture()
	otree := fixture.GetTree()
	if otree == nil {
		t.Fatalf("Failed to parse GEDCOM file: %s", fixture.TestGedcomFilePath)
	}
	data, err := otree.FixReferences()
	if err != nil {
		t.Fatalf("Error fixing references: %v", err)
	}
	if data == nil {
		t.Fatal("Error fixing references")
	}
	PrintParseData(&data)
} // TestFixReferences
func TestCleanAndExport(t *testing.T) {
	fixture := NewTestFixture()
	otree := fixture.GetTree()
	if otree == nil {
		t.Fatalf("Failed to parse GEDCOM file: %s", fixture.TestGedcomFilePath)
	}
	oMapInitial, err := otree.Describe()
	if err != nil {
		t.Fatalf("Error cleaning and exporting: %v", err)
	}
	fmt.Printf("Initial map: %v\n", oMapInitial)
	err = otree.CleanAndExport(fixture.TestOutputFilePath)
	if err != nil {
		t.Fatalf("Error cleaning and exporting: %v", err)
	}
	oMapFinal, err := otree.Describe()
	if err != nil {
		t.Fatalf("Error cleaning and exporting: %v", err)
	}
	fmt.Printf("Final map: %v\n", oMapFinal)
}

func TestGetAllFiles(t *testing.T) {
	fixture := NewTestFixture()
	otree := fixture.GetTree()
	if otree == nil {
		t.Fatalf("Failed to parse GEDCOM file: %s", fixture.TestGedcomFilePath)
	}
	files, err := otree.GetAllFiles()
	if err != nil {
		t.Fatalf("Error getting all files: %v", err)
	}
	if len(files) < 1 {
		t.Fatalf("Error getting all files: %v", err)
	}
} // TestGetAllFiles
func TestNonExistingFiles(t *testing.T) {
	fixture := NewTestFixture()
	otree := fixture.GetTree()
	if otree == nil {
		t.Fatalf("Failed to parse GEDCOM file: %s", fixture.TestGedcomFilePath)
	}
	files, err := otree.GetNonExistingFiles(fixture.TestMediasDirPath)
	if err != nil {
		t.Fatalf("Error getting all files: %v", err)
	}
	if len(files) > 0 {
		fmt.Printf("Non existing files: %v\n", files)
	}
}
func TestMediasDirs(t *testing.T) {
	fixture := NewTestFixture()
	otree := fixture.GetTree()
	if otree == nil {
		t.Fatalf("Failed to parse GEDCOM file: %s", fixture.TestGedcomFilePath)
	}
	files, err := otree.GetAllFiles()
	if err != nil {
		t.Fatalf("Error getting all files: %v", err)
	}
	filesDir := fixture.TestMediasDirPath
	if filesDir == "" {
		t.Error("Expected a non-empty TestMediasDirPath, got an empty string")
	}
	info, err := os.Stat(filesDir)
	if errors.Is(err, os.ErrNotExist) {
		t.Errorf("Expected the directory at TestMediasDirPath to exist, but it does not: %s", filesDir)
	}
	if info.IsDir() == false {
		t.Errorf("Expected the directory at TestMediasDirPath to be a directory, but it is not: %s", filesDir)
	}
	// Read directory entries
	entries, err := os.ReadDir(filesDir)
	if err != nil {
		t.Fatalf("Error reading directory: %v", err)
	}
	filesToDelete := []string{}
	// Iterate over directory entries
	for _, entry := range entries {
		if !entry.IsDir() {
			name := entry.Name()
			_, ok := files[name]
			if !ok {
				filesToDelete = append(filesToDelete, name)
			}
		} // false
	} // entry
	for _, name := range filesToDelete {
		err := os.Remove(filepath.Join(filesDir, name))
		if err != nil {
			t.Fatalf("Error deleting file: %v", err)
		}
	}
} // TestMediasDisr
