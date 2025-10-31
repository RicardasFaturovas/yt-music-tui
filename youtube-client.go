package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type SearchResult struct {
	Title   string `json:"title"`
	VideoId string `json:"videoId"`
}

type YoutubeVideo struct {
	AdaptiveFormats []AdaptiveFormat `json:"adaptiveFormats"`
}

type AdaptiveFormat struct {
	Url  string `json:"url"`
	Type string `json:"type"`
}

func getSearchResults(query string) []SearchResult {
	var results []SearchResult
	err := getRequest("/api/v1/search", query, &results)
	if err != nil {
		log.Fatal(err)
	}

	return results
}

func getUrl(videoId string) string {
	var result YoutubeVideo
	err := getRequest("/api/v1/videos/"+videoId, "", &result)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range result.AdaptiveFormats {
		if len(f.Type) > 5 && f.Type[:5] == "audio" {
			log.Println("AUDIO FORMAT FOUND")
			return f.Url
		}
	}

	return ""
}

// TODO: Cleanup error handling
func getRequest(endpoint string, query string, v any) error {
	youtubeClient := &http.Client{}

	var baseUrl = getConfig().InvidiousUrl
	req, err := http.NewRequest(http.MethodGet, baseUrl+endpoint, nil)
	if err != nil {
		log.Fatal(err)
	}

	q := req.URL.Query()
	q.Add("q", query)
	req.URL.RawQuery = q.Encode()
	res, getErr := youtubeClient.Do(req)
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
