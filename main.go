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
	temp, wind, clouds, rain, tempUnit := getWeather(lat, long)

	fmt.Printf("It is currently %g%s with gusts to %gmph. There is %g%% cloud cover and you can expect %gin of rain.", temp, tempUnit, wind, clouds, rain)
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

func getWeather(lat string, long string) (temp float32, wind float32, clouds float32, rain float32, tempUnit string) {
	type Weather struct {
		Unit struct {
			FeelsLike string `json:"apparent_temperature"`
		} `json:"current_units"`
		Current struct {
			FeelsLike     float32 `json:"apparent_temperature"`
			Precipitation float32 `json:"precipitation"`
			Clouds        float32 `json:"cloud_cover"`
			Wind          float32 `json:"wind_gusts_10m"`
		} `json:"current"`
	}

	baseurl := "https://api.open-meteo.com/v1/forecast?"
	params := fmt.Sprintf("latitude=%s&longitude=%s", lat, long)
	data := "current=temperature_2m,apparent_temperature,precipitation,cloud_cover,wind_gusts_10m"
	units := "temperature_unit=fahrenheit&wind_speed_unit=mph&precipitation_unit=inch"
	url := fmt.Sprintf("%s%s&%s&%s", baseurl, params, data, units)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Could not fetch weather data: %s", err)
	}

	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var weather Weather
	err = json.Unmarshal(body, &weather)
	if err != nil {
		panic(err)
	}

	temp = weather.Current.FeelsLike
	wind = weather.Current.Wind
	clouds = weather.Current.Clouds
	rain = weather.Current.Precipitation
	tempUnit = weather.Unit.FeelsLike
	return
}
