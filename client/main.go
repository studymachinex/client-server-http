package main

import (
	"client-server-http/shared"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*300)
	defer cancel()

	getCotacaoDolar(ctx)
	log.Println("Requisição concluída")
}

func getCotacaoDolar(ctx context.Context) (string, error) {
	apiUrl := "https://economia.awesomeapi.com.br/json/last/USD-BRL"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiUrl, nil)
	if err != nil {
		log.Fatalf("Erro ao fazer requisição: %v", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("erro ao obter cotação: %s", resp.Status)
	}

	var response shared.DolarApiResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", err
	}

	return response.USDBRL.Bid, nil
}
