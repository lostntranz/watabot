package bot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type AQI struct {
	Datas  Data   `json:"data"`
	Status string `json:"status"`
}

type Data struct {
	City      string   `json:"city"`
	Country   string   `json:"country"`
	Current   Currents `json:"current"`
	Locations Location `json:"location"`
	State     string   `json:"state"`
}

type Currents struct {
	Pollutions  Pollution  `json:"pollution"`
	AQIWeathers AQIWeather `json:"weather"`
}

type AQIWeather struct {
	HU int32   `json:"hu"`
	IC string  `json:"ic"`
	PR int32   `json:"pr"`
	TP int32   `json:"tp"`
	TS string  `json:"ts"`
	WD int32   `json:"wd"`
	WS float64 `json:"ws"`
}

type Location struct {
	Coordinates []float64 `json:"coordinates"`
	Type        string    `json:"type"`
}

type Pollution struct {
	AqiCN     int32  `json:"aqicn"`
	AqiUS     int32  `json:"aqius"`
	MainCN    string `json:"maincn"`
	MainUS    string `json:"mainus"`
	TimeStamp string `json:"ts"`
}

func (u *Update) GetAQI() (string, error) {

	aqiSearch := strings.ReplaceAll(u.Message.Text, "/aqi ", "")
	if aqiSearch == "" {
		return "wait, what am I searching here?", nil
	}

	city := url.QueryEscape(aqiSearch)
	url := fmt.Sprintf("http://api.airvisual.com/v2/city?city=%s&state=California&country=USA&key=%s", city, os.Getenv("AQI_API_KEY"))
	fmt.Println(url)
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	fmt.Println(string(content))
	var a AQI
	if err := json.Unmarshal(content, &a); err != nil {
		return "", err
	}

	aqi := fmt.Sprintf("Current AQI for %s, %s: %d", a.Datas.City, a.Datas.State, a.Datas.Current.Pollutions.AqiUS)
	return aqi, nil

}
