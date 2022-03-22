package data

import (
	"encoding/json"

	"github.com/pkg/errors"
)

type PeerList struct {
	SelfId  string
	PeerIds []string
	Length  int
}

// PeerIdsToJson ...
func PeerIdsToJson(peerList []string) []byte {
	value, err := json.Marshal(peerList)
	if err != nil {
		return []byte{}
	}
	return value
}

// JsonToPeerIds ...
func JsonToPeerIds(inputJson []byte) ([]string, error) {
	var ques []string
	err := json.Unmarshal(inputJson, &ques)
	if err != nil {
		return ques, errors.New("Cannot decode Json to Registered Data")
	}
	return ques, nil
}

func (pl *PeerList) InsertToList(peerList []string) {

	for _, port := range peerList {

		if port != pl.SelfId {

			pl.PeerIds = append(pl.PeerIds, port)
		}
	}
}

func (pl *PeerList) ToJson() []byte {
	value, _ := json.Marshal(pl)
	return value
}

func DecodeJson(temp []byte) PeerList {

	var templist PeerList
	_ = json.Unmarshal(temp, &templist)

	return templist

}
