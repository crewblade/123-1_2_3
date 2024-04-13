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
            duration: '30s',
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
    //console.log("banners:", JSON.stringify(banners));

    const banner = randomItem(banners);
    const tagID = randomItem(banner.tag_ids);

    const pattern = randomIntBetween(0, 2);

    let url = `http://0.0.0.0:8080/banner?feature_id=${banner.feature_id}&tag_id=${tagID}`;

    // let url = `http://0.0.0.0:8080//banner?feature_id=${banner.feature_id}&tag_id=${tagID}&limit=${limit}&offset=${offset}`;
    // if (pattern === 0) {
    //     url = `http://0.0.0.0:8080/banner?tag_id=${tagID}&limit=${limit}&offset=${offset}`;
    // }
    // if (pattern === 1) {
    //     url = `http://0.0.0.0:8080/banner?feature_id=${banner.feature_id}&limit=${limit}&offset=${offset}`;
    // }
    // if (pattern === 2) {
    //     url = `http://0.0.0.0:8080/banner?&limit=${limit}&offset=${offset}`;
    // }


    const params = {
        headers: {
            'token': token,
        },
    };

    http.get(url, params);

}
