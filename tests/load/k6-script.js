import http from 'k6/http';
import { check, group, sleep } from 'k6';
import { Rate } from 'k6/metrics';

// Custom metrics
const errorRate = new Rate('errors');

export let options = {
    stages: [
        { duration: '30s', target: 50 }, // Ramp up
        { duration: '1m', target: 500 }, // Load: 500 RPS target approx (depends on VUs)
        { duration: '30s', target: 0 },  // Ramp down
    ],
    thresholds: {
        http_req_duration: ['p(99)<200'], // 99% of requests must be < 200ms
        errors: ['rate<0.01'],            // < 1% error rate
    },
};

const BASE_URL = 'http://grocery-platform-gateway';

// User Behaviors
export default function () {
    const roll = Math.random();

    if (roll < 0.8) {
        // 80% Browse / Search
        group('Browse', function () {
            const res = http.get(`${BASE_URL}/api/search?q=milk`, {
                tags: { type: 'search' },
            });
            check(res, { 'status is 200': (r) => r.status === 200 }) || errorRate.add(1);
        });
    } else if (roll < 0.95) {
        // 15% Add to Cart
        group('Cart', function () {
            const payload = JSON.stringify({ itemId: 'item-123', quantity: 1 });
            const params = { headers: { 'Content-Type': 'application/json' } };
            const res = http.post(`${BASE_URL}/api/cart`, payload, params);
            check(res, { 'status is 200': (r) => r.status === 200 }) || errorRate.add(1);
        });
    } else {
        // 5% Checkout
        group('Checkout', function () {
            const payload = JSON.stringify({ userId: 'user-1', items: [{ itemId: 'item-123', quantity: 1 }] });
            const params = { headers: { 'Content-Type': 'application/json' } };
            const res = http.post(`${BASE_URL}/api/checkout/finalize`, payload, params);
            check(res, {
                'status is 200': (r) => r.status === 200,
                'no stock conflict': (r) => r.status !== 409,
            }) || errorRate.add(1);
        });
    }

    sleep(1);
}
