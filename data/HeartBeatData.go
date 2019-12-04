package data

import (
	"encoding/json"
	//"github.com/pkg/errors"
)

type HeartBeatData struct {

	SendersId string
	BlockJson string
}


//HeartBeatData ----> JSON
func (q *HeartBeatData) HeartBeatDataToJson() []byte {
	value, _ := json.Marshal(q)
	return value
}



