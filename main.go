package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

func main() {
	lat, long := getLocation()
	temp, wind, clouds := getWeather(lat, long)

	fmt.Printf("It is currently %gF with gusts to %gmph and %g%% cloud cover.", temp, wind, clouds)
}

func getLocation() (lat string, long string) {
	type Location struct {
		Latitude  float32 `json:"latitude"`
		Longitude float32 `json:"longitude"`
	}

	resp, err := http.Get("https://ipapi.co/json")
	if err != nil {
		fmt.Printf("error getting location: %s", err)
	}

	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var location Location
	err = json.Unmarshal(body, &location)
	if err != nil {
		panic(err)
	}

	lat = strconv.FormatFloat(float64(location.Latitude), 'g', -1, 32)
	long = strconv.FormatFloat(float64(location.Longitude), 'g', -1, 32)
	return
}

func getWeather(lat string, long string) (temp float32, wind float32, clouds float32) {
	type Conditions struct {
		Main struct {
			Temperature float32 `json:"feels_like"`
		} `json:"main"`
		Wind struct {
			Gust float32 `json:"gust"`
		} `json:"wind"`
		Clouds struct {
			CloudCover float32 `json:"all"`
		} `json:"clouds"`
	}

	apikey := "fd02e3f101de34b392e1deebe9e19d1e"
	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?lat=%s&lon=%s&appid=%s&units=imperial", lat, long, apikey)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Could not fetch weather data: %s", err)
	}

	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var weather Conditions
	err = json.Unmarshal(body, &weather)
	if err != nil {
		panic(err)
	}
	temp = weather.Main.Temperature
	wind = weather.Wind.Gust
	clouds = weather.Clouds.CloudCover

	return
}
