package id

import (
	"encoding/json"
)

// IdentityVersion for updates
const IdentityVersion = 1

// Identity structure
type Identity struct {
	Name    string
	Version int
}

// ToIdentity converts the json input to an Identity object or returns an error on failure
func ToIdentity(input string) (Identity, error) {
	id := Identity{}
	err := json.Unmarshal([]byte(input), &id)
	if err != nil {
		return Identity{}, err
	}
	return id, nil
}

func (identity Identity) String() string {
	output, err := json.Marshal(identity)
	if err != nil {
		return ""
	}
	return string(output)
}
