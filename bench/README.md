# Ener-gentic Benchmark Suite

Load generator and synthetic data tooling for comparing JSON, Arrow-HTTP, and Arrow Flight performance.

## Usage

```bash
# Start the mesh
cd ../mesh
docker compose up

# Run benchmark with 10 series (default)
./bench.sh

# Run with 50 series
./bench.sh 50

# Run with 100 series
./bench.sh 100
```

## What it measures

- **JSON Single**: Individual `/predict` requests (baseline)
- **Arrow-HTTP Batch**: Batch `/predict_batch_arrow` with Arrow IPC over HTTP
- **Arrow Flight Batch**: Batch `/predict_batch_flight` with Arrow Flight (preferred)

## Expected results

Flight should show significant throughput improvements for batch operations, especially with larger series counts (50+).

## Synthetic data

The script generates random time series data for testing. For real gust patterns, use the CSV tooling in `/mesh/forecaster/data/`.
