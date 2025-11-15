# Ener-gentic
Polyglot AI+Energy Mesh: Python forecasts → Rust optimizes (C++ sim) → Go gateways. gRPC/Arrow batches for scalable kW heartbeats. Docker up, curl predict—interlang niche for grid futurists.

## Mesh & Bench

### /mesh
V3.2 implementation with **Python → Rust(+C++) → Go** architecture. Gateway→Forecaster communication over **Arrow Flight** with graceful fallback to Arrow-HTTP. Supports **cadence switching** (minutes|seconds) for 60/3600-step horizons with gust injection.

**Quickstart:**
```bash
cd mesh
docker compose build && docker compose up

# Flight (preferred)
curl -s -X POST localhost:8080/predict_batch_flight \
  -H 'Content-Type: application/json' \
  -d '{"x":[[4,5,6,7,6,5,6,7,8,7],[3,4.8,7.2,2.1,5.5,6.0,7.8,6.3,5.1,4.9]]}' | jq .

# Arrow-HTTP (baseline)
curl -s -X POST localhost:8080/predict_batch_arrow \
  -H 'Content-Type: application/json' \
  -d '{"x":[[4,5,6,7,6,5,6,7,8,7],[3,4.8,7.2,2.1,5.5,6.0,7.8,6.3,5.1,4.9]]}' | jq .
```

See [mesh/README.md](mesh/README.md) for architecture details and stack configuration.

### /bench
curl-based load generator for comparing JSON vs Arrow-HTTP vs Flight performance. Includes synthetic gust CSV tooling.

**Run benchmarks:**
```bash
cd bench
./bench.sh 10   # 10 series
./bench.sh 50   # 50 series
./bench.sh 100  # 100 series
```

See [bench/README.md](bench/README.md) for usage and expected results.

## License
MIT License - see [LICENSE](LICENSE) for details.

Dependencies (Apache-2.0/BSD): Arrow, gRPC, protobuf.
