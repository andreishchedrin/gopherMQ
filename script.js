import http from 'k6/http';
import { check, sleep } from 'k6';

export default function () {
    const url = 'http://127.0.0.1:8888/push';
    const payload = JSON.stringify({
        channel: "Hello1",
        message: "dfgdfgdfg34634-test"
    });

    const params = {
        headers: {
            'Content-Type': 'application/json',
        },
    };

    const res = http.post(url, payload, params);
    check(res, { 'status was 200': (r) => r.status == 200 });
}