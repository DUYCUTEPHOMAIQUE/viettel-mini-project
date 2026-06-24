# TCP Connection Analysis

**Quan sát:** Số TCP connection ESTABLISHED vào HAProxy :8000 trong lúc chạy load test.

```bash
# Lệnh quan sát
watch -n 1 'ss -tn state established dst :8000 | wc -l'
```

## Dữ liệu đo

| Config | RPS | TCP Connections | Connections/RPS |
|:------:|:---:|:---------------:|:---------------:|
| 1 pod  |  4 500 |  ~200 | 0.044 |
| 2 pods |  4 500 |  ~200 | 0.044 |
| 2 pods |  9 000 |  ~440 | 0.049 |
| 3 pods |  4 500 |  ~183 | 0.041 |
| 3 pods |  9 000 |  ~420 | 0.047 |
| 3 pods | 12 000 |  ~525 | 0.044 |
| 3 pods | 13 000 |  ~989 | **0.076** |

## Nhận xét

### 1. Connection count phụ thuộc vào RPS, không phụ thuộc số pod

Cùng 4 500 RPS: 1 pod → 200, 2 pod → 200, 3 pod → 183 — gần như bằng nhau.
K6 tạo VU dựa trên latency (Little's Law), không quan tâm đến topology phía server.

```
N = λ × W
N = số connection
λ = RPS
W = avg latency (giây)
```

### 2. Nhiều pod → ít connection hơn tại cùng RPS

| RPS | 2 pods | 3 pods |
|:---:|:------:|:------:|
| 9 000 | ~440 | ~420 |

Nhiều pod → latency thấp hơn → VU hoàn thành nhanh hơn → K6 cần ít VU hơn → ít connection hơn. Đúng với dự đoán từ Little's Law.

### 3. Connection là dấu hiệu saturation — xác định cliff point

| Config | RPS | Connections | Connections/RPS | Trạng thái |
|:------:|:---:|:-----------:|:---------------:|:----------:|
| 3 pods |  9 000 | ~420 | 0.047 | ✅ Bình thường |
| 3 pods | 12 000 | ~525 | 0.044 | ✅ Bình thường |
| 3 pods | 13 000 | ~989 | **0.076** | ❌ Bão hòa |

Từ 12k → 13k RPS (chỉ +8%), connection tăng từ 525 → 989 (+88%). Tỉ lệ connections/RPS nhảy vọt từ 0.044 lên 0.076 — **cliff point nằm trong khoảng 12 000–13 000 RPS**.

Điều này khớp với benchmark P95: tại 12k RPS P95 vẫn trong ngưỡng, tại 13k RPS bắt đầu vượt KPI.

**Cơ chế:**
```
Server saturation → latency tăng → VU bị block lâu hơn
→ K6 spawn thêm VU → số TCP connection bùng nổ
```

### 4. Ngưỡng cảnh báo thực tế

Nếu `connections/RPS > 0.06` → hệ thống đang vào vùng stress, cần scale thêm pod trước khi P95 vượt KPI.

Từ data: **saturation point thực của 3 pods ≈ 12 000–12 500 RPS** (không phải 13 500 như ước tính từ benchmark đơn thuần).

---

*Quan sát: 2026-06-25*
