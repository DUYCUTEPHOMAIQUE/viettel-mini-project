# Báo cáo: Thiết kế giao tiếp HTTP hiệu năng cao (100.000 TPS)

---

## 1. Bối cảnh & Mục tiêu

Thiết kế giao tiếp HTTP giữa Node A và Node B trong hệ thống UDM-like với mục tiêu:

- Throughput: **100.000 TPS**
- P95 latency: **< 20ms**
- P99 latency: **< 50ms**
- Error rate: **< 0.1%**

API duy nhất: `GET /subscriber/{supi}` — trả JSON tối giản.

---

## 2. Kiến trúc

```
Node A (K6)  ──HTTP/1.1 keep-alive──▶  HAProxy L4 :8000  ──roundrobin──▶  Node B Pod 1 :8080
                                        (HAProxy 2.9)                   ──▶  Node B Pod 2 :8080
                                                                         ──▶  Node B Pod N :8080
```

| Thành phần | Công nghệ | Vai trò |
|--|--|--|
| Node A | K6 | Load generator |
| Load Balancer | HAProxy 2.9 — TCP mode (L4) | Phân phối traffic |
| Node B | Go + fasthttp | API server stateless |

**Lý do chọn L4 thay L7:** HAProxy L4 forward TCP thuần, không parse HTTP header → overhead thấp hơn đáng kể ở TPS cao.

---

## 3. Tối ưu Tầng 1 — Code & HTTP

### 3.1 Baseline: net/http

Server Go dùng `net/http` + `json.Marshal` mỗi request.

| Target RPS | Actual RPS | P95 | P99 | Trạng thái |
|:--:|:--:|:--:|:--:|:--:|
| 3 000 | 2 980 | 15.78 ms | 33.04 ms | ✅ |
| 3 200 | 3 181 | 16.66 ms | 30.57 ms | ✅ |
| 3 400 | 3 361 | 31.62 ms | 49.72 ms | ❌ |
| 3 600 | 3 554 | 39.34 ms | 55.29 ms | ❌ |

**Saturation point: ~3 200 RPS** (1 pod, 0.5 CPU)

### 3.2 Optimized: fasthttp + pre-serialized JSON

Hai thay đổi chính:

**1. Thay net/http bằng fasthttp**
fasthttp có HTTP parser riêng, dùng object pool cho request/response → tránh allocation và GC pressure ở mỗi request.

**2. Pre-serialize JSON lúc khởi động**
```go
// Serialize 1 lần khi init, hot path chỉ đọc []byte
cache[supi] = []byte(`{"supi":"...","status":"REGISTERED","plmnId":"00101"}`)
```
Mỗi request chỉ cần `w.Write([]byte)` — không gọi `json.Marshal`, không tạo allocation.

| Target RPS | Actual RPS | P95 | P99 | Trạng thái |
|:--:|:--:|:--:|:--:|:--:|
| 4 500 | 4 498 | 757 µs | 2.32 ms | ✅ |
| 4 700 | 4 696 | 837 µs | 3.06 ms | ✅ |
| 5 000 | 4 986 | 4.09 ms | 22.18 ms | ❌ |

**Saturation point: ~4 700 RPS** (1 pod, 0.5 CPU)

### 3.3 So sánh

| Metric | net/http | fasthttp | Cải thiện |
|--|:--:|:--:|:--:|
| Saturation | ~3 200 RPS | ~4 700 RPS | **+43%** |
| P95 tại saturation | ~17 ms | < 1 ms | **~17×** |
| Alloc/request | 1 | 0 | **zero-alloc** |

---

## 4. Tối ưu Tầng 2 — Vertical Scaling

Tăng CPU limit của 1 pod để đánh giá hiệu quả scale dọc.

| CPU limit | Saturation point | Tỉ lệ |
|:--:|:--:|:--:|
| 0.5 CPU | ~4 700 RPS | 1× |
| 1.0 CPU | ~9 000 RPS | ~1.9× |

2× CPU → ~1.9× throughput. Không đạt 2× hoàn toàn do overhead cố định của HAProxy và Docker networking.

**Nhận xét:** Vertical scaling có giới hạn vật lý (không thể tăng CPU mãi) và không giải quyết được single point of failure. Horizontal scaling là hướng đúng cho bài toán 100k TPS.

---

## 5. Tối ưu Tầng 3 — Horizontal Scaling

Giữ nguyên 0.5 CPU/pod, tăng số pod và đo throughput tổng.

