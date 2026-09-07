#!/bin/bash
set -e

# Benchmark script for JSON vs Arrow-HTTP vs Flight
# Usage: ./bench.sh [series_count]

SERIES_COUNT=${1:-10}
GATEWAY_URL=${GATEWAY_URL:-http://localhost:8080}

echo "=== Ener-gentic Mesh Benchmark ==="
echo "Series count: $SERIES_COUNT"
echo "Gateway: $GATEWAY_URL"
echo ""

# Generate synthetic data
generate_series() {
    local count=$1
    python3 -c "
import json
import random
series = [[random.uniform(3, 8) for _ in range(20)] for _ in range($count)]
print(json.dumps({'x': series}))
"
}

# Benchmark JSON single
echo "--- JSON Single Predict ---"
PAYLOAD='{"x":[4,5,6,7,6,5,5,6,8,7,6,6,7,7,6,5,5,6,7,8]}'
time for i in $(seq 1 $SERIES_COUNT); do
    curl -s -X POST $GATEWAY_URL/predict \
        -H 'Content-Type: application/json' \
        -d "$PAYLOAD" > /dev/null
done
echo ""

# Benchmark Arrow-HTTP batch
echo "--- Arrow-HTTP Batch ---"
BATCH_PAYLOAD=$(generate_series $SERIES_COUNT)
time curl -s -X POST $GATEWAY_URL/predict_batch_arrow \
    -H 'Content-Type: application/json' \
    -d "$BATCH_PAYLOAD" > /dev/null
echo ""

# Benchmark Flight batch
echo "--- Arrow Flight Batch ---"
time curl -s -X POST $GATEWAY_URL/predict_batch_flight \
    -H 'Content-Type: application/json' \
    -d "$BATCH_PAYLOAD" > /dev/null
echo ""

echo "=== Benchmark Complete ==="
echo ""
echo "Summary:"
echo "- JSON: $SERIES_COUNT individual requests"
echo "- Arrow-HTTP: 1 batch request ($SERIES_COUNT series)"
echo "- Flight: 1 batch request ($SERIES_COUNT series, preferred)"
