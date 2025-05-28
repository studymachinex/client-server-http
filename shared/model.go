package shared

type DolarApiResponse struct {
	USDBRL struct {
		Bid string `json:"bid"`
	} `json:"USDBRL"`
}

// Nova struct para resposta simples da API
// Apenas o campo bid

type SimpleBidResponse struct {
	Bid string `json:"bid"`
}