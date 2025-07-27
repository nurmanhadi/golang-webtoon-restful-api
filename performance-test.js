import http from 'k6/http';
import { check } from 'k6';

export let options = {
    thresholds: {
        http_req_duration: ['p(95)<200'], // 95th percentile harus di bawah 200ms
    },
    scenarios: {
        load_test: {
            executor: 'ramping-vus',
            startVUs: 0,
            stages: [
                { duration: '10s', target: 50 },  // naik ke 50 VU
                { duration: '50s', target: 100 }, // naik ke 100 VU dan tahan
                { duration: '10s', target: 0 },   // turun ke 0 VU
            ],
        },
    },
};

export default function () {
    const url = 'http://localhost:3000/api/comics/9053b057-1db4-4e08-941c-63d5826523c5/chapters/10?number=10';

    let res = http.get(url);

    check(res, {
        'status is 200': (r) => r.status === 200,
    });
}
