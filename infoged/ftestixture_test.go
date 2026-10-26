package infoged

import ( 
	"errors"
    "os"
	"testing"
)

func TestNewTestFixture(t *testing.T) {
	fixture := NewTestFixture()
	if fixture == nil {
		t.Error("Expected a non-nil TestFixture, got nil")
	}
	filePath := fixture.TestGedcomFilePath
	if filePath == "" {
		t.Error("Expected a non-empty TestGedcomFilePath, got an empty string")
	}
	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		t.Errorf("Expected the file at TestGedcomFilePath to exist, but it does not: %s", filePath)
	}
	//
	finalPath := fixture.TestFinalFilePath
	if finalPath == "" {
		t.Error("Expected a non-empty TestFinalFilePath, got an empty string")
	}
	if _, err := os.Stat(finalPath); errors.Is(err, os.ErrNotExist) {
		t.Errorf("Expected the file at TestFinalFilePath to exist, but it does not: %s", finalPath)
	}
	//
	if fixture.TestOutputFilePath == "" {
		t.Error("Expected a non-empty TestOutputFilePath, got an empty string")
	}
	mediasDir := fixture.TestMediasDirPath
	if mediasDir == "" {
		t.Error("Expected a non-empty TestMediasDirPath, got an empty string")
	}
	info, err := os.Stat(mediasDir)
	if errors.Is(err, os.ErrNotExist) {
		t.Errorf("Expected the directory at TestMediasDirPath to exist, but it does not: %s", mediasDir)
	}
	if info.IsDir() == false {
		t.Errorf("Expected the directory at TestMediasDirPath to be a directory, but it is not: %s", mediasDir)
	}
}// TestNewTestFixture
func TestMisc(t *testing.T) {
	fixture := NewTestFixture()
	if fixture == nil {
		t.Error("Expected a non-nil TestFixture, got nil")
	}
	filepath := fixture.TestGedcomFile
	if filepath == "" {
		t.Error("Expected a non-empty TestGedcomFile, got an empty string")
	}
	if _, err := os.Stat(filepath); errors.Is(err, os.ErrNotExist) {
		t.Errorf("Expected the file at TestGedcomFile to exist, but it does not: %s", filepath)
	}
}// TestMisc