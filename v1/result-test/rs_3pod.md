# Load Test Results — 3 Pods

**Config:** 3 pods Node B · HAProxy L4 · `cpus: 0.5` · `memory: 256M` · duration 30s/run

## Kết quả

| Run | Target RPS | Actual RPS | Avg      | P95      | P99      | Error Rate | Dropped  |
|:---:|:----------:|:----------:|:--------:|:--------:|:--------:|:----------:|:--------:|
|  1  |   13 500   |   13 290   |  2.27 ms | 11.60 ms | 24.41 ms |   0.00%    |   6 266  |
|  2  |   13 500   |   13 392   |  2.16 ms | 11.07 ms | 22.41 ms |   0.00%    |   3 194  |
|  3  |   13 500   |   13 180   |  3.62 ms | 18.51 ms | 30.73 ms |   0.00%    |   9 557  |
|  4  |   14 000   |   13 558   |  4.74 ms | 22.80 ms | 36.44 ms |   0.00%    |  13 106  |
|  5  |   14 000   |   13 233   |  6.05 ms | 27.23 ms | 40.45 ms |   0.00%    |  22 495  |
|  6  |   14 000   |   12 424   | 10.17 ms | 39.36 ms | 49.82 ms |   0.00%    |  47 163  |
|  7  |   14 500   |   12 902   | 10.59 ms | 41.21 ms | 52.95 ms |   0.00%    |  47 516  |
|  8  |   14 500   |   13 629   |  7.18 ms | 30.94 ms | 42.93 ms |   0.00%    |  25 800  |
|  9  |   14 500   |   13 868   |  6.22 ms | 26.57 ms | 41.36 ms |   0.00%    |  18 870  |

## Nhận xét

| Target RPS | Trạng thái | Ghi chú |
|:----------:|:----------:|---------|
| 13 500     | ⚠️ Cận ngưỡng | P95 11–18ms đạt KPI, P99 22–30ms — dao động lớn |
| 14 000     | ❌ Vượt KPI  | P95 23–39ms, P99 gần chạm 50ms, dropped tăng mạnh |
| 14 500     | ❌ Vượt KPI  | P95 và P99 vượt ngưỡng, dropped 18k–47k/30s, kết quả không ổn định |

> **Saturation point** của 3 pods (0.5 CPU/pod): ~**13 000–13 500 RPS**
>
> ⚠️ Dropped cao ngay từ 13 500 RPS cho thấy K6 đang dùng ~500% CPU — bottleneck có phần ở K6, không hoàn toàn ở Node B

## So sánh tổng hợp

| Cấu hình | Saturation point | Tỉ lệ scale |
|:--------:|:----------------:|:-----------:|
| 1 pod    | ~4 700 RPS       | 1×          |
| 2 pods   | ~9 000 RPS       | ~1.9×       |
| 3 pods   | ~13 000 RPS      | ~2.8×       |

> Scale ngang **tuyến tính** — mỗi pod thêm vào đóng góp ~4 600–4 700 RPS.
> Tỉ lệ chưa đạt 3× hoàn toàn do CPU contention giữa K6 và Node B trên cùng máy.

## Scripts

**3 pods:**
```bash
docker compose -f docker-compose.2pod.yml down
docker compose up -d
cd k6 && TARGET_RPS=13500 k6 run load-test.js
```

---

_Test run: 2026-06-25 03:58:51_
