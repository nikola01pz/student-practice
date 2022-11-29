package http

import (
	"encoding/json"
	"net/http"
	"time"
)

func GetJson(url string, target interface{}) error {
	var client = &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(target)
}
