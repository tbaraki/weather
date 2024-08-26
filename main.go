package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Location struct {
	Latitude  float32 `json:"lat"`
	Longitude float32 `json:"lon"`
}

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

func main() {
	location := getLocation()
	weather := getWeather(location)

	fmt.Printf("It is currently %g%s with gusts to %gmph. There is %g%% cloud cover and you can expect %gin of rain.",
		weather.Current.FeelsLike, weather.Unit.FeelsLike, weather.Current.Wind, weather.Current.Clouds, weather.Current.Precipitation)
}

func getLocation() (location Location) {
	resp, err := http.Get("http://ip-api.com/json?fields=lat,lon")
	if err != nil {
		fmt.Printf("error getting location: %s", err)
	}

	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	// var location Location
	err = json.Unmarshal(body, &location)
	if err != nil {
		panic(err)
	}
	return
}

func getWeather(location Location) (weather Weather) {
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
	params := fmt.Sprintf("latitude=%g&longitude=%g", location.Latitude, location.Longitude)
	data := "current=temperature_2m,apparent_temperature,precipitation,cloud_cover,wind_gusts_10m"
	units := "temperature_unit=fahrenheit&wind_speed_unit=mph&precipitation_unit=inch"
	url := fmt.Sprintf("%s%s&%s&%s", baseurl, params, data, units)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Could not fetch weather data: %s", err)
	}

	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	err = json.Unmarshal(body, &weather)
	if err != nil {
		panic(err)
	}
	return
}
