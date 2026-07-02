Link docs report: https://docs.google.com/document/d/1xMXLwuEkd8_BDtOI4U86cnfCsxofBwGgpv46WwWp8is/edit?usp=sharing

### How to run 

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