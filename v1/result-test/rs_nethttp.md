# Load Test Results — Baseline (net/http)

**Config:** 1 pod Node B (net/http) · HAProxy L4 · `cpus: 0.5` · `memory: 256M` · duration 30s/run

## Kết quả

| Run | Target RPS | Actual RPS | Avg     | P95      | P99      | Error Rate | Dropped |
|:---:|:----------:|:----------:|:-------:|:--------:|:--------:|:----------:|:-------:|
|  1  |   3 000    |   2 980    | 2.36 ms | 15.78 ms | 33.04 ms |   0.00%    |   596   |
|  2  |   3 200    |   3 181    | 2.69 ms | 16.66 ms | 30.57 ms |   0.00%    |   567   |
|  3  |   3 400    |   3 361    | 6.45 ms | 31.62 ms | 49.72 ms |   0.00%    | 1 144   |
|  4  |   3 600    |   3 554    | 9.37 ms | 39.34 ms | 55.29 ms |   0.00%    | 1 373   |

## Nhận xét

| Target RPS | Trạng thái | Ghi chú |
|:----------:|:----------:|---------|
| 3 000      | ✅ Đạt KPI  | P95 15.78ms, P99 33ms — còn headroom |
| 3 200      | ✅ Đạt KPI  | P95 16.66ms, P99 30ms — cận ngưỡng P95 |
| 3 400      | ❌ Vượt KPI | P95 31.62ms vượt 20ms, P99 gần chạm 50ms |
| 3 600      | ❌ Vượt KPI | P95 39ms, P99 55ms vượt cả 2 ngưỡng |

> **Saturation point** của 1 pod net/http (0.5 CPU): ~**3 200 RPS**

## So sánh với fasthttp

| Framework  | Saturation point | P95 tại saturation | Alloc/request |
|:----------:|:----------------:|:------------------:|:-------------:|
| net/http   | ~3 200 RPS       | ~17 ms             | 1 (json.Marshal) |
| fasthttp   | ~4 700 RPS       | < 1 ms             | 0 (pre-serialized) |
| **Cải thiện** | **+43%**      | **~17×**           | **zero-alloc** |

---

_Test run: 2026-06-25 04:21:35_
