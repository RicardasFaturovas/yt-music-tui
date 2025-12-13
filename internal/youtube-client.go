package internal

import (
	"encoding/json"
	"errors"
	"fmt"
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

var GetSongsError = errors.New("GetSongsError")

func NewYoutubeClient(baseUrl string) *YoutubeClient {
	youtubeClient := &http.Client{}

	return &YoutubeClient{
		youtubeClient,
		baseUrl,
	}
}

func (y *YoutubeClient) GetSearchResults(query string) ([]YoutubeSong, error) {
	var results []YoutubeSong
	err := y.getRequest("/api/v1/search", query, &results)

	return results, err
}

func (y *YoutubeClient) getRequest(endpoint string, query string, v any) error {
	req, err := http.NewRequest(http.MethodGet, y.baseUrl+endpoint, nil)
	if err != nil {
		log.Println(fmt.Errorf("Building request: %w", err))
		return GetSongsError
	}

	q := req.URL.Query()
	q.Add("q", query)
	req.URL.RawQuery = q.Encode()
	res, getErr := y.client.Do(req)
	if getErr != nil {
		log.Println(fmt.Errorf("Executing request: %w", err))
		return GetSongsError
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		log.Println(fmt.Errorf("Reading response body: %w", err))
		return GetSongsError
	}

	jsonErr := json.Unmarshal(body, &v)
	if jsonErr != nil {
		log.Println(fmt.Errorf("Unmarshalling: %w", err))
		return GetSongsError
	}
	return nil
}
