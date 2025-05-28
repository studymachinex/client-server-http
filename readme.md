# Desafio Go: Client-Server HTTP com Context, Banco de Dados e Manipulação de Arquivos

## Descrição

Este projeto consiste em dois sistemas escritos em Go: um servidor (`server.go`) e um cliente (`client.go`). O objetivo é consumir a cotação do dólar em tempo real, persistir os dados em um banco SQLite e manipular arquivos, utilizando contextos para controle de timeout.

---

## Requisitos

### 1. Servidor (`server.go`)

- Expõe o endpoint HTTP `/cotacao` na porta `8080`.
- Ao receber uma requisição, consome a API pública:  
  `https://economia.awesomeapi.com.br/json/last/USD-BRL`
- Retorna ao cliente apenas o valor do campo `bid` em formato JSON.
- Utiliza o package `context` para:
  - Timeout de 200ms para chamada da API de cotação.
  - Timeout de 10ms para persistir a cotação no banco SQLite.
- Registra cada cotação recebida no banco de dados SQLite.
- Loga erro caso algum contexto ultrapasse o tempo limite.

### 2. Cliente (`client.go`)

- Realiza uma requisição HTTP para o endpoint `/cotacao` do servidor.
- Utiliza o package `context` com timeout de 300ms para aguardar a resposta.
- Recebe apenas o valor do campo `bid`.
- Salva o valor recebido em um arquivo `cotacao.txt` no formato:  
  `Dólar: {valor}`
- Loga erro caso o contexto ultrapasse o tempo limite.

---

## Como Executar

### 1. Inicie o Servidor

No diretório `server/`:

```powershell
cd server
go run main.go
```

O servidor estará disponível em `http://localhost:8080/cotacao`.

### 2. Execute o Cliente

No diretório `client/`:

```powershell
cd client
go run main.go
```

O cliente irá:
- Solicitar a cotação ao servidor.
- Salvar o valor em `cotacao.txt`.

---

## Observações

- Certifique-se de que o servidor esteja rodando antes de executar o cliente.
- O banco de dados SQLite será criado automaticamente pelo servidor.
- Todos os contextos possuem timeouts rígidos, conforme especificado nos requisitos.
- Em caso de timeout, mensagens de erro serão exibidas nos logs.

---

## Tecnologias Utilizadas

- [Go](https://golang.org/)
- [SQLite](https://www.sqlite.org/index.html)
- [HTTP](https://pkg.go.dev/net/http)
- [Context](https://pkg.go.dev/context)

---

## Autor

Desafio proposto por [Full Cycle](https://fullcycle.com.br/).