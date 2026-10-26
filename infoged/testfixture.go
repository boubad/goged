package infoged

import (
	"runtime"
)

type TestFixture struct {
	TestGedcomFile  string
	TestGedcomFilePath string
	TestOutputFilePath string
	TestFinalFilePath  string
	TestMediasDirPath  string
	TestOutputFilePath2 string
}

func NewTestFixture() *TestFixture {
	//
	testgedcomfile := "../testdata/tree1.ged"
	testoutputfile := "../testdata/test_go.ged"
	testfinalfile := "../testdata/test_final.ged"
	testoutputfilepath2 := "../testdata/test_go_2.ged"
	//
	currentOS := runtime.GOOS
	var testMediasDirPath string
	if currentOS == "windows" {
		testMediasDirPath = "C:\\testdata\\medias\\genea"
	} else {
		testMediasDirPath = "/home/boubad/testdata/medias/genea"
	}
	return &TestFixture{
		TestGedcomFile:     testgedcomfile,
		TestGedcomFilePath:  testgedcomfile,
		TestOutputFilePath: testoutputfile,
		TestFinalFilePath:  testfinalfile,
		TestMediasDirPath:  testMediasDirPath,
		TestOutputFilePath2: testoutputfilepath2,
	}
} // NewTestFixture
func (t *TestFixture) GetGedcomFilePath() string {
	return t.TestGedcomFilePath
}
func (t *TestFixture) GetOutputFilePath() string {
	return t.TestOutputFilePath
}
func (t *TestFixture) GetTree() *Tree {
	testGedcomFile := t.TestGedcomFilePath
	doc, err := NewTreeFromFile(testGedcomFile)
	if err != nil {
		return nil
	}
	return doc
}
