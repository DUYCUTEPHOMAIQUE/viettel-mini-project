# Load Test Results — 1 Pod

**Config:** 1 pod Node B · HAProxy L4 · `cpus: 0.5` · `memory: 256M` · duration 30s/run

## Kết quả

| Run | Target RPS | Actual RPS | Avg     | P95      | P99      | Error Rate | Dropped |
|:---:|:----------:|:----------:|:-------:|:--------:|:--------:|:----------:|:-------:|
|  1  |   4 500    |   4 498    | 379 µs  |  757 µs  | 2.32 ms  |   0.00%    |   39    |
|  2  |   4 500    |   4 499    | 355 µs  |  647 µs  | 1.64 ms  |   0.00%    |    4    |
|  3  |   4 500    |   4 497    | 416 µs  |  941 µs  | 2.36 ms  |   0.00%    |   58    |
|  4  |   4 700    |   4 696    | 447 µs  |  837 µs  | 3.06 ms  |   0.00%    |  112    |
|  5  |   4 700    |   4 697    | 449 µs  |  811 µs  | 3.33 ms  |   0.00%    |   81    |
|  6  |   4 700    |   4 694    | 597 µs  | 1.11 ms  | 10.6 ms  |   0.00%    |  171    |
|  7  |   5 000    |   4 986    | 1.11 ms | 4.09 ms  | 22.18 ms |   0.00%    |  393    |
|  8  |   5 000    |   4 972    | 5.15 ms | 30.54 ms | 43.58 ms |   0.00%    |  814    |
|  9  |   5 000    |   4 972    | 3.59 ms | 25.36 ms | 42.81 ms |   0.00%    |  822    |

## Nhận xét

| Target RPS | Trạng thái | Ghi chú |
|:----------:|:----------:|---------|
| 4 500      | ✅ Đạt KPI | P95 < 1ms, P99 < 3ms, ổn định qua 3 lần |
| 4 700      | ⚠️ Cận ngưỡng | P95 < 1.2ms nhưng P99 dao động mạnh ở run 3 (10.6ms) |
| 5 000      | ❌ Vượt KPI | P95 lên 25–30ms, P99 gần chạm 50ms, dropped tăng mạnh |

> **Saturation point** của 1 pod (0.5 CPU): ~**4 700–4 800 RPS**

## Scripts

**1 pod:**
```bash
docker compose -f docker-compose.1pod.yml up -d
cd k6 && TARGET_RPS=4500 k6 run load-test.js
```

**2 pods:**
```bash
docker compose -f docker-compose.1pod.yml down
docker compose -f docker-compose.2pod.yml up -d
cd k6 && TARGET_RPS=9000 k6 run load-test.js
```

**3 pods:**
```bash
docker compose -f docker-compose.2pod.yml down
docker compose up -d
cd k6 && TARGET_RPS=13500 k6 run load-test.js
```

---

_Test run: 2026-06-25 03:17:15_

# Terminal 2 — mở trước khi chạy k6                      
watch -n 1 'docker stats --no-stream --format "table {{.Name}}\t{{.CPUPerc}}\t{{.MemUsage}}" | grep -E "NAME|node-b|haproxy"'

# coi cpu usage của k6
watch -n 1 'ps -p $(pgrep k6) -o pid,pcpu,pmem,vsz,rss,comm'