package entity

import (
	"strings"

	mapset "github.com/deckarep/golang-set/v2"
)

type IDSet struct {
	Set mapset.Set[ID]
}

// NewIDSet create a new entity IDSet
func NewIDSet() IDSet {
	return IDSet{
		Set: mapset.NewSet[ID](),
	}
}

// Add
func (i *IDSet) Add(id ID) IDSet {
	i.Set.Add(id)
	return *i
}

// StringToIDSet convert a string to an entity IDSet
func StringToIDSet(IDSetSplitedByComma string) (IDSet, error) {
	IDSetString := strings.Split(IDSetSplitedByComma, ",")
	IDSet := NewIDSet()
	for _, IDString := range IDSetString {
		newID, err := StringToID(IDString)
		if err != nil {
			return NewIDSet(), err
		}
		IDSet = IDSet.Add(newID)
	}
	return IDSet, nil
}

// Strings
func (i *IDSet) Strings() []string {
	idSetString := []string{}
	for id := range i.Set.Iter() {
		idSetString = append(idSetString, id.String())
	}
	return idSetString
}
