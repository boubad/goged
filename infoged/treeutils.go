package infoged

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/elliotchance/gedcom"
)

// //////////////////////
type RootNodeStats struct {
	Tag           string
	Ids           map[string]string
	Refs          map[string]string
	Missing       []string
	NotReferenced []string
} //
func newRootNodeStats(tag string) *RootNodeStats {
	return &RootNodeStats{
		Tag:           tag,
		Ids:           map[string]string{},
		Refs:          map[string]string{},
		Missing:       []string{},
		NotReferenced: []string{}}
} // newRootNodeStats
func PrintParseData(pdata *map[string]*RootNodeStats) {
	for _, v := range *pdata {
		if len(v.Missing) > 0 || len(v.NotReferenced) > 0 {
			fmt.Printf("NodeType: %s\n", v.Tag)
			if len(v.Missing) > 0 {
				fmt.Printf("\tmissing IDs: %v\n", v.Missing)
			}
			if len(v.NotReferenced) > 0 {
				fmt.Printf("\tnotrefrenced IDs: %v\n", v.NotReferenced)
			}
		}
	}
} // printParseData
// /////////////////////////
func parseReference(node gedcom.Node) (string, bool) {
	if node != nil {
		value := strings.TrimSpace(node.Value())
		if len(value) > 2 && strings.HasPrefix(value, "@") && strings.HasSuffix(value, "@") {
			return strings.Trim(value, "@"), true
		}
	}
	return "", false
} // parseReference

// /////////////////////////
func replaceReference(doc *gedcom.Document, oldref string, newref string) error {
	if doc == nil {
		return  errors.New("node is nil")
	}
	//
	filterFunc := func(node gedcom.Node) (gedcom.Node, bool) {
		if node == nil {
			return nil, false
		}
		id := node.Pointer()
		if id == oldref {
			return gedcom.NewNode(node.Tag(), node.Value(), newref), false
		}
		if xref, ok := parseReference(node); ok {
			if xref == oldref {
				return gedcom.NewNode(node.Tag(), node.Value(), ""), false
			}
		}
		return node, true
	} // filterFunc
	newList := []gedcom.Node{}
	for _, node := range doc.Nodes() {
		xnode := gedcom.Filter(node, filterFunc)
		if xnode != nil {
			newList = append(newList, xnode)
		}
	}
	doc.SetNodes(newList)
	return nil
} // replaceReference

