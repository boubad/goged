package infoged

import (
	"bytes"
	"errors"
	"os"

	"github.com/iand/gedcom"
)

type InfoTree struct {
	Doc *gedcom.Gedcom
}

func NewInfoTreeFromFile(filename string) (*InfoTree, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	d := gedcom.NewDecoder(bytes.NewReader(data))
	doc, err := d.Decode()
	if err != nil {
		return nil, err
	}
	return &InfoTree{Doc: doc}, nil
} // NewInfoTreeFromFile
func NewInfoTreeFromString(data string) (*InfoTree, error) {
	d := gedcom.NewDecoder(bytes.NewReader([]byte(data)))
	doc, err := d.Decode()
	if err != nil {
		return nil, err
	}
	return &InfoTree{Doc: doc}, nil
} // NewInfoTreeFromString
func (t *InfoTree) IsValid() bool {
	if t.Doc == nil {
		return false
	}
	if t.Doc.Header == nil || t.Doc.Trailer == nil {
		return false
	}
	return true
} // IsValid
func (t *InfoTree) String() (string, error) {
	if t.Doc == nil {
		return "", errors.New("document is nil")
	}
	b := new(bytes.Buffer)
	e := gedcom.NewEncoder(b)
	err := e.Encode(t.Doc)
	if err != nil {
		return "", err
	}
	return b.String(), nil
} // String
func (t *InfoTree) WriteToFile(filePath string) error {
	if t.Doc == nil {
		return errors.New("document is nil")
	}
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	e := gedcom.NewEncoder(file)
	err = e.Encode(t.Doc)
	if err != nil {
		return err
	}

	return nil
} // WriteToFile
