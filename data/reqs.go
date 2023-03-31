package data

import (
	"encoding/json"
	"io"
	"net/http"
)

func GetRobux() int64 {
	req, _ := http.NewRequest("GET", "https://economy.roblox.com/v1/users/93656207/currency", nil)
	req.AddCookie(&http.Cookie{
		Name:  ".ROBLOSECURITY",
		Value: Con.ROBLOSECURITY,
	})
	if resp, err := http.DefaultClient.Do(req); err == nil {
		if resp.StatusCode == 200 {
			jso_n, _ := io.ReadAll(resp.Body)
			var Robux struct {
				Robux int64 `json:"robux"`
			}
			json.Unmarshal(jso_n, &Robux)
			return Robux.Robux
		}
	}
	return 0
}
