/**
 * Ramp test: tăng dần RPS để tìm điểm bão hòa (saturation point).
 * Dùng cho báo cáo before/after tối ưu.
 *
 * Usage:
 *   k6 run ramp-test.js
 *   BASE_URL=http://localhost:8000 MAX_RPS=80000 k6 run ramp-test.js
 */

import http from 'k6/http';
import { check } from 'k6';

const BASE_URL         = __ENV.BASE_URL         || 'http://localhost:8000';
const MAX_RPS          = parseInt(__ENV.MAX_RPS          || '50000');
const SUBSCRIBER_COUNT = parseInt(__ENV.SUBSCRIBER_COUNT || '100000');

export const options = {
  scenarios: {
    ramp_load: {
      executor: 'ramping-arrival-rate',
      startRate: 1000,
      timeUnit: '1s',
      preAllocatedVUs: Math.ceil(MAX_RPS / 50),
      maxVUs:          Math.ceil(MAX_RPS / 5),
      stages: [
        { target: Math.floor(MAX_RPS * 0.25), duration: '30s' }, // 0→25% max
        { target: Math.floor(MAX_RPS * 0.50), duration: '30s' }, // 25→50%
        { target: Math.floor(MAX_RPS * 0.75), duration: '30s' }, // 50→75%
        { target: MAX_RPS,                    duration: '60s' }, // 75→100% sustained
        { target: 0,                          duration: '10s' }, // ramp down
      ],
    },
  },

  thresholds: {
    'http_req_duration': ['p(95)<20', 'p(99)<50'],
    'http_req_failed':   ['rate<0.001'],
  },
};

function randomSupi() {
  const n = Math.floor(Math.random() * SUBSCRIBER_COUNT);
  return `imsi-00101${String(n).padStart(10, '0')}`;
}

export default function () {
  const res = http.get(`${BASE_URL}/subscriber/${randomSupi()}`);

  check(res, {
    'status 200': (r) => r.status === 200,
  });
}
