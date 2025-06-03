// services/weather.go
package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

type WeatherAPIResponse struct {
	Current struct {
		TempC float64 `json:"temp_c"`
	} `json:"current"`
}

// GetTemperatureByCityWithClient busca a temperatura usando um cliente HTTP customizado (facilita testes)
func GetTemperatureByCityWithClient(client *http.Client, city string) (float64, error) {
	apiKey := os.Getenv("WEATHER_API_KEY")
	if apiKey == "" {
		return 0, fmt.Errorf("weather api key not set")
	}

	escapedCity := url.QueryEscape(city)
	url := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%s&q=%s", apiKey, escapedCity)

	resp, err := client.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("weather api error: status code %d", resp.StatusCode)
	}

	var data WeatherAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, err
	}

	return data.Current.TempC, nil
}

// função helper
func GetTemperatureByCity(city string) (float64, error) {
	return GetTemperatureByCityWithClient(http.DefaultClient, city)
}
