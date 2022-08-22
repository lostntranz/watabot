package bot

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

type JokeResults struct {
	CurrentPage  int     `json:"current_page"`
	Limit        int     `json:"limit"`
	NextPage     int     `json:"next_page"`
	PreviousPage int     `json:"previous_page"`
	Results      []Jokes `json:"results"`
	SearchTerm   string  `json:"search_term"`
	Status       int     `json:"status"`
	TotalJokes   int     `json:"total_jokes"`
	TotalPages   int     `json:"total_pages"`
}

type Jokes struct {
	Id   string `json:"id"`
	Joke string `json:"joke"`
}

func (u *Update) GetJoke() (string, error) {

	// strip /joke endpoint from string
	jokeSearch := strings.ReplaceAll(u.Message.Text, "/joke ", "")
	if jokeSearch == "" {
		return "wait, what am I searching here?", nil
	}

	client := &http.Client{}

	req, err := http.NewRequest("GET", "https://icanhazdadjoke.com/search", nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("Accept", "application/json")

	searchParams := req.URL.Query()
	searchParams.Add("term", jokeSearch)

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

	var j JokeResults
	if err := json.Unmarshal(content, &j); err != nil {
		return "", err
	}

	if j.TotalJokes < 1 {
		return "got no jokes for " + jokeSearch, nil
	}
	rand.NewSource(time.Now().UnixNano())
	random := rand.Intn(j.TotalJokes)

	return j.Results[random].Joke, nil
}
