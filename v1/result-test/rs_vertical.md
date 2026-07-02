# Load Test Results — Vertical Scaling (fasthttp, 1 pod)

**Config:** 1 pod Node B (fasthttp) · HAProxy L4 · duration 30s/run

## 1.0 CPU

| Run | Target RPS | Actual RPS | Avg      | P95      | P99      | Error Rate | Dropped |
|:---:|:----------:|:----------:|:--------:|:--------:|:--------:|:----------:|:-------:|
|  1  |   7 000    |   6 990    | 602 µs   |  1.95 ms |  3.40 ms |   0.00%    |   274   |
|  2  |   8 000    |   7 989    | 768 µs   |  2.73 ms |  5.84 ms |   0.00%    |   306   |
|  3  |   9 000    |   8 999    | 687 µs   |  2.58 ms |  4.71 ms |   0.00%    |     0   |
|  4  |  10 000    |   9 936    |  2.87 ms | 16.72 ms | 26.88 ms |   0.00%    | 1 795   |
|  5  |  11 000    |  10 667    | 11.66 ms | 42.38 ms | 52.32 ms |   0.00%    | 9 802   |
|  6  |  12 000    |  11 244    | 22.86 ms | 65.66 ms | 77.30 ms |   0.00%    |22 653   |

## Nhận xét

| Target RPS | Trạng thái | Ghi chú |
|:----------:|:----------:|---------|
|  7 000     | ✅ Đạt KPI  | P95 1.95ms — rất thoải mái |
|  8 000     | ✅ Đạt KPI  | P95 2.73ms — ổn định |
|  9 000     | ✅ Đạt KPI  | P95 2.58ms — dropped = 0, peak tốt nhất |
| 10 000     | ⚠️ Cận ngưỡng | P95 16.72ms, bắt đầu có dropped |
| 11 000     | ❌ Vượt KPI | P95 42ms, P99 52ms |
| 12 000     | ❌ Vượt KPI | P95 65ms, P99 77ms, dropped cao |

> **Saturation point** của 1 pod fasthttp (1.0 CPU): ~**9 000–10 000 RPS**

## So sánh Vertical Scaling

| CPU limit | Saturation point | Tỉ lệ vs 0.5 CPU |
|:---------:|:----------------:|:----------------:|
| 0.5 CPU   | ~4 700 RPS       | 1×               |
| 1.0 CPU   | ~9 000 RPS       | ~1.9×            |

> Scale dọc 2× CPU → throughput tăng ~**1.9×** — gần tuyến tính nhưng không hoàn toàn do overhead Docker + HAProxy cố định.

---

_Test run: 2026-06-25 04:27:00_
