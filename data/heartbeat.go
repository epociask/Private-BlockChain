package data

import (
	"encoding/json"
)

type HeartBeatData struct {
	SendersId string
	BlockJson string
}

// HeartBeatDataToJson ...
func (q *HeartBeatData) HeartBeatDataToJson() []byte {
	value, _ := json.Marshal(q)
	return value
}
