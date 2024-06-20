package helper

import "encoding/json"

type ResponseBody struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func ResponseJSON(s interface{}) []byte {
	b, _ := json.Marshal(s)

	return b
}
