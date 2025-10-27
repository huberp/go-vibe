// Basic smoke test - verifies system works under minimal load
import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  vus: 1, // 1 virtual user
  duration: '30s', // for 30 seconds
  thresholds: {
    http_req_duration: ['p(95)<500'], // 95% of requests should be below 500ms
    http_req_failed: ['rate<0.01'], // Error rate should be less than 1%
  },
};

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';

export default function () {
  // Test health endpoint
  let healthRes = http.get(`${BASE_URL}/health`);
  check(healthRes, {
    'health check status is 200': (r) => r.status === 200,
    'health check has status field': (r) => JSON.parse(r.body).status !== undefined,
  });

  // Test info endpoint
  let infoRes = http.get(`${BASE_URL}/info`);
  check(infoRes, {
    'info status is 200': (r) => r.status === 200,
  });

  sleep(1);
}
