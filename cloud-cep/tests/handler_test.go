package tests

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"jhonathannc/pos-go-desafios/cloud-cep/handler"
)

func TestWeatherHandler_InvalidCEP(t *testing.T) {
	req := httptest.NewRequest("GET", "/weather?cep=123", nil)
	rr := httptest.NewRecorder()

	handler.WeatherByCEPHandler(rr, req)

	if rr.Code != http.StatusUnprocessableEntity {
		t.Errorf("Expected status 422, got %d", rr.Code)
	}
}

func TestWeatherHandler_NotFound(t *testing.T) {
	req := httptest.NewRequest("GET", "/weather?cep=00000000", nil)
	rr := httptest.NewRecorder()

	handler.WeatherByCEPHandler(rr, req)

	if rr.Code != http.StatusUnprocessableEntity {
		t.Errorf("Expected status 404, got %d", rr.Code)
	}
}

func TestWeatherHandler_Success(t *testing.T) {
	os.Setenv("WEATHER_API_KEY", "test") // Use mock in real case

	req := httptest.NewRequest("GET", "/weather?cep=01001000", nil)
	rr := httptest.NewRecorder()

	handler.WeatherByCEPHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}
}
