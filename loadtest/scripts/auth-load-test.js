// Load test for authentication endpoints
import http from 'k6/http';
import { check, sleep } from 'k6';
import { SharedArray } from 'k6/data';

// Load test users from fixture data
const users = new SharedArray('users', function () {
  return [
    { email: 'test@example.com', password: 'password123', role: 'user' },
    { email: 'admin@example.com', password: 'password123', role: 'admin' },
    { email: 'john@example.com', password: 'password123', role: 'user' },
    { email: 'jane@example.com', password: 'password123', role: 'admin' },
    { email: 'bob@example.com', password: 'password123', role: 'user' },
  ];
});

export const options = {
  stages: [
    { duration: '30s', target: 10 }, // Ramp up to 10 users
    { duration: '1m', target: 10 },  // Stay at 10 users for 1 minute
    { duration: '30s', target: 0 },  // Ramp down to 0 users
  ],
  thresholds: {
    http_req_duration: ['p(95)<1000'], // 95% of requests should be below 1s
    http_req_failed: ['rate<0.05'],    // Error rate should be less than 5%
    'http_req_duration{endpoint:login}': ['p(95)<1500'], // Login can be slower
  },
};

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';

export default function () {
  // Select a random user for load testing
  // Note: Math.random() is acceptable here as this is for test data selection,
  // not for any security-sensitive operations
  const user = users[Math.floor(Math.random() * users.length)];

  // Test login
  const loginPayload = JSON.stringify({
    email: user.email,
    password: user.password,
  });

  const loginParams = {
    headers: {
      'Content-Type': 'application/json',
    },
    tags: { endpoint: 'login' },
  };

  const loginRes = http.post(`${BASE_URL}/v1/login`, loginPayload, loginParams);
  
  const loginSuccess = check(loginRes, {
    'login status is 200': (r) => r.status === 200,
    'login returns token': (r) => JSON.parse(r.body).token !== undefined,
  });

  if (loginSuccess) {
    const token = JSON.parse(loginRes.body).token;

    // Test authenticated endpoints
    const authParams = {
      headers: {
        'Authorization': `Bearer ${token}`,
      },
      tags: { endpoint: 'get_users' },
    };

    // Get all users (admin only)
    if (user.role === 'admin') {
      const getUsersRes = http.get(`${BASE_URL}/v1/users`, authParams);
      check(getUsersRes, {
        'get users status is 200 or 403': (r) => r.status === 200 || r.status === 403,
      });
    }

    // Get user by ID
    const getUserRes = http.get(`${BASE_URL}/v1/users/1`, authParams);
    check(getUserRes, {
      'get user status is 200 or 403': (r) => r.status === 200 || r.status === 403,
    });
  }

  sleep(1);
}
