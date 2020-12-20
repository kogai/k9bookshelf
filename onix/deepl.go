package onix

import (
	"encoding/json"
	"net/http"
	"os"
)

const deeplEndpoint = "https://api.deepl.com/v2"

// DeepLTranslate represents response of translate API.
type DeepLTranslate struct {
	Translations []struct {
		// "detected_source_language":"EN",
		Text string `json:"text"`
	} `json:"translations"`
}

// Translate send request to DeepL
func Translate(raw string) (*string, error) {
	authKey := os.Getenv("DEEPL_AUTH_KEY")
	req, err := http.NewRequest("GET", deeplEndpoint+"/translate", nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Add("auth_key", authKey)
	q.Add("text", raw)
	q.Add("source_lang", "EN")
	q.Add("target_lang", "JA")
	req.URL.RawQuery = q.Encode()

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	var r DeepLTranslate
	err = json.NewDecoder(res.Body).Decode(&r)
	if err != nil {
		return nil, err
	}
	return &r.Translations[0].Text, nil
}
