package messages

import (
	"fmt"
)

type ServiceAddr struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

func (s ServiceAddr) String() string {
	return fmt.Sprintf("%v:%v", s.Host, s.Port)
}

type FrequencyRequest struct {
	Url  string `json:"url"`
	Text string `json:"text"`
}
func (r FrequencyRequest) String() string {
	text := "empty"
	if r.Text != ""{
		text = "..."
	}
	return fmt.Sprintf("Url: %v, Text: %v", r.Url, text)
}

type FrequencyResponse struct {
	Result map[string]int64 `json:"result"`
	ErrorCode int `json:"errorCode"`
	Error  string `json:"error"`
}
