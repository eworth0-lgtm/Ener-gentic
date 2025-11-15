package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type PredictRequest struct {
	X []float32 `json:"x"`
}

type BatchPredictRequest struct {
	X [][]float32 `json:"x"`
}

type PredictResponse struct {
	Forecast []float32 `json:"forecast"`
	Dispatch []float32 `json:"dispatch"`
}

type BatchPredictResponse struct {
	Forecasts [][]float32 `json:"forecasts"`
	Dispatches [][]float32 `json:"dispatches"`
}

var (
	forecasterAddr      string
	forecasterFlightAddr string
	optimizerAddr       string
)

func main() {
	forecasterAddr = getEnv("FORECASTER_ADDR", "localhost:50051")
	forecasterFlightAddr = getEnv("FORECASTER_FLIGHT_ADDR", "localhost:8815")
	optimizerAddr = getEnv("OPTIMIZER_ADDR", "localhost:50052")

	http.HandleFunc("/predict", handlePredict)
	http.HandleFunc("/predict_batch_flight", handlePredictBatchFlight)
	http.HandleFunc("/predict_batch_arrow", handlePredictBatchArrow)

	log.Println("Go Gateway running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handlePredict(w http.ResponseWriter, r *http.Request) {
	var req PredictRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Call forecaster via gRPC
	forecast, err := callForecaster(req.X)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Call optimizer via gRPC
	dispatch, err := callOptimizer(forecast)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := PredictResponse{
		Forecast: forecast,
		Dispatch: dispatch,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func handlePredictBatchFlight(w http.ResponseWriter, r *http.Request) {
	var req BatchPredictRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Call forecaster via Arrow Flight (preferred)
	forecasts, err := callForecasterBatchFlight(req.X)
	if err != nil {
		// Fallback to Arrow-HTTP
		log.Printf("Flight failed, falling back to Arrow-HTTP: %v", err)
		forecasts, err = callForecasterBatchArrow(req.X)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Call optimizer batch via gRPC
	dispatches, err := callOptimizerBatch(forecasts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := BatchPredictResponse{
		Forecasts:  forecasts,
		Dispatches: dispatches,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func handlePredictBatchArrow(w http.ResponseWriter, r *http.Request) {
	var req BatchPredictRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Call forecaster via Arrow-HTTP
	forecasts, err := callForecasterBatchArrow(req.X)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Call optimizer batch via gRPC
	dispatches, err := callOptimizerBatch(forecasts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := BatchPredictResponse{
		Forecasts:  forecasts,
		Dispatches: dispatches,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func callForecaster(x []float32) ([]float32, error) {
	conn, err := grpc.Dial(forecasterAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// Simplified: would use generated proto client
	// For now, return mock data
	forecast := make([]float32, 60)
	for i := range forecast {
		forecast[i] = x[len(x)-1] + float32(i)*0.1
	}
	return forecast, nil
}

func callForecasterBatchFlight(series [][]float32) ([][]float32, error) {
	// Arrow Flight implementation
	// For now, fallback to Arrow-HTTP
	return callForecasterBatchArrow(series)
}

func callForecasterBatchArrow(series [][]float32) ([][]float32, error) {
	// Arrow-HTTP implementation
	forecasts := make([][]float32, len(series))
	for i, x := range series {
		forecast, err := callForecaster(x)
		if err != nil {
			return nil, err
		}
		forecasts[i] = forecast
	}
	return forecasts, nil
}

func callOptimizer(forecast []float32) ([]float32, error) {
	conn, err := grpc.Dial(optimizerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// Simplified: would use generated proto client
	dispatch := make([]float32, len(forecast))
	for i, v := range forecast {
		dispatch[i] = v * 0.95 // Simple optimization
	}
	return dispatch, nil
}

func callOptimizerBatch(forecasts [][]float32) ([][]float32, error) {
	dispatches := make([][]float32, len(forecasts))
	for i, forecast := range forecasts {
		dispatch, err := callOptimizer(forecast)
		if err != nil {
			return nil, err
		}
		dispatches[i] = dispatch
	}
	return dispatches, nil
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
