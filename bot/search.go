package bot

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

type Results struct {
	Kind               string            `json:"kind"`
	Urls               Url               `json:"url"`
	Items              []Item            `json:"items"`
	Queries            Query             `json:"queries"`
	Contexts           Context           `json:"context"`
	SearchInformations SearchInformation `json:"searchInformation"`
}

type SearchInformation struct {
	SearchTime            float64 `json:"searchTime"`
	FormattedSearchTime   string  `json:"formattedSearchTime"`
	TotalResults          string  `json:"totalResults"`
	FormattedTotalResults string  `json:"formattedTotalResults"`
}

type Context struct {
	Title string `json:"title"`
}

type Query struct {
	Requests []Request `json:"requests"`
}

type Request struct{}

type Url struct {
	Type     string `json:"type"`
	Template string `json:"template"`
}

type Item struct {
	Kind        string `json:"kind"`
	Title       string `json:"title"`
	HtmlTitle   string `json:"htmlTitle"`
	Link        string `json:"link"`
	DisplayLink string `json:"displayLink"`
	Snippet     string `json:"snippet"`
	HtmlSnippet string `json:"htmlSnippet"`
	Mime        string `json:"mime"`
	FileFormat  string `json:"fileFormat"`
	Images      Image  `json:"image"`
}

type Image struct {
	ContextLink     string `json:"contextLink"`
	Height          int    `json:"height"`
	Width           int    `json:"width"`
	ByteSize        int    `json:"byteSize"`
	ThumbnailLink   string `json:"thumbnailLink"`
	ThumbnailHeight int    `json:"thumbnailHeight"`
	ThumbnailWidth  int    `json:"thumbnailWidth"`
}

func (u *Update) Search() (string, error) {
	// strip /search endpoint from string
	searchText := strings.ReplaceAll(u.Message.Text, "/search ", "")
	if searchText == "" {
		return "wait, what am I searching here?", nil
	}

	client := &http.Client{}

	req, err := http.NewRequest("GET", "https://www.googleapis.com/customsearch/v1", nil)
	if err != nil {
		return "", err
	}

	searchParams := req.URL.Query()
	searchParams.Add("key", os.Getenv("SEARCH_API_KEY"))
	searchParams.Add("cx", "a7921bdc020aa4524")
	searchParams.Add("searchType", "image")
	searchParams.Add("q", searchText)

	req.URL.RawQuery = searchParams.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var r Results
	if err := json.Unmarshal(content, &r); err != nil {
		return "", err
	}

	// randomize search results of 10
	rand.Seed(time.Now().UnixNano())
	random := rand.Intn(10)

	for index, img := range r.Items {
		if index == random {
			return img.Link, nil
		}
	}
	//return string(content), nil
	return "", nil
}
