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
	for _, node := range t.doc.Nodes() {
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
func (t *Tree) IndivsCount() int {
	if t == nil {
		return 0
	}
	if t.doc == nil {
		return 0
	}
	return len(t.doc.Individuals())
}
func (t *Tree) FamiliesCount() int {
	if t == nil {
		return 0
	}
	if t.doc == nil {
		return 0
	}
	return len(t.doc.Families())
}
func (t *Tree) NotesCount() int {
	if t == nil {
		return 0
	}
	if t.doc == nil {
		return 0
	}
	nret := 0
	for _, node := range t.doc.Nodes() {
		if node.Tag().Is(gedcom.TagNote) {
			nret++
		}
	}
	return nret
} //NotesCount
func (t *Tree) ObjectsCount() int {
	if t == nil {
		return 0
	}
	if t.doc == nil {
		return 0
	}
	nret := 0
	for _, node := range t.doc.Nodes() {
		if node.Tag().Is(gedcom.TagObject) {
			nret++
		}
	}
	return nret
} //ObjectsCount
func (t *Tree) SourcesCount() int {
	if t == nil {
		return 0
	}
	if t.doc == nil {
		return 0
	}
	nret := 0
	for _, node := range t.doc.Nodes() {
		if node.Tag().Is(gedcom.TagSource) {
			nret++
		}
	}
	return nret
} //SourcesCount
func (t *Tree) RepositoriesCount() int {
	if t == nil {
		return 0
	}
	if t.doc == nil {
		return 0
	}
	nret := 0
	for _, node := range t.doc.Nodes() {
		if node.Tag().Is(gedcom.TagRepository) {
			nret++
		}
	}
	return nret
} //RepositoriesCount
func (t *Tree) SubmissionsCount() int {
	if t == nil {
		return 0
	}
	if t.doc == nil {
		return 0
	}
	nret := 0
	for _, node := range t.doc.Nodes() {
		if node.Tag().Is(gedcom.TagSubmission) {
			nret++
		}
	}
	return nret
} //ReferencesCount
func (t *Tree) SubmittersCount() int {
	if t == nil {
		return 0
	}
	if t.doc == nil {
		return 0
	}
	nret := 0
	for _, node := range t.doc.Nodes() {
		if node.Tag().Is(gedcom.TagSubmitter) {
			nret++
		}
	}
	return nret
} //ReferencesCount
func (t *Tree) GetNonExistingFiles(baseDir string) ([]string, error) {
	if t == nil {
		return nil, errors.New("tree is nil")
	}
	if t.doc == nil {
		return nil, errors.New("document is nil")
	}
	oMap := map[string]string{}
	for _, node := range t.doc.Nodes() {
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
	for _, node := range t.doc.Nodes() {
		findAllFiles(&oMap, node)
	} // doc
	return oMap, nil
} // GetNonExistingFiles
func (t *Tree) FixNodesUniqueness() (notuniques map[string][]string, err error) {
	if t == nil {
		return nil, errors.New("tree is nil")
	}
	if t.doc == nil {
		return nil, errors.New("document is nil")
	}
	notuniques = map[string][]string{}
	oMap := map[string]string{}
	for _, node := range t.doc.Nodes() {
		id := node.Pointer()
		if len(id) > 0 {
			if _, ok := oMap[id]; ok {
				if _, ok := notuniques[id]; !ok {
					notuniques[id] = []string{}
				}
				notuniques[id] = append(notuniques[id], node.String())
				continue
			}
			oMap[id] = id
		} // id
	} // doc
	return notuniques, nil
} // FixNodesUniqueness
func (t *Tree) FilterChanges() error {
	if t == nil {
		return errors.New("tree is nil")
	}
	if t.doc == nil {
		return errors.New("document is nil")
	}
	filterFunc := func(node gedcom.Node) (gedcom.Node, bool) {
		if node == nil {
			return nil, false
		}
		tag := node.Tag()
		if !tag.IsOfficial() || tag.Is(gedcom.TagChange) {
			return nil, false
		}
		initial_value := strings.TrimSpace(node.Value())
		if tag.Is(gedcom.TagName) {
			if strings.Contains(initial_value, "Mere") ||
				strings.Contains(initial_value, "Mère") ||
				strings.Contains(initial_value, "(") ||
				strings.Contains(initial_value, "Fille") ||
				strings.Contains(initial_value, "Fils") ||
				strings.Contains(initial_value, "Mère") ||
				strings.Contains(initial_value, "Pere") ||
				strings.Contains(initial_value, "Père") {
				return nil, false
			}
			value := strings.ToUpper(initial_value)
			if strings.Contains(value, "??") ||
				strings.Contains(value, "(INC)") ||
				strings.Contains(value, "UNKNOWN") {
				return nil, false
			}
		} // name
		if tag.Is(gedcom.TagNote) && len(node.Pointer()) < 1 && len(initial_value) < 1 {
			return nil, false
		}
		if len(node.Pointer()) == 0 && len(initial_value) == 0 && len(node.Nodes()) == 0 {
			if !tag.Is(gedcom.TagTrailer) && !tag.Is(gedcom.TagHeader) {
				return nil, false
			}
		}
		if tag.Is(gedcom.TagPedigree) || tag.Is(gedcom.TagType) {
			value := strings.ToUpper(initial_value)
			if value == "PHOTO" {
				return gedcom.NewNode(tag, "photo", ""), true
			}
			if value != initial_value {
				return gedcom.NewNode(tag, value, ""), true
			}
		}
		if tag.Is(gedcom.TagFormat) {
			value := strings.ToUpper(initial_value)
			if value == "JPG" {
				return gedcom.NewNode(tag, "jpeg", ""), true
			}
		}
		if tag.Is(gedcom.TagDestination) {
			return gedcom.NewNode(tag, "ANY", ""), true
		}
		return node, true
	} // filterFunc
	newList := []gedcom.Node{}
	for _, node := range t.doc.Nodes() {
		xnode := gedcom.Filter(node, filterFunc)
		if xnode != nil {
			newList = append(newList, xnode)
		}
	}
	t.doc.SetNodes(newList)
	return nil
} // FilterChanges
func (t *Tree) FilterOfficialTags() error {
	if t == nil {
		return errors.New("tree is nil")
	}
	if t.doc == nil {
		return errors.New("document is nil")
	}

	newList := []gedcom.Node{}
	for _, node := range t.doc.Nodes() {
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
func (t *Tree) PrepareExport() error {
	if t == nil {
		return errors.New("tree is nil")
	}
	if t.doc == nil {
		return errors.New("document is nil")
	}
	err := t.FilterChanges()
	if err != nil {
		return err
	}
	//
	notuniques, err := t.FixNodesUniqueness()
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
func (t *Tree) ReplaceReferences(oldref string, newref string) error {
	if t == nil {
		return errors.New("tree is nil")
	}
	return replaceReference(t.doc, oldref, newref)
}
