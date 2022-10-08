package bot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

const (
	TODAY_TIDE = "https://www.deltaboating.com/tides/sac.php"
)

type Update struct {
	UpdateId int     `json:"update_id"`
	Message  Message `json:"message"`
}

type Message struct {
	Text string `json:"text"`
	Chat Chat   `json:"chat"`
}

type Chat struct {
	Id int `json:"id"`
}

func HandleTelegramWebhook(w http.ResponseWriter, r *http.Request) {
	incoming, err := parseChatMessages(r)
	if err != nil {
		return
	}
	log.Printf("echoing incoming message: %v", incoming.Message.Text)

	switch true {
	case strings.Contains(incoming.Message.Text, "/tide"):
		incoming.Respond(TODAY_TIDE)
	case strings.Contains(incoming.Message.Text, "/joke"):
		results, err := incoming.GetJoke()
		if err != nil {
			incoming.Respond("Ain't no joke, something is wrong here")
		}
		incoming.Respond(results)
	case strings.Contains(incoming.Message.Text, "/search"):
		results, err := incoming.Search()
		if err != nil {
			incoming.Respond("Some Ting Wong")
		}
		incoming.Respond(results)
	case strings.Contains(incoming.Message.Text, "/weather"):
		results, err := incoming.GetWeather()
		if err != nil {
			incoming.Respond("Can't get weather info")
		}
		incoming.Respond(results)
	case strings.Contains(incoming.Message.Text, "/aqi"):
		results, err := incoming.GetAQI()
		if err != nil {
			incoming.Respond("Can't get AQI info")
		}
		incoming.Respond(results)
	default:
		incoming.Respond("嘥撚氣,，算吧啦!")
	}
}

func parseChatMessages(r *http.Request) (*Update, error) {
	var u Update
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		log.Printf("couldn't decode incoming message")
		return nil, err
	}
	return &u, nil
}

func (u *Update) Respond(response string) (string, error) {

	var apiUrl = fmt.Sprintf("https://api.telegram.org/bot" + os.Getenv("TELEGRAM_BOT_TOKEN") + "/sendMessage")
	resp, err := http.PostForm(
		apiUrl,
		url.Values{
			"chat_id": {strconv.Itoa(u.Message.Chat.Id)},
			"text":    {response},
		},
	)
	if err != nil {
		log.Printf("couldnt send response back to telegram: %v", err)
		return "", err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("couldn't read post response body: %v", err)
		return "", err
	}

	content := string(respBody)
	log.Printf("Response post body: %s", content)

	return content, nil

}
