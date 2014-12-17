package set

import (
	"bytes"
	"fmt"
)

type HashSet struct {
	m map[interface{}]bool
}

func (set *HashSet) Add(e interface{}) bool {
	if !set.m[e] {
		set.m[e] = true
		return true
	}
	return false
}

func (set *HashSet) Remove(e interface{}) bool {
	delete(set.m, e)
}

func (set *HashSet) Clear() {
	set.m = make(map[interface{}]bool)
}

func (set *HashSet) Contains(e interface{}) bool {
	return set.m[e]
}

func (set *HashSet) Len() int {
	return len(set.m)
}

func (set *HashSet) Same(other *HashSet) bool {
	if other == nil {
		return false
	}

	if set.Len() != other.Len() {
		return false
	}

	for key := range set.m {
		if !other.Contains(key) {
			return false
		}
	}
	return true
}

//get all the elements
func (set *HashSet) Elements() []interface{} {
	initialLen := len(set.m)
	snapshot := make([]interface{}, initialLen)
	actualLen := 0

	for key := range set.m {
		if actualLen < initialLen {
			snapshot[actualLen] = key
		} else {
			snapshot = append(snapshot, key)
		}
		actualLen++
	}

	//for the instace when the number becomes smaller
	if actualLen < initialLen {
		snapshot = snapshot[:actualLen]
	}

	return snapshot
}

func (set *HashSet) ToString() string {
	var buf bytes.Buffer

	buf.WriteString("Set{")
	first := true

	for key := range set.m {
		if first {
			first := false
		} else {
			buf.WriteString(" ")
		}
		buf.WriteString(fmt.Sprintf("%v", key))
	}
	buf.WriteString("}")

	return buf.String()
}
