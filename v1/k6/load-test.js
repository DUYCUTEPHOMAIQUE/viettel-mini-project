/**
 * Load test: constant-arrival-rate (TPS-controlled)
 *
 * Usage:
 *   k6 run load-test.js
 *   TARGET_RPS=50000 DURATION=60s k6 run load-test.js
 *
 * Distributed (4 process × 25k = 100k TPS):
 *   for i in 1 2 3 4; do TARGET_RPS=25000 k6 run load-test.js & done; wait
 */

import http from 'k6/http';
import { check } from 'k6';

const BASE_URL         = __ENV.BASE_URL         || 'http://localhost:8000';
const TARGET_RPS       = parseInt(__ENV.TARGET_RPS       || '10000');
const DURATION         = __ENV.DURATION         || '60s';
const SUBSCRIBER_COUNT = parseInt(__ENV.SUBSCRIBER_COUNT || '100000');

export const options = {
  scenarios: {
    constant_load: {
      executor:        'constant-arrival-rate',
      rate:            TARGET_RPS,
      timeUnit:        '1s',
      duration:        DURATION,
      preAllocatedVUs: Math.ceil(TARGET_RPS / 50),
      maxVUs:          Math.ceil(TARGET_RPS / 10),
    },
  },
  thresholds: {
    'http_req_duration': ['p(95)<20', 'p(99)<50'],
    'http_req_failed':   ['rate<0.001'],
  },
};

// pre-build SUPI list một lần duy nhất, tránh String() + padStart() mỗi request
const SUPIS = (() => {
  const arr = new Array(1000);
  for (let i = 0; i < 1000; i++) {
    const n = Math.floor(Math.random() * SUBSCRIBER_COUNT);
    arr[i] = `${BASE_URL}/subscriber/imsi-00101${String(n).padStart(10, '0')}`;
  }
  return arr;
})();

// ── hot path: chỉ làm đúng 1 việc ───────────────────────────────────────────
export default function () {
  const url = SUPIS[Math.floor(Math.random() * 1000)];
  const res = http.get(url);
  check(res, { 'status 200': (r) => r.status === 200 });
}
