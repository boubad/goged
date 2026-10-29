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
// //////////////////////////////
func collectNodesByFilter(node gedcom.Node, pResult *[]gedcom.Node, filterFunc func(node gedcom.Node) bool) {
	if node != nil {
		nodes := node.Nodes()
		for _, child := range nodes {
			collectNodesByFilter(child, pResult, filterFunc)
		}
		if filterFunc(node) {
			*pResult = append(*pResult, node)
		}
	}
} // collectNodesByFilter
func CollectDocumentNodesByFilter(doc *gedcom.Document, pResult *[]gedcom.Node, filterFunc func(node gedcom.Node) bool) error {
	if pResult == nil {
		return errors.New("result is nil")
	}
	*pResult = []gedcom.Node{}
	if doc == nil {
		return errors.New("document is nil")
	}
	nodes := doc.Nodes()
	for _, node := range nodes {
		collectNodesByFilter(node, pResult, filterFunc)
	}
	return nil
} // CollectDocumentNodes
// ////////////////////////
func findNotExistingFiles(pMap *map[string]string, node gedcom.Node, baseDir string) {
	if node != nil {
		nodes := node.Nodes()
		for _, child := range nodes {
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
	}
} // findNonExistringFiles
func findAllFiles(pMap *map[string]string, node gedcom.Node) {
	if node != nil {
		nodes := node.Nodes()
		for _, child := range nodes {
			findAllFiles(pMap, child)
		}
		if !node.Tag().Is(gedcom.TagFile) {
			return
		}
		value := strings.TrimSpace(node.Value())
		if value != "" {
			(*pMap)[value] = value
		}
	}
} // findAllExistingFiles
func deleteReferences(node gedcom.Node, id string) {
	if node != nil {
		nodes := node.Nodes()
		if len(nodes) < 1 {
			return
		}
		for _, child := range nodes {
			deleteReferences(child, id)
		}
		foundNodes := []gedcom.Node{}
		for _, child := range nodes {
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
		node.SetNodes(foundNodes)
	} // node
} //deleteReferences
func deleteDocReferences(doc *gedcom.Document, id string) {
	if doc == nil {
		return
	}
	nodes := doc.Nodes()
	for _, node := range nodes {
		deleteReferences(node, id)
	}
} //deleteDocReferences
func deleteDocRootNode(doc *gedcom.Document, id string) {
	if doc == nil {
		return
	}
	sid := strings.TrimSpace(id)
	if len(sid) < 1 {
		return
	}
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
// ////////////////////
func parseNode(pdata *map[string]*RootNodeStats, node gedcom.Node) {
	if node != nil {
		nodes := node.Nodes()
		for _, child := range nodes {
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
				if !ok || pcur == nil {
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
				if !ok || pcur == nil {
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
			if !ok || pcur == nil {
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
			if !ok || pcur == nil {
				pcur = newRootNodeStats(key)
				pcur.Refs[sref] = sref
				(*pdata)[key] = pcur
			} else {
				pcur.Refs[sref] = sref
				(*pdata)[key] = pcur
			}
		}
	}
} // parseNode
func parseDocIds(doc *gedcom.Document) (pdata map[string]*RootNodeStats, err error) {
	if doc == nil {
		return nil, errors.New("document is nil")
	}
	pdata = map[string]*RootNodeStats{}
	nodes := doc.Nodes()
	for _, node := range nodes {
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
func cleanNodesFilter(node gedcom.Node) (gedcom.Node, bool) {
	if node == nil {
		return nil, false
	}
	tag := node.Tag()
	if !tag.IsOfficial() || tag.Is(gedcom.TagFact) {
		return nil, false
	}
	/*
		if !tag.IsOfficial() || tag.Is(gedcom.TagChange) || tag.Is(gedcom.TagFact) {
			return nil, false
		}
	*/
	if tag.Is(gedcom.TagFormat) {
		value := strings.ToUpper(node.Value())
		if value == "JPEG" {
			return gedcom.NewNode(tag, "jpg", ""), false
		}
	}
	if tag.Is(gedcom.TagDestination) {
		return gedcom.NewNode(tag, "ANY", ""), false
	}
	if tag.Is(gedcom.TagPedigree) {
		return gedcom.NewNode(tag, strings.ToLower(node.Value()), ""), false
	}
	if tag.Is(gedcom.TagName) {
		value := strings.ToUpper(node.Value())
		if strings.Contains(value, "?") ||
			strings.Contains(value, "N/A") ||
			strings.Contains(value, "(INC)") ||
			strings.Contains(value, "UNKNOWN") {
			return nil, false
		}
		nameNode, ok := node.(*gedcom.NameNode)
		if !ok {
			return nil, false
		}
		sgiven := nameNode.GivenName()
		ssurname := nameNode.Surname()
		snewName := ""
		if len(sgiven) < 1 && len(ssurname) < 1 {
			return nil, false
		} else if len(sgiven) < 1 {
			snewName = fmt.Sprintf("/%s/", ssurname)
		} else if len(ssurname) < 1 {
			snewName = sgiven
		} else {
			snewName = fmt.Sprintf("%s /%s/", sgiven, ssurname)
		}
		if snewName != node.Value() {
			return gedcom.NewNode(tag, snewName, ""), false

		} else {
			return node, true
		}
	} // Name tag
	if tag.Is(gedcom.TagNote) {
		initial_value := strings.TrimSpace(node.Value())
		if len(node.Pointer()) < 1 && len(initial_value) < 1 {
			return nil, false
		}
	}
	if len(node.Nodes()) == 0 && len(node.Pointer()) == 0 {
		if !tag.Is(gedcom.TagTrailer) && !tag.Is(gedcom.TagHeader) {
			initial_value := strings.TrimSpace(node.Value())
			if len(initial_value) < 1 {
				return nil, false
			}
		}
	}
	// nested tags
	if len(node.Nodes()) > 0 {
		found := false
		for _, child := range node.Nodes() {
			if child != nil && child.Tag().Is(tag) {
				found = true
				break
			}
		}
		if found {
			return nil, false
		}
	} // nodes
	return node, true
} // cleanNodesFilter
func cleanDuplicatesFilter(node gedcom.Node) (gedcom.Node, bool) {
	if node == nil {
		return node, false
	}
	oldList := node.Nodes()
	n := len(oldList)
	if n < 1 {
		return node, false
	}
	newList := []gedcom.Node{}
	for i := 0; i < n; i++ {
		inode := oldList[i]
		if inode == nil {
			continue
		}
		if len(inode.Nodes()) > 0 {
			newList = append(newList, inode)
			continue
		}
		istring := inode.GEDCOMLine(gedcom.NoIndent)
		found := false
		for j := i + 1; j < n; j++ {
			jnode := oldList[j]
			if jnode == nil {
				continue
			}
			if len(jnode.Nodes()) > 0 {
				continue
			}
			jstring := jnode.GEDCOMLine(gedcom.NoIndent)
			if istring == jstring {
				found = true
				break
			}
		} // j
		if !found {
			newList = append(newList, inode)
		}
	} // i
	if len(newList) != n {
		node.SetNodes(newList)
	}
	return node, true
} // cleanDuplicates
