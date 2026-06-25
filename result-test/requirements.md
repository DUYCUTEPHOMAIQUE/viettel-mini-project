# MINIPROJECT: THIẾT KẾ GIAO TIẾP HTTP HIỆU NĂNG CAO (100.000 TPS)

## 1. Bối cảnh

Thiết kế giao tiếp HTTP giữa 2 node trong hệ thống UDM-like (Node A ↔ Node B), hướng tới tải lớn **100.000 requests/second (TPS)**.

---

## 2. Mục tiêu

Thiết kế giao tiếp HTTP hiệu năng cao, tối ưu latency, throughput và khả năng scale ngang.

---

## 3. Phạm vi

- **Node A**: client/NF gửi request
- **Node B**: API server xử lý request
- Giao thức: REST HTTP stateless

---

## 4. API chính

```
GET /subscriber/{supi}
```

- Response JSON tối giản, không payload dư thừa

---

## 5. Yêu cầu hiệu năng

- HTTP keep-alive bắt buộc
- Connection pooling phía client
- Stateless server
- Payload tối giản
- Hỗ trợ scale ngang nhiều instance
- Không tạo TCP connection per request

---

## 6. Kiến trúc

```
Node A → Load Balancer → Node B cluster
                     (Optional: Redis cache)
```

---

## 7. KPI mục tiêu

| Chỉ số | Mục tiêu |
|--------|----------|
| Throughput | 100.000 TPS (toàn hệ thống) |
| P95 latency | < 20ms |
| P99 latency | < 50ms |
| Error rate | < 0.1% |

---

## 8. Deliverables

- [ ] Thiết kế kiến trúc (diagram + mô tả)
- [ ] Source code demo HTTP service
- [ ] Load test script (k6/JMeter)
- [ ] Báo cáo before/after tối ưu

---

## 9. Thời gian

**2 tuần**

---

## 10. Tiêu chí đánh giá

| Hạng mục | Tỉ trọng |
|----------|----------|
| Thiết kế HTTP & tối ưu | 60% |
| Khả năng scale & phân tích bottleneck | 20% |
| Code & tooling | 20% |
