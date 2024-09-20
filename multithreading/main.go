package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type ViaCepAddress struct {
	Cep        string `json:"cep"`
	Logradouro string `json:"logradouro"`
	Bairro     string `json:"bairro"`
	Localidade string `json:"localidade"`
	UF         string `json:"uf"`
}

type BrasilApiAddress struct {
	Cep          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
}

type ApiResponse struct {
	APIName string
	Data    interface{}
}

func main() {
	cep := "13409120"

	ch := make(chan ApiResponse)

	endpoint := "https://brasilapi.com.br/api/cep/v1/" + cep
	go getCepFromApi("brasilapi", endpoint, ch)

	endpoint = "http://viacep.com.br/ws/" + cep + "/json/"
	go getCepFromApi("viacep", endpoint, ch)

	select {
	case res := <-ch:
		fmt.Printf("Resposta mais rapida: %s\n", res.APIName)
		fmt.Println("Endereço:")
		b, _ := json.MarshalIndent(res.Data, "", "  ")
		fmt.Println(string(b))
	case <-time.After(1 * time.Second):
		fmt.Println("timeout 1s atingido.")
	}
}

func getCepFromApi(apiName, endpoint string, resultChan chan ApiResponse) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		panic(err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			fmt.Println("Erro: Tempo limit excedido, a requisição levou mais de 1s")
		}
		panic(err)
	}
	defer res.Body.Close()

	var data interface{}
	if apiName == "brasilapi" {
		data = &BrasilApiAddress{}
	} else if apiName == "viacep" {
		data = &ViaCepAddress{}
	}

	err = json.NewDecoder(res.Body).Decode(data)
	if err != nil {
		panic(err)
	}

	resultChan <- ApiResponse{APIName: apiName, Data: data}
}