| Pods | Saturation point | Tỉ lệ |
|:--:|:--:|:--:|
| 1 pod | ~4 700 RPS | 1× |
| 2 pods | ~9 000 RPS | ~1.9× |
| 3 pods | ~13 000 RPS | ~2.8× |

Throughput tăng **gần tuyến tính** theo số pod — chứng minh kiến trúc stateless hoạt động đúng: không có shared state giữa các pod, HAProxy phân phối đều.

### Extrapolate lên 100.000 TPS

Từ dữ liệu thực đo:
```
Mỗi pod (0.5 CPU, fasthttp) ≈ 4 700 RPS
100 000 RPS ÷ 4 700 RPS/pod ≈ 22 pods
```

Với infrastructure thực (Node B cluster riêng, K6 trên máy riêng), 22 pod × 0.5 CPU = 11 core đủ để đạt 100k TPS trong KPI.

---

## 6. Phân tích Bottleneck

### 6.1 CPU Throttling (Docker cgroup)

Khi container đạt giới hạn CPU (`cpus: 0.5`), kernel Linux throttle theo chu kỳ 100ms:

```
Container được 50ms CPU / 100ms window
→ burst request → hết quota sau 50ms
→ bị freeze 50ms còn lại
→ latency median tăng lên ~100–140ms
```

Dấu hiệu nhận biết: median latency ≈ 100–140ms, CPU container đứng ở đúng mức giới hạn (50%, 100%...).

**Giải pháp:** Scale ngang (thêm pod) thay vì tăng CPU limit của 1 pod.

### 6.2 K6 CPU Saturation

Ở 13.000 RPS, K6 tiêu thụ ~500% CPU (5 core) do JavaScript engine (goja) có overhead cao. Trên môi trường single machine, K6 tranh CPU với Node B làm kết quả benchmark bị nhiễu.

**Giải pháp production:** Chạy K6 trên máy riêng hoặc dùng distributed K6 (nhiều instance song song).

### 6.3 HAProxy

HAProxy L4 không phải bottleneck ở mức TPS này. CPU cao (~85%) quan sát được trong benchmark là do chạy chung máy với K6, không phải do HAProxy thiếu capacity.

---

## 7. Monitoring

Trong quá trình benchmark, theo dõi 3 chỉ số sau là đủ để phát hiện saturation sớm:

### 7.1 TCP Connection Count

```bash
watch -n 1 'ss -tn state established dst :8000 | wc -l'
```

Tỉ lệ `connections/RPS` ổn định trong vùng bình thường, tăng đột biến khi hệ thống bão hòa:

| Config | RPS | Connections | Connections/RPS | Trạng thái |
|:------:|:---:|:-----------:|:---------------:|:----------:|
| 3 pods |  9 000 | ~420 | 0.047 | ✅ |
| 3 pods | 12 000 | ~525 | 0.044 | ✅ |
| 3 pods | 13 000 | ~989 | 0.076 | ❌ |

**Ngưỡng cảnh báo:** `connections/RPS > 0.06` → cần scale thêm pod.

### 7.2 CPU per Pod

```bash
docker stats --format "table {{.Name}}\t{{.CPUPerc}}\t{{.MemUsage}}"
```

Khi CPU container chạm đúng mức giới hạn (`cpus` limit) và đứng yên ở đó → đang bị cgroup throttle → latency sẽ spike. Scale ngang ngay thay vì chờ error rate tăng.

### 7.3 K6 Dropped Iterations

`dropped_iterations` trong K6 output là chỉ số nhanh nhất để phát hiện hệ thống không theo kịp target RPS. Khi dropped tăng đột biến so với các run trước ở cùng RPS → hệ thống đang vào vùng stress.

---

## 8. Kết luận

| Hạng mục | Kết quả |
|--|--|
| Framework tối ưu | fasthttp + pre-serialized JSON: +43% throughput, P95 giảm 17× |
| Vertical scaling | 2× CPU → ~1.9× throughput |
| Horizontal scaling | Tuyến tính, 3 pods → ~2.8× throughput |
| Kiến trúc stateless | Verified: không có shared state, scale ngang không cần sync |
| Projected 100k TPS | ~22 pods × 0.5 CPU = 11 core Node B cluster |

**Giới hạn môi trường demo:** Toàn bộ benchmark chạy trên 1 laptop (6 core). K6, HAProxy và Node B tranh CPU nhau làm số liệu ở mức TPS cao bị dao động. Trên môi trường production với các node riêng biệt, kết quả sẽ ổn định hơn và saturation point của mỗi pod sẽ cao hơn.

---

*Ngày thực hiện: 2026-06-25*
