package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"gonum.org/v1/gonum/mat"

	"api-factorizacion-matriz/pkg/entities"
)

func FactorizarMatriz(c *fiber.Ctx) error {
	var request entities.FactorizacionRequest

	// Parsear la solicitud
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No se pudo parsear la solicitud",
		})
	}
	matriz := request.Matriz

	// Validar que la matriz sea un rectángulo horizontal
	if !tieneMasColumnasQueFilas(matriz) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "La matriz no representa un rectangulo horizontal",
		})
	}

	// Rotar la matriz 90 grados a la derecha
	matrizRotada := rotarMatriz90Grados(matriz)

	// Calcular la factorización QR de la matriz rotada
	matrizQ, matrizR, err := factorizacionQR(matrizRotada)
	if err != nil {
		fmt.Println("Error al calcular la factorización QR:", err)
		return err
	}

	// Analizar la factorización QR
	estadistica, err := analizarFactorizacion(matrizQ, matrizR)
	if err != nil {
		fmt.Println("Error al analizar la factorización QR:", err)
		return err
	}

	resultado := entities.FactorizacionResponse{
		MatrizRotada: matrizRotada,
		MatrizQ:      matrizQ,
		MatrizR:      matrizR,
		Estadistica:  estadistica,
	}

	return c.JSON(resultado)
}

func tieneMasColumnasQueFilas(matriz [][]float64) bool {
	// Obtiene el número de filas
	numFilas := len(matriz)
	// Asume que todas las filas tienen la misma cantidad de columnas y obtiene ese número
	numColumnas := 0
	if numFilas > 0 {
		numColumnas = len(matriz[0])
	}
	// Compara el número de columnas con el número de filas
	return numColumnas > numFilas
}

func rotarMatriz90Grados(matriz [][]float64) [][]float64 {
	if len(matriz) == 0 {
		return nil
	}
	m, n := len(matriz), len(matriz[0])
	// Crear una nueva matriz con dimensiones invertidas
	matrizRotada := make([][]float64, n)
	for i := range matrizRotada {
		matrizRotada[i] = make([]float64, m)
	}

	// Copiar los valores ajustando los índices
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			matrizRotada[j][m-i-1] = matriz[i][j]
		}
	}

	return matrizRotada
}

func factorizacionQR(matriz [][]float64) ([][]float64, [][]float64, error) {
	// Convertir la matriz de entrada a un tipo mat.Dense para su procesamiento
	r, c := len(matriz), len(matriz[0])
	flat := make([]float64, r*c)
	for i, row := range matriz {
		for j, val := range row {
			flat[i*c+j] = val
		}
	}
	A := mat.NewDense(r, c, flat)

	// Realizar la factorización QR
	var qr mat.QR
	qr.Factorize(A)

	// Extraer la matriz Q.
	var Q mat.Dense
	qr.QTo(&Q)

	// Extraer la matriz R.
	var R mat.Dense
	qr.RTo(&R)

	// Convertir las matrices Q y R a slices para su devolución
	matrizQ := denseToSlices(&Q)
	matrizR := denseToSlices(&R)

	return matrizQ, matrizR, nil
}

func denseToSlices(d *mat.Dense) [][]float64 {
	r, c := d.Dims()
	out := make([][]float64, r)
	for i := range out {
		out[i] = make([]float64, c)
		for j := range out[i] {
			out[i][j] = d.At(i, j)
		}
	}
	return out
}

func analizarFactorizacion(matrizQ [][]float64, matrizR [][]float64) (entities.EstadisticaResponse, error) {
	// Crear la estructura de data para enviar al servidor
	data := entities.EstadisticaRequest{
		MatrizQ: matrizQ,
		MatrizR: matrizR,
	}

	// Serializar a JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error serializando los datos:", err)
		return entities.EstadisticaResponse{}, err
	}

	// Crear la petición POST
	// url := "http://localhost:3100/api/estadistica-factorizacion"
	url := "https://cr-api-estadistica-factorizacion-5l45uowjra-uc.a.run.app/api/estadistica-factorizacion"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creando la petición:", err)
		return entities.EstadisticaResponse{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Enviar la petición
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error enviando la petición:", err)
		return entities.EstadisticaResponse{}, err
	}
	defer resp.Body.Close()

	// Leer y manejar la respuesta
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error leyendo la respuesta:", err)
		return entities.EstadisticaResponse{}, err
	}

	// Deserializar el cuerpo de la estadistica
	var estadistica entities.EstadisticaResponse
	err = json.Unmarshal(body, &estadistica)
	if err != nil {
		fmt.Println("Error deserializando la respuesta:", err)
		return entities.EstadisticaResponse{}, err
	}

	return estadistica, nil
}
