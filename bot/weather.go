package bot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type Weather struct {
	LocalObservationDateTime string `json:"LocalObservationDateTime"`
	EpochTime                int64  `json:"EpochTime"`
	WeatherText              string `json:"WeatherText"`
	WeatherIcon              int32  `json:"WeatherIcon"`
	HasPrecipitation         bool   `json:"HasPrecipitation"`
	PrecipitationType        string `json:"PrecipitationType"`
	IsDayTime                bool   `json:"IsDayTime"`
	Temperature              struct {
		Metric struct {
			Value    float64
			Unit     string
			UnitType int32
		} `json:"Metric"`
		Imperial struct {
			Value    float64
			Unit     string
			UnitType int32
		} `json:"Imperial"`
	} `json:"Temperature"`
	MobileLink string `json:"MobileLink"`
	Link       string `json:"Link"`
}

func (u *Update) GetWeather() (string, error) {

	// strip /search endpoint from string
	searchText := strings.ReplaceAll(u.Message.Text, "/weather ", "")
	if searchText == "" {
		return "wait, what am I searching here?", nil
	}

	locationKey, city, err := getLocationKey(searchText)
	if err != nil {
		return "", nil
	}

	client := &http.Client{}
	weatherURL := fmt.Sprintf("http://dataservice.accuweather.com/currentconditions/v1/%s", locationKey)
	req, err := http.NewRequest("GET", weatherURL, nil)
	if err != nil {
		return "", err
	}

	searchParams := req.URL.Query()
	searchParams.Add("apikey", os.Getenv("WEATHER_API_KEY"))

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

	var d []Weather
	if err := json.Unmarshal(content, &d); err != nil {
		return "", err
	}
	result := fmt.Sprintf("Weather for %s is %v - currently %v deg", city, d[0].WeatherText, d[0].Temperature.Imperial.Value)
	return result, nil
}

func getLocationKey(zipcode string) (string, string, error) {

	client := &http.Client{}

	req, err := http.NewRequest("GET", "http://dataservice.accuweather.com/locations/v1/postalcodes/US/search", nil)
	if err != nil {
		return "", "", err
	}

	searchParams := req.URL.Query()
	searchParams.Add("apikey", os.Getenv("WEATHER_API_KEY"))
	searchParams.Add("q", zipcode)

	req.URL.RawQuery = searchParams.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}

	var c []map[string]json.RawMessage

	if err := json.Unmarshal(content, &c); err != nil {
		return "", "", nil
	}

	if len(c) != 1 {
		return "", "", fmt.Errorf("couldn't determine location")
	}
	sanitizeLocationKey := strings.Trim(string(c[0]["Key"]), "\"")
	city := strings.Trim(string(c[0]["LocalizedName"]), "\"")

	return sanitizeLocationKey, city, nil
}
