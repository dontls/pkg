package tree

import (
	"encoding/json"
	"fmt"
	"testing"
)

type Menu struct {
	ID       uint
	PID      uint
	Name     string
	Children []Menu `json:"children,omitempty"`
}

func TestTree(t *testing.T) {
	v := []Menu{{
		ID:   1,
		PID:  0,
		Name: "0-1",
	}, {
		ID:   2,
		PID:  1,
		Name: "1-2",
	}, {
		ID:  3,
		PID: 2,
		Name: "2-3",
	}, {
		ID:   4,
		PID:  1,
		Name: "1-4",
	}, {
		ID:   5,
		PID:  2,
		Name: "1-5",
	}, {
		ID:   6,
		PID:  0,
		Name: "0-6",
	}}
	r := Slice(v, func(i int) bool { return v[i].PID == 0 }, func(i, j int) bool { return v[i].ID == v[j].PID })
	b, _ := json.MarshalIndent(r, "", "  ")
	fmt.Printf("%s\n", b)
}
