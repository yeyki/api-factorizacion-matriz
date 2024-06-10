package entities

type EstadisticaRequest struct {
	MatrizQ [][]float64 `json:"matrizQ"`
	MatrizR [][]float64 `json:"matrizR"`
}

type EstadisticaResponse struct {
	ValorMaximo    float64 `json:"valorMaximo"`
	ValorMinimo    float64 `json:"valorMinimo"`
	Promedio       float64 `json:"promedio"`
	SumaTotal      float64 `json:"sumaTotal"`
	MatrizDiagonal string  `json:"matrizDiagonal"`
}
