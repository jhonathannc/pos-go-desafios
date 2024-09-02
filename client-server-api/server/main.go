package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Cotacao struct {
	USDBRL struct {
		Bid    string `json:"bid"`
		Code   string `json:"code"`
		Codein string `json:"codein"`
		Name   string `json:"name"`
	} `json:"USDBRL"`
}

func main() {
	http.HandleFunc("/cotacao", GetCotacaoHandler)
	http.ListenAndServe(":8080", nil)
}

func GetCotacaoHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/cotacao" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	cotacao, err := FindCotacao()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(cotacao.USDBRL.Bid)
}

func FindCotacao() (*Cotacao, error) {
	// Criando context e definido 200ms para o timeout da api economia
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			fmt.Println("Erro: tempo limite excedido na chamada da API")
		}
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var cotacao Cotacao
	err = json.Unmarshal(body, &cotacao)
	if err != nil {
		return nil, err
	}

	err = SaveDatabase(&cotacao)
	if err != nil {
		return nil, err
	}

	return &cotacao, nil
}

func SaveDatabase(cotacao *Cotacao) error {
	// Criando context e definido 10ms para o timeout para persistir dados no database
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	// Abrir conexão com o banco de dados
	db, err := sql.Open("sqlite3", "./cotacoes.db")
	if err != nil {
		return err
	}
	defer db.Close()

	// Criar tabela se não existir
	_, err = db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS cotacoes (
			id INTEGER PRIMARY KEY,
			bid TEXT,
			code TEXT,
			codein TEXT,
			name TEXT
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			fmt.Println("Erro: tempo limite excedido ao criar tabela")
		}
		return err
	}

	// Inserir cotação no banco de dados
	_, err = db.ExecContext(ctx, `
		INSERT INTO cotacoes (bid, code, codein, name)
		VALUES (?, ?, ?, ?);
	`, cotacao.USDBRL.Bid, cotacao.USDBRL.Code, cotacao.USDBRL.Codein, cotacao.USDBRL.Name)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			fmt.Println("Erro: tempo limite excedido ao inserir cotação no banco de dados")
		}
		return err
	}

	return nil
}
