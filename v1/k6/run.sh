#!/usr/bin/env bash
# Chạy distributed K6: chia TARGET_RPS đều cho N process song song
# Usage: ./run.sh [total_rps] [instances] [duration]
#   ./run.sh 100000 4 60s   → 4 process × 25000 RPS = 100k TPS

set -euo pipefail

TOTAL_RPS=${1:-100000}
INSTANCES=${2:-4}
DURATION=${3:-60s}
BASE_URL=${BASE_URL:-http://localhost:8000}

RPS_PER_INSTANCE=$(( TOTAL_RPS / INSTANCES ))

echo "========================================"
echo "Total RPS   : ${TOTAL_RPS}"
echo "Instances   : ${INSTANCES}"
echo "RPS/instance: ${RPS_PER_INSTANCE}"
echo "Duration    : ${DURATION}"
echo "Target      : ${BASE_URL}"
echo "========================================"

pids=()
for i in $(seq 1 "$INSTANCES"); do
  echo "[instance $i] starting → ${RPS_PER_INSTANCE} RPS"
  BASE_URL="$BASE_URL" \
  TARGET_RPS="$RPS_PER_INSTANCE" \
  DURATION="$DURATION" \
  k6 run --tag "instance=$i" load-test.js &
  pids+=($!)
done

# Chờ tất cả xong
for pid in "${pids[@]}"; do
  wait "$pid"
done

echo "All instances done."
