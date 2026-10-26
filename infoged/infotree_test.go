package infoged

import "testing"

func TestNewInfoTreeFromFile(t *testing.T) {
	fixture := NewTestFixture()
	testGedcomFile := fixture.TestGedcomFilePath
	otree, err := NewInfoTreeFromFile(testGedcomFile)
	if err != nil {
		t.Fatalf("Error parsing GEDCOM file: %v", err)
	}
	if otree == nil {
		t.Fatalf("Failed to parse GEDCOM file: %s", testGedcomFile)
	}
	if !otree.IsValid() {
		t.Fatalf("Document is nil in the parsed tree")
	}
} // TestNewInfoTreeFromFile
func TestNewInfoTreeFromString(t *testing.T) {
	fixture := NewTestFixture()
	testGedcomFile := fixture.TestGedcomFilePath
	init_otree, err := NewTreeFromFile(testGedcomFile)
	if init_otree == nil || err != nil {
		t.Fatalf("Error parsing GEDCOM file: %v", err)
	}
	if !init_otree.IsValid() {
		t.Fatalf("Document is nil in the parsed tree")
	}
	data := init_otree.String()
	otree, err := NewInfoTreeFromString(data)
	if err != nil {
		t.Fatalf("Error parsing GEDCOM file: %v", err)
	}
	if otree == nil {
		t.Fatalf("Failed to parse GEDCOM file: %s", testGedcomFile)
	}
	if !otree.IsValid() {
		t.Fatalf("Document is nil in the parsed tree")
	}
} // TestNewInfoTreeFromString
func TestString(t *testing.T) {
	fixture := NewTestFixture()
	testGedcomFile := fixture.TestGedcomFilePath
	otree, err := NewInfoTreeFromFile(testGedcomFile)
	if err != nil {
		t.Fatalf("Error parsing GEDCOM file: %v", err)
	}
	if otree == nil {
		t.Fatalf("Failed to parse GEDCOM file: %s", testGedcomFile)
	}
	if !otree.IsValid() {
		t.Fatalf("Document is nil in the parsed tree")
	}
	data, err := otree.String()
	if data == "" || err != nil {
		t.Fatalf("Failed to parse GEDCOM file: %s", testGedcomFile)
	}
}// TestString
func  TestWriteInfoTreeToFile(t *testing.T) {
	fixture := NewTestFixture()
	testGedcomFile := fixture.TestGedcomFilePath
	otree, err := NewInfoTreeFromFile(testGedcomFile)
	if err != nil {
		t.Fatalf("Error parsing GEDCOM file: %v", err)
	}
	if otree == nil {
		t.Fatalf("Failed to parse GEDCOM file: %s", testGedcomFile)
	}
	if !otree.IsValid() {
		t.Fatalf("Document is nil in the parsed tree")
	}
	err = otree.WriteToFile(fixture.TestOutputFilePath2)
	if err != nil {
		t.Fatalf("Error writing to file: %v", err)
	}
}// TestWriteToFile
