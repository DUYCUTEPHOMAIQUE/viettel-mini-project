# Load Test Results — 2 Pods

**Config:** 2 pods Node B · HAProxy L4 · `cpus: 0.5` · `memory: 256M` · duration 30s/run

## Kết quả

| Run | Target RPS | Actual RPS | Avg      | P95      | P99      | Error Rate | Dropped |
|:---:|:----------:|:----------:|:--------:|:--------:|:--------:|:----------:|:-------:|
|  1  |   9 000    |   8 970    |  2.35 ms | 12.46 ms | 23.94 ms |   0.00%    |   866   |
|  2  |   9 000    |   8 992    |  1.81 ms |  8.57 ms | 21.61 ms |   0.00%    |   212   |
|  3  |   9 000    |   8 978    |  1.58 ms |  7.10 ms | 20.96 ms |   0.00%    |   517   |
|  4  |   9 400    |   9 369    |  4.35 ms | 27.80 ms | 48.86 ms |   0.00%    |   787   |
|  5  |   9 400    |   9 289    |  5.12 ms | 27.28 ms | 41.93 ms |   0.00%    | 3 296   |
|  6  |   9 400    |   9 326    |  4.09 ms | 22.90 ms | 43.36 ms |   0.00%    | 2 115   |
|  7  |  10 000    |   9 863    | 11.67 ms | 54.80 ms | 63.83 ms |   0.00%    | 3 578   |
|  8  |  10 000    |   9 714    | 15.44 ms | 59.74 ms | 69.60 ms |   0.00%    | 8 303   |
|  9  |  10 000    |   9 877    | 15.47 ms | 60.40 ms | 68.23 ms |   0.00%    | 3 396   |

## Nhận xét

| Target RPS | Trạng thái | Ghi chú |
|:----------:|:----------:|---------|
| 9 000      | ⚠️ Cận ngưỡng | P95 7–12ms đạt KPI, nhưng P99 ~21–24ms — dao động cao |
| 9 400      | ❌ Vượt KPI  | P95 23–28ms, P99 gần chạm 50ms, dropped tăng mạnh |
| 10 000     | ❌ Vượt KPI  | P95 và P99 vượt ngưỡng, dropped 3k–8k/30s |

> **Saturation point** của 2 pods (0.5 CPU/pod): ~**9 000–9 200 RPS**

## So sánh với 1 Pod

| Cấu hình | Saturation point | Tỉ lệ scale |
|:--------:|:----------------:|:-----------:|
| 1 pod    | ~4 700 RPS       | 1×          |
| 2 pods   | ~9 000 RPS       | ~1.9×       |

> Scale ngang 2× pod → throughput tăng gần **2× (1.9×)** — linear scaling hoạt động tốt.

## Scripts

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

_Test run: 2026-06-25 03:42:31_
