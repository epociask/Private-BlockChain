package data

import (
	"encoding/json"

	"github.com/pkg/errors"
)

type RegisterData struct {
	AssignedId  string
	PeerMapJson string
}

// ToJson ...
func (q *RegisterData) ToJson() ([]byte, error) {
	value, err := json.Marshal(q)
	if err != nil {
		return []byte{}, errors.New("Cannot encode Registered Data to Json")
	}
	return value, nil
}

// RegisteredDataFromJson ...
func RegisteredDataFromJson(inputJson []byte) (RegisterData, error) {
	ques := RegisterData{}
	err := json.Unmarshal(inputJson, &ques)
	if err != nil {
		return ques, errors.New("Cannot decode Json to Registered Data")
	}
	return ques, nil
}
