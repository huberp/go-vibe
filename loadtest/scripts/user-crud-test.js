// Full user CRUD operations load test
import http from 'k6/http';
import { check, sleep } from 'k6';
import { randomString } from 'https://jslib.k6.io/k6-utils/1.2.0/index.js';

export const options = {
  stages: [
    { duration: '1m', target: 20 },  // Ramp up to 20 users
    { duration: '3m', target: 20 },  // Stay at 20 users
    { duration: '1m', target: 0 },   // Ramp down
  ],
  thresholds: {
    http_req_duration: ['p(95)<1500'],
    http_req_failed: ['rate<0.05'],
    'http_req_duration{endpoint:create_user}': ['p(95)<2000'],
    'http_req_duration{endpoint:login}': ['p(95)<1500'],
    'http_req_duration{endpoint:get_users}': ['p(95)<1000'],
  },
};

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';
const ADMIN_EMAIL = __ENV.ADMIN_EMAIL || 'admin@example.com';
const ADMIN_PASSWORD = __ENV.ADMIN_PASSWORD || 'password123';

let adminToken = null;

export function setup() {
  // Login as admin to get token for the test
  const loginPayload = JSON.stringify({
    email: ADMIN_EMAIL,
    password: ADMIN_PASSWORD,
  });

  const loginRes = http.post(`${BASE_URL}/v1/login`, loginPayload, {
    headers: { 'Content-Type': 'application/json' },
  });

  if (loginRes.status === 200) {
    return { adminToken: JSON.parse(loginRes.body).token };
  }
  
  console.warn('Setup: Could not get admin token, some tests may fail');
  return { adminToken: null };
}

export default function (data) {
  const token = data.adminToken;

  // 1. Create a new user
  const randomEmail = `user-${randomString(10)}@loadtest.com`;
  const createUserPayload = JSON.stringify({
    name: `Load Test User ${randomString(5)}`,
    email: randomEmail,
    password: 'testpassword123',
    role: 'user',
  });

  const createParams = {
    headers: { 'Content-Type': 'application/json' },
    tags: { endpoint: 'create_user' },
  };

  const createRes = http.post(`${BASE_URL}/v1/users`, createUserPayload, createParams);
  const createSuccess = check(createRes, {
    'create user status is 201': (r) => r.status === 201,
    'create user returns user object': (r) => JSON.parse(r.body).id !== undefined,
  });

  sleep(1);

  // 2. Login with the new user
  if (createSuccess) {
    const loginPayload = JSON.stringify({
      email: randomEmail,
      password: 'testpassword123',
    });

    const loginRes = http.post(`${BASE_URL}/v1/login`, loginPayload, {
      headers: { 'Content-Type': 'application/json' },
      tags: { endpoint: 'login' },
    });

    check(loginRes, {
      'login status is 200': (r) => r.status === 200,
      'login returns token': (r) => JSON.parse(r.body).token !== undefined,
    });

    sleep(1);
  }

  // 3. Get all users (as admin)
  if (token) {
    const getUsersRes = http.get(`${BASE_URL}/v1/users`, {
      headers: { 'Authorization': `Bearer ${token}` },
      tags: { endpoint: 'get_users' },
    });

    check(getUsersRes, {
      'get users status is 200': (r) => r.status === 200,
      'get users returns array': (r) => Array.isArray(JSON.parse(r.body)),
    });

    sleep(1);
  }

  // 4. Test health and metrics
  http.get(`${BASE_URL}/health`, { tags: { endpoint: 'health' } });
  http.get(`${BASE_URL}/metrics`, { tags: { endpoint: 'metrics' } });

  sleep(2);
}
