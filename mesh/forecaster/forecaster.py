#!/usr/bin/env python3
import os
import time
import numpy as np
import grpc
from concurrent import futures
import pyarrow as pa
import pyarrow.flight as flight
from flask import Flask, request, jsonify

# Import generated proto
import core_pb2
import core_pb2_grpc

CADENCE = os.getenv("CADENCE", "minutes")
HORIZON = 3600 if CADENCE == "seconds" else 60

app = Flask(__name__)


class ForecasterServicer(core_pb2_grpc.ForecasterServicer):
    def Predict(self, request, context):
        x = np.array(request.values)
        forecast = simple_forecast(x, HORIZON)
        return core_pb2.Forecast(values=forecast.tolist())

    def PredictBatch(self, request, context):
        forecasts = []
        for series in request.series:
            x = np.array(series.values)
            forecast = simple_forecast(x, HORIZON)
            forecasts.append(core_pb2.Forecast(values=forecast.tolist()))
        return core_pb2.BatchForecast(forecasts=forecasts)


class FlightServer(flight.FlightServerBase):
    def __init__(self, location, **kwargs):
        super().__init__(location, **kwargs)

    def do_action(self, context, action):
        if action.type == "predict_batch":
            # Deserialize Arrow IPC from body
            reader = pa.ipc.open_stream(action.body)
            table = reader.read_all()
            
            # Process each series
            results = []
            for batch in table.to_batches():
                for row in batch.to_pylist():
                    series_data = row['series']
                    forecast = simple_forecast(np.array(series_data), HORIZON)
                    results.append(forecast.tolist())
            
            # Return Arrow table
            result_table = pa.table({'forecast': results})
            sink = pa.BufferOutputStream()
            writer = pa.ipc.new_stream(sink, result_table.schema)
            writer.write_table(result_table)
            writer.close()
            
            yield flight.Result(sink.getvalue())


def simple_forecast(x, horizon):
    """Simple moving average forecast with optional gust injection"""
    mean = np.mean(x[-10:]) if len(x) >= 10 else np.mean(x)
    std = np.std(x[-10:]) if len(x) >= 10 else np.std(x)
    
    forecast = np.full(horizon, mean)
    
    # Add micro-gusts for seconds cadence
    if CADENCE == "seconds":
        gust_indices = np.random.choice(horizon, size=max(1, horizon // 100), replace=False)
        forecast[gust_indices] += np.random.normal(0, std * 2, size=len(gust_indices))
    
    # Add noise
    forecast += np.random.normal(0, std * 0.1, size=horizon)
    
    return forecast


@app.route('/predict', methods=['POST'])
def predict_http():
    data = request.json
    x = np.array(data['x'])
    forecast = simple_forecast(x, HORIZON)
    return jsonify({'forecast': forecast.tolist()})


@app.route('/predict_batch_arrow', methods=['POST'])
def predict_batch_arrow_http():
    """Arrow-HTTP fallback endpoint"""
    data = request.json
    results = []
    for series in data['x']:
        x = np.array(series)
        forecast = simple_forecast(x, HORIZON)
        results.append(forecast.tolist())
    return jsonify({'forecasts': results})


def serve():
    # Start gRPC server
    grpc_server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    core_pb2_grpc.add_ForecasterServicer_to_server(ForecasterServicer(), grpc_server)
    grpc_server.add_insecure_port('[::]:50051')
    grpc_server.start()
    print(f"gRPC Forecaster running on :50051 (cadence={CADENCE}, horizon={HORIZON})")

    # Start Arrow Flight server
    location = flight.Location.for_grpc_tcp("0.0.0.0", 8815)
    flight_server = FlightServer(location)
    flight_server.serve()
    print(f"Arrow Flight server running on :8815")

    # Start Flask HTTP server
    app.run(host='0.0.0.0', port=5000, threaded=True)


if __name__ == '__main__':
    serve()
