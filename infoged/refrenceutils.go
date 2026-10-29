package infoged

///////////////////
import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/elliotchance/gedcom"
)

func getNextAvailableNumberID(doc *gedcom.Document, prefix string) int {
	nmax := 0
	nodes := doc.Nodes()
	for _, node := range nodes {
		id := strings.TrimSpace(node.Pointer())
		if strings.HasPrefix(id, prefix) {
			n, err := strconv.Atoi(id[len(prefix):])
			if err == nil {
				if n > nmax {
					nmax = n
				}
			}
		}
	} // node
	nmax = nmax + 1
	return nmax
} // getNextAvailableNumberID
func fixNonConformIds(doc *gedcom.Document, tag gedcom.Tag, prefix string) *gedcom.Document {
	oMap := map[string]string{}
	nmax := getNextAvailableNumberID(doc, prefix)
	nlen := len(prefix) + 5
	nodes := doc.Nodes()
	for _, node := range nodes {
		if node == nil {
			continue
		}
		if !node.Tag().Is(tag) {
			continue
		}
		id := strings.TrimSpace(node.Pointer())
		if len(id) > 0 && (!strings.HasPrefix(id, prefix) || len(id) < nlen) {
			oldId := "@" + id + "@"
			newId := fmt.Sprintf("@%s%05d@", prefix, nmax)
			nmax = nmax + 1
			oMap[oldId] = newId
		}
	} // node
	if len(oMap) < 1 {
		return doc
	}
	sdoc := doc.String()
	for oldId, newId := range oMap {
		sdoc = strings.ReplaceAll(sdoc, oldId, newId)
	}
	newDoc, err := gedcom.NewDocumentFromString(sdoc)
	if err != nil {
		return nil
	}
	return newDoc
} // fixNonConformIds
func fixDocIDs(doc *gedcom.Document) (*gedcom.Document, error) {
	if doc == nil {
		return nil, errors.New("document is nil")
	}
	xdoc := fixNonConformIds(doc, gedcom.TagIndividual, "I")
	if xdoc == nil {
		return nil, errors.New("failed to fix individual IDs")
	}
	xdoc = fixNonConformIds(xdoc, gedcom.TagFamily, "F")
	if xdoc == nil {
		return nil, errors.New("failed to fix family IDs")
	}
	xdoc = fixNonConformIds(xdoc, gedcom.TagObject, "M")
	if xdoc == nil {
		return nil, errors.New("failed to fix object IDs")
	}
	xdoc = fixNonConformIds(xdoc, gedcom.TagNote, "N")
	if xdoc == nil {
		return nil, errors.New("failed to fix note IDs")
	}
	return xdoc, nil
} // fixDocIDs

// ////////////////////////
