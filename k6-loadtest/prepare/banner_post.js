import { sleep } from 'k6';
import http from 'k6/http';
import { randomIntBetween } from 'https://jslib.k6.io/k6-utils/1.2.0/index.js';

export const options = {
    scenarios: {
        load_test: {
            executor: 'constant-arrival-rate',
            rate: 100,
            timeUnit: '1s',
            duration: '3m',
            preAllocatedVUs: 10,
            maxVUs: 200,
        },
    },
    thresholds: {
        http_req_failed: ['rate<0.0001'],
        http_req_duration: ['p(95)<200'],
    },
};

export default function () {
    const token = 'admin_token';
    const numTags = randomIntBetween(1, 3);
    const tagIds = [];
    for (let i = 0; i < numTags; i++) {
        tagIds.push(randomIntBetween(1, 600));
    }
    const featureId = randomIntBetween(1, 600);
    const content = { title: "some_title", text: "some_text", url: "some_url" };
    let isActive = false;
    if (randomIntBetween(0, 100) < 90) {
        isActive = true;
    }

    http.post(
        'http://0.0.0.0:8080/banner',
        JSON.stringify({
            tag_ids: tagIds,
            feature_id: featureId,
            content: content,
            is_active: isActive
        }),
        {
            headers: {
                'Content-Type': 'application/json',
                'token': token,
            },
        },
    );

    sleep(1);
}
