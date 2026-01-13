import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate, Trend } from 'k6/metrics';

// Custom metrics
const errorRate = new Rate('errors');
const requestDuration = new Trend('request_duration');

export const options = {
  stages: [
    { duration: '30s', target: 50 },   // Ramp up to 50 users
    { duration: '1m', target: 50 },     // Stay at 50 users
    { duration: '30s', target: 100 },   // Ramp up to 100 users
    { duration: '1m', target: 100 },   // Stay at 100 users
    { duration: '30s', target: 200 },   // Ramp up to 200 users
    { duration: '1m', target: 200 },   // Stay at 200 users
    { duration: '30s', target: 0 },     // Ramp down
  ],
  thresholds: {
    http_req_duration: ['p(95)<500'], // 95% of requests should be below 500ms
    errors: ['rate<0.1'],              // Error rate should be less than 10%
  },
};

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';

export default function () {
  // Test Token Bucket algorithm
  const payload = JSON.stringify({
    key: `user:${__VU}:${__ITER}`,
    algorithm: 'token_bucket',
    limit: 100,
    window: '1m',
  });

  const params = {
    headers: {
      'Content-Type': 'application/json',
    },
  };

  const startTime = Date.now();
  const response = http.post(`${BASE_URL}/api/v1/limit-check`, payload, params);
  const duration = Date.now() - startTime;

  requestDuration.add(duration);

  const success = check(response, {
    'status is 200 or 429': (r) => r.status === 200 || r.status === 429,
    'response has allowed field': (r) => {
      try {
        const body = JSON.parse(r.body);
        return 'allowed' in body;
      } catch {
        return false;
      }
    },
  });

  if (!success) {
    errorRate.add(1);
  } else {
    errorRate.add(0);
  }

  sleep(0.1); // 100ms between requests
}

