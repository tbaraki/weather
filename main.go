package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

type Location struct {
	Latitude  float32 `json:"lat"`
	Longitude float32 `json:"lon"`
}

type ZipLocation struct {
	Places []Place `json:"places"`
}

type Place struct {
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
	City      string `json:"place name"`
}

type Weather struct {
	Unit struct {
		FeelsLike string `json:"temperature_2m"`
	} `json:"current_units"`
	Current struct {
		FeelsLike     float32 `json:"temperature_2m"`
		Precipitation float32 `json:"precipitation"`
		Clouds        float32 `json:"cloud_cover"`
		Wind          float32 `json:"wind_gusts_10m"`
	} `json:"current"`
}

func main() {
	zip := flag.String("z", "", "5-digit US zipcode. Example: 'weather -z 90210'")
	flag.Parse()

	var location Location

	if *zip != "" {
		ziplocation := getLocationByZip(*zip)
		latFloat, _ := strconv.ParseFloat(ziplocation.Places[0].Latitude, 32)
		longFloat, _ := strconv.ParseFloat(ziplocation.Places[0].Longitude, 32)
		location.Latitude = float32(latFloat)
		location.Longitude = float32(longFloat)
		fmt.Printf("Fetching weather for %s...\n", ziplocation.Places[0].City)
	} else {
		iplocation := getLocationByIp()
		location.Latitude = iplocation.Latitude
		location.Longitude = iplocation.Longitude
	}

	weather := getWeather(location)

	fmt.Printf("It is currently %g%s with gusts to %gmph.\n",
		weather.Current.FeelsLike,
		weather.Unit.FeelsLike,
		weather.Current.Wind)

	fmt.Printf("There is %g%% cloud cover. You can expect %gin of rain.",
		weather.Current.Clouds,
		weather.Current.Precipitation)
}

func getLocationByIp() (iplocation Location) {
	resp, err := http.Get("http://ip-api.com/json?fields=lat,lon")
	if err != nil {
		fmt.Printf("error getting location: %s", err)
	}

	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(body, &iplocation)
	if err != nil {
		panic(err)
	}
	return
}

func getLocationByZip(zip string) (ziplocation ZipLocation) {
	baseurl := "https://api.zippopotam.us/us/"
	url := fmt.Sprint(baseurl, zip)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("error getting location: %s", err)
	}

	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(body, &ziplocation)
	if err != nil {
		panic(err)
	}
	return
}

func getWeather(location Location) (weather Weather) {
	baseurl := "https://api.open-meteo.com/v1/forecast?"
	params := fmt.Sprintf("latitude=%g&longitude=%g", location.Latitude, location.Longitude)
	data := "current=temperature_2m,precipitation,cloud_cover,wind_gusts_10m"
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
