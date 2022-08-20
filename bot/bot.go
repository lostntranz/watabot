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
	incoming, err := parseTelegramRequest(r)
	if err != nil {
		return
	}

	log.Printf("echoing incoming message: %v", incoming.Message.Text)
	sendRespTelegram(incoming.Message.Chat.Id, "嘥撚氣,，算吧啦!")
}

func parseTelegramRequest(r *http.Request) (*Update, error) {
	var u Update
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		log.Printf("couldn't decode incoming message")
		return nil, err
	}
	return &u, nil
}

func sendRespTelegram(id int, response string) (string, error) {

	var apiUrl = fmt.Sprintf("https://api.telegram.org/bot" + os.Getenv("TELEGRAM_BOT_TOKEN") + "/sendMessage")
	resp, err := http.PostForm(
		apiUrl,
		url.Values{
			"chat_id": {strconv.Itoa(id)},
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
