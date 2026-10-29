package infoged

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/elliotchance/gedcom"
)

// /////////////////////////
type Tree struct {
	doc *gedcom.Document
}

// //////////////////////////////
func NewTreeFromFile(filename string) (*Tree, error) {
	document, err := gedcom.NewDocumentFromGEDCOMFile(filename)
	if err != nil {
		return nil, err
	}
	if document == nil {
		return nil, errors.New("failed to create a gedcom document from the file")
	}
	return &Tree{doc: document}, nil
} // NewTree
func NewTreeFromString(data string) (*Tree, error) {
	document, err := gedcom.NewDocumentFromString(data)
	if err != nil {
		return nil, err
	}
	if document == nil {
		return nil, errors.New("failed to create a gedcom document from the string")
	}
	return &Tree{doc: document}, nil
} // NewTreeFromString
func (t *Tree) IsValid() bool {
	if t == nil {
		return false
	}
	if t.doc == nil {
		return false
	}
	return true
} // IsValid
func (t *Tree) GetDocument() *gedcom.Document {
	if t == nil {
		return nil
	}
	return t.doc
} //ge
func (t *Tree) Describe() (oMap map[string]int, err error) {
	if t == nil {
		return nil, errors.New("tree is nil")
	}
	if t.doc == nil {
		return nil, errors.New("document is nil")
	}
	oMap = make(map[string]int)
	nodes := t.doc.Nodes()
	for _, node := range nodes {
		if node == nil {
			continue
		}
		id := strings.TrimSpace(node.Pointer())
		if len(id) < 1 {
			continue
		}
		stag := strings.TrimSpace(node.Tag().String())
		n, ok := oMap[stag]
		if !ok {
			oMap[stag] = 1
		} else {
			oMap[stag] = n + 1
		}
	} // doc
	return oMap, nil
} //
func (t *Tree) WriteToFile(filePath string) error {
	if t == nil {
		return errors.New("tree is nil")
	}
	if t.doc == nil {
		return errors.New("document is nil")
	}
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()
	data := t.doc.String()
	_, err = file.WriteString(data)
	if err != nil {
		return fmt.Errorf("failed to write to file: %v", err)
	}
	return nil
} // WriteToFile
func (t *Tree) String() string {
	if t == nil {
		return ""
	}
	if t.doc == nil {
		return "document is nil"
	}
	return t.doc.String()
}
func (t *Tree) GetNonExistingFiles(baseDir string) ([]string, error) {
	if t == nil {
		return nil, errors.New("tree is nil")
	}
	if t.doc == nil {
		return nil, errors.New("document is nil")
	}
	oMap := map[string]string{}
	nodes := t.doc.Nodes()
	for _, node := range nodes {
		findNotExistingFiles(&oMap, node, baseDir)
	} // doc
	lret := []string{}
	for _, v := range oMap {
		lret = append(lret, v)
	}
	return lret, nil
} // GetNonExistingFiles
func (t *Tree) GetAllFiles() (map[string]string, error) {
	if t == nil {
		return nil, errors.New("tree is nil")
	}
	if t.doc == nil {
		return nil, errors.New("document is nil")
	}
	oMap := map[string]string{}
	nodes := t.doc.Nodes()
	for _, node := range nodes {
		findAllFiles(&oMap, node)
	} // doc
	return oMap, nil
} // GetNonExistingFiles
func (t *Tree) CheckNodesUniqueness() (notuniques map[string][]string, err error) {
	if t == nil {
		return nil, errors.New("tree is nil")
	}
	if t.doc == nil {
		return nil, errors.New("document is nil")
	}
	notuniques = map[string][]string{}
	oMap := map[string]string{}
	nodes := t.doc.Nodes()
	for _, node := range nodes {
		id := node.Pointer()
		if len(id) > 0 {
			if _, ok := oMap[id]; ok {
				if _, ok := notuniques[id]; !ok {
					notuniques[id] = []string{}
				}
				notuniques[id] = append(notuniques[id], node.GEDCOMLine(gedcom.NoIndent))
				continue
			}
			oMap[id] = id
		} // id
	} // doc
	return notuniques, nil
} // CheckNodesUniqueness
func (t *Tree) FilterDuplicates() error {
	if t == nil {
		return errors.New("tree is nil")
	}
	if t.doc == nil {
		return errors.New("document is nil")
	}
	newList := []gedcom.Node{}
	nodes := t.doc.Nodes()
	for _, node := range nodes {
		xnode := gedcom.Filter(node, cleanDuplicatesFilter)
		if xnode != nil {
			newList = append(newList, xnode)
		}
	}
	t.doc.SetNodes(newList)
	return nil
} // FilterDuplicates

func (t *Tree) CleanNodes() error {
	if t == nil {
		return errors.New("tree is nil")
	}
	if t.doc == nil {
		return errors.New("document is nil")
	}
	newList := []gedcom.Node{}
	nodes := t.doc.Nodes()
	for _, node := range nodes {
		xnode := gedcom.Filter(node, cleanNodesFilter)
		if xnode != nil {
			newList = append(newList, xnode)
		}
	}
	t.doc.SetNodes(newList)
	return nil
} // CleanNodes
func (t *Tree) FilterOfficialTags() error {
	if t == nil {
		return errors.New("tree is nil")
	}
	if t.doc == nil {
		return errors.New("document is nil")
	}

	newList := []gedcom.Node{}
	nodes := t.doc.Nodes()
	for _, node := range nodes {
		xnode := gedcom.Filter(node, gedcom.OfficialTagFilter())
		if xnode != nil {
			newList = append(newList, xnode)
		}
	}
	t.doc.SetNodes(newList)
	return nil
} //
func (t *Tree) FixReferences() (map[string]*RootNodeStats, error) {
	if t == nil {
		return nil, errors.New("tree is nil")
	}
	if t.doc == nil {
		return nil, errors.New("document is nil")
	}
	return fixReferences(t.doc)
} // FixReferences
func (t *Tree) NormalizeIds() error {
	if t == nil {
		return errors.New("tree is nil")
	}
	if t.doc == nil {
		return errors.New("document is nil")
	}
	xdoc, err := fixDocIDs(t.doc)
	if err != nil {
		return err
	}
	if xdoc == nil {
		return errors.New("failed to fix IDs")
	}
	t.doc = xdoc
	return nil
} // NormalizeIds
func (t *Tree) PrepareExport() error {
	if t == nil {
		return errors.New("tree is nil")
	}
	if t.doc == nil {
		return errors.New("document is nil")
	}
	err := t.CleanNodes()
	if err != nil {
		return err
	}
	err = t.FilterDuplicates()
	if err != nil {
		return err
	}
	notuniques, err := t.CheckNodesUniqueness()
	if err != nil {
		return err
	}
	if len(notuniques) > 0 {
		for _, node := range notuniques {
			fmt.Printf("Not unique node: %s\n", node)
		}
	}
	//
	pdata, err := fixReferences(t.doc)
	if err != nil {
		return err
	}
	if pdata == nil {
		return errors.New("failed to fix references")
	}
	//
	xdoc, err := fixDocIDs(t.doc)
	if err != nil {
		return err
	}
	if xdoc == nil {
		return errors.New("failed to fix document IDs")
	}
	t.doc = xdoc
	return nil
} // PrepareExport
func (t *Tree) CleanAndExport(filePath string) error {
	if t == nil {
		return errors.New("tree is nil")
	}
	err := t.PrepareExport()
	if err != nil {
		return err
	}
	return t.WriteToFile(filePath)
} // CleanAndExport
