package entities

type FactorizacionRequest struct {
	Matriz [][]float64 `json:"matriz"`
}

type FactorizacionResponse struct {
	MatrizRotada [][]float64         `json:"matrizRotada"`
	MatrizQ      [][]float64         `json:"matrizQ"`
	MatrizR      [][]float64         `json:"matrizR"`
	Estadistica  EstadisticaResponse `json:"estadistica"`
}
