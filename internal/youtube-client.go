package internal

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type YoutubeSong struct {
	Title   string `json:"title"`
	VideoId string `json:"videoId"`
}

type YoutubeClient struct {
	client  *http.Client
	baseUrl string
}

func NewYoutubeClient(baseUrl string) *YoutubeClient {
	youtubeClient := &http.Client{}

	return &YoutubeClient{
		youtubeClient,
		baseUrl,
	}
}

func (y *YoutubeClient) GetSearchResults(query string) []YoutubeSong {
	var results []YoutubeSong
	err := y.getRequest("/api/v1/search", query, &results)
	if err != nil {
		log.Fatal(err)
	}

	return results
}

// TODO: Cleanup error handling
func (y *YoutubeClient) getRequest(endpoint string, query string, v any) error {
	req, err := http.NewRequest(http.MethodGet, y.baseUrl+endpoint, nil)
	if err != nil {
		log.Fatal(err)
	}

	q := req.URL.Query()
	q.Add("q", query)
	req.URL.RawQuery = q.Encode()
	res, getErr := y.client.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
		return readErr
	}

	jsonErr := json.Unmarshal(body, &v)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	return nil
}
