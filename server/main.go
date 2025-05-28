package main

import (
	"client-server-http/shared"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func main() {
	log.Println("Iniciando servidor...")

	var err error
	db, err = sql.Open("mysql", "root:root@tcp(localhost:3306)/goexpert")
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}
	defer db.Close()

	err = createTableIfNotExists()
	if err != nil {
		log.Fatalf("Erro ao criar tabela: %v", err)
	}

	http.HandleFunc("/cotacao", cotacaoHandler)
	log.Println("Servidor escutando na porta 8080...")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Erro ao iniciar servidor: %v", err)
	}
}



func cotacaoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		log.Printf("Método %s não permitido", r.Method)
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}
	apiUrl := "https://economia.awesomeapi.com.br/json/last/USD-BRL"

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiUrl, nil)
	if err != nil {
		log.Printf("Erro ao criar request: %v", err)
		http.Error(w, "Erro interno", http.StatusInternalServerError)
		return
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			log.Println("Timeout: tempo de execução insuficiente para buscar dados da API")
		}
		log.Printf("Erro ao buscar dados da API: %v", err)
		http.Error(w, "Erro ao buscar dados", http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	jsonResponse, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("Erro ao ler corpo da resposta: %v", err)
		http.Error(w, "Erro ao ler resposta", http.StatusInternalServerError)
		return
	}

	var response shared.DolarApiResponse
	err = json.Unmarshal(jsonResponse, &response)
	if err != nil {
		log.Printf("Erro ao fazer unmarshal da resposta: %v", err)
		http.Error(w, "Erro ao processar resposta", http.StatusInternalServerError)
		return
	}

	jsonResp, err := json.Marshal(response.USDBRL)
	if err != nil {
		log.Printf("Erro ao serializar resposta: %v", err)
		http.Error(w, "Erro ao processar resposta", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)

	go func(valor string) {
		if err := persistirCotacao(valor); err != nil {
			log.Printf("Erro ao persistir cotação: %v", err)
		}
	}(response.USDBRL.Bid)
}

func persistirCotacao(valorDolar string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			log.Println("Transaction Timeout")
			return fmt.Errorf("transaction timeout")
		}
		return fmt.Errorf("erro ao iniciar transação: %w", err)
	}

	stmt, err := tx.PrepareContext(ctx, "INSERT INTO cotacoes (dolar) VALUES (?)")
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			log.Println("Prepare Timeout")
			return fmt.Errorf("prepare timeout")
		}
		tx.Rollback()
		return fmt.Errorf("erro ao preparar statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, valorDolar)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			log.Println("Execute Timeout")
			tx.Rollback()
			return fmt.Errorf("execute timeout")
		}
		tx.Rollback()
		return fmt.Errorf("erro ao executar statement: %w", err)
	}

	if err = tx.Commit(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			log.Println("Commit Timeout")
			return fmt.Errorf("commit Timeout")
		}
		tx.Rollback()
		return fmt.Errorf("erro ao commitar transação: %w", err)
	}

	log.Printf("Cotação persistida com sucesso: %s", valorDolar)
	return nil
}

func createTableIfNotExists() error {
	tableQuery := `CREATE TABLE IF NOT EXISTS cotacoes (
		id INT AUTO_INCREMENT PRIMARY KEY,
		dolar VARCHAR(20) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`
	_, err := db.Exec(tableQuery)
	if err != nil {
		return fmt.Errorf("erro ao criar tabela: %w", err)
	}
	return nil
}
