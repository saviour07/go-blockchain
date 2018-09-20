package id

import (
	"encoding/json"
)

// Identity structure
type Identity struct {
	Name string
}

func (identity Identity) String() string {
	output, err := json.Marshal(identity)
	if err != nil {
		return ""
	}
	return string(output)
}
