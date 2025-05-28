package main

import (
	"client-server-http/shared"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	resp, err := getCotacaoDolar()
	if err != nil {
		log.Fatal(err)
	}
	salvarCotacaoEmArquivo(resp)
	log.Printf("Cotação do dólar: R$ %s", resp)

}

func salvarCotacaoEmArquivo(resp string) {
	file, err := os.Create("cotacao.txt")
	if err != nil {
		log.Fatalf("Erro ao criar arquivo: %v", err)
	}
	defer file.Close()

	_, err = file.WriteString(resp)
	if err != nil {
		log.Fatalf("Erro ao escrever no arquivo: %v", err)
	}
	log.Println("Cotação salva em cotacao.txt")
}

func getCotacaoDolar() (string, error) {
	apiUrl := "http://localhost:8080/cotacao"
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*300)
	defer cancel()
	
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, apiUrl, nil)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			log.Println("Timeout: tempo de execução insuficiente para buscar dados da API")
		}
		return "erro ao realizar a requisicao: ", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("erro ao obter cotação: %s", resp.Status)
	}

	byteData, err := io.ReadAll(resp.Body)
	if err != nil {
		return "erro ao ler resposta: ", err
	}
	var c shared.SimpleBidResponse
	err = json.Unmarshal(byteData, &c)
	if err != nil {
		return "erro ao decodificar resposta: ", err
	}
	return fmt.Sprintf("Dólar: %s", c.Bid), nil
}
