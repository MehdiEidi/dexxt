package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func getFarsiAPI(finglish string) (farsi string, err error) {
	body := strings.NewReader(finglish)

	req, err := http.NewRequest("POST", "https://9mkhzfaym3.execute-api.us-east-1.amazonaws.com/production/convert?", body)
	if err != nil {
		err = fmt.Errorf("err creating request: %w", err)
		return
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:102.0) Gecko/20100101 Firefox/102.0")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Content-Type", "text/plain")
	req.Header.Set("Origin", "https://behnevis.com")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Referer", "https://behnevis.com/")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "cross-site")
	req.Header.Set("Te", "trailers")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		err = fmt.Errorf("err sending request: %w", err)
		return
	}
	defer resp.Body.Close()

	resp_body, err := io.ReadAll(resp.Body)
	if err != nil {
		err = fmt.Errorf("err reading response body: %w", err)
		return
	}

	var result map[string]string

	err = json.Unmarshal(resp_body, &result)
	if err != nil {
		err = fmt.Errorf("err unMarshaling response body: %w", err)
		return
	}

	for _, v := range result {
		farsi += v
	}

	return
}