// ////////////////////////
func findNotExistingFiles(pMap *map[string]string, node gedcom.Node, baseDir string) {
	if node == nil {
		return
	}
	for _, child := range node.Nodes() {
		findNotExistingFiles(pMap, child, baseDir)
	}
	if !node.Tag().Is(gedcom.TagFile) {
		return
	}
	value := strings.TrimSpace(node.Value())
	if value == "" {
		return
	}
	filePath := filepath.Join(baseDir, value)
	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		(*pMap)[value] = value
	}
} // findNonExistringFiles
func findAllFiles(pMap *map[string]string, node gedcom.Node) {
	if node == nil {
		return
	}
	for _, child := range node.Nodes() {
		findAllFiles(pMap, child)
	}
	if !node.Tag().Is(gedcom.TagFile) {
		return
	}
	value := strings.TrimSpace(node.Value())
	if value != "" {
		(*pMap)[value] = value
	}
} // findAllExistingFiles
func deleteReferences(node gedcom.Node, id string) {
	if node == nil {
		return
	}
	if len(node.Nodes()) < 1 {
		return
	}
	for _, child := range node.Nodes() {
		deleteReferences(child, id)
	}
	foundNodes := []gedcom.Node{}
	for _, child := range node.Nodes() {
		bFound := false
		if sref, ok := parseReference(child); ok {
			if sref == id {
				bFound = true
				break
			}
		}
		if !bFound {
			foundNodes = append(foundNodes, child)
		}
	} // child
	if len(foundNodes) > 0 {
		node.SetNodes(foundNodes)
	} // child
} //deleteReferences
func deleteDocReferences(doc *gedcom.Document, id string) {
	for _, node := range doc.Nodes() {
		deleteReferences(node, id)
	}
} //deleteDocReferences
func deleteDocRootNode(doc *gedcom.Document, id string) {
	deleteDocReferences(doc, id)
	oldNodes := doc.Nodes()
	newNodes := []gedcom.Node{}
	for _, node := range oldNodes {
		if node == nil {
			continue
		}
		if node.Pointer() != id {
			newNodes = append(newNodes, node)
		}
	}
	if len(newNodes) != len(oldNodes) {
		doc.SetNodes(newNodes)
	}
} //deleteDocRootNode
func parseNode(pdata *map[string]*RootNodeStats, node gedcom.Node) {
	if node != nil {
		for _, child := range node.Nodes() {
			parseNode(pdata, child)
		}
		tag := node.Tag()

		if tag.Is(gedcom.TagHusband) ||
			tag.Is(gedcom.TagWife) ||
			tag.Is(gedcom.TagChild) ||
			tag.Is(gedcom.TagAlias) ||
			tag.Is(gedcom.TagAssociates) {
			if sref, ok := parseReference(node); ok {
				key := strings.TrimSpace(gedcom.TagIndividual.String())
				pcur, ok := (*pdata)[key]
				if !ok {
					pcur = newRootNodeStats(key)
					pcur.Refs[sref] = sref
					(*pdata)[key] = pcur
				} else {
					pcur.Refs[sref] = sref
					(*pdata)[key] = pcur
				}
			}
			return
		}
		if tag.Is(gedcom.TagFamilyChild) || tag.Is(gedcom.TagFamilySpouse) {
			if sref, ok := parseReference(node); ok {
				key := strings.TrimSpace(gedcom.TagFamily.String())
				pcur, ok := (*pdata)[key]
				if !ok {
					pcur = newRootNodeStats(key)
					pcur.Refs[sref] = sref
					(*pdata)[key] = pcur
				} else {
					pcur.Refs[sref] = sref
					(*pdata)[key] = pcur
				}
			}
			return
		}
		id := strings.TrimSpace(node.Pointer())
		if len(id) > 0 {
			key := strings.TrimSpace(tag.String())
			pcur, ok := (*pdata)[key]
			if !ok {
				pcur = newRootNodeStats(key)
				pcur.Ids[id] = id
				(*pdata)[key] = pcur
			} else {
				pcur.Ids[id] = id
				(*pdata)[key] = pcur
			}
			return
		}
		if sref, ok := parseReference(node); ok {
			key := strings.TrimSpace(tag.String())
			pcur, ok := (*pdata)[key]
			if !ok {
				pcur = newRootNodeStats(key)
				pcur.Refs[sref] = sref
				(*pdata)[key] = pcur
			} else {
				pcur.Refs[sref] = sref
				(*pdata)[key] = pcur
			}
			return
		}

	}
} // parseNode
func parseDocIds(doc *gedcom.Document) (pdata map[string]*RootNodeStats, err error) {
	if doc == nil {
		return nil, errors.New("document is nil")
	}
	pdata = map[string]*RootNodeStats{}
	for _, node := range doc.Nodes() {
		parseNode(&pdata, node)
	}
	for _, pitem := range pdata {
		for _, id := range pitem.Refs {
			if _, ok := pitem.Ids[id]; !ok {
				pitem.Missing = append(pitem.Missing, id)
			}
		}
		for _, id := range pitem.Ids {
			if _, ok := pitem.Refs[id]; !ok {
				pitem.NotReferenced = append(pitem.NotReferenced, id)
			}
		}

	} // pitem
	return pdata, nil
} // parseDoc
func fixReferences(doc *gedcom.Document) (map[string]*RootNodeStats, error) {
	pdata, err := parseDocIds(doc)
	if err != nil {
		return nil, err
	}
	//printParseData(&pdata)
	for _, pitem := range pdata {
		if len(pitem.Missing) > 0 {
			for _, id := range pitem.Missing {
				deleteDocReferences(doc, id)
			}
		}
		if len(pitem.NotReferenced) > 0 {
			for _, id := range pitem.NotReferenced {
				deleteDocRootNode(doc, id)
			}
		}
	} // item
	return pdata, nil
} // fixReferences
