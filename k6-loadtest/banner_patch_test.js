import { sleep } from 'k6';
import http from 'k6/http';
import { randomIntBetween, randomItem } from 'https://jslib.k6.io/k6-utils/1.2.0/index.js';

const banners = JSON.parse(open(`./prepare/banners.json`));


export const options = {
    scenarios: {
        load_test: {
            executor: 'constant-arrival-rate',
            rate: 1000,
            timeUnit: '1s',
            duration: '60s',
            preAllocatedVUs: 10,
            maxVUs: 40,
        },
    },
    thresholds: {
        http_req_failed: ['rate<0.0001'],
        http_req_duration: ['p(95)<200'],
    },
};

export default function () {

    const token = 'admin_token';

    const banner = randomItem(banners);
    const numTags = randomIntBetween(1, 3);
    const tagIds = [];
    const content = { title: "updated_title", text: "updated_text", url: "updated_url" };
    for (let i = 0; i < numTags; i++) {
        tagIds.push(randomIntBetween(0, 10000));
    }
    const featureId = randomIntBetween(0, 10000)
    const id = banner.banner_id
    let isActive = false;
    if (randomIntBetween(0, 100) < 50) {
        isActive = true;
    }

    http.patch(
        `http://0.0.0.0:8080/banner/${id}`,
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



}
