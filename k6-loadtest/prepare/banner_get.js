import http from 'k6/http';

export const options = {
    vus: 1,
    iterations: 1,
};

export default function () {
    const url = 'http://0.0.0.0:8080/banner';
    const token = 'admin_token';
    const headers = {
        'token': token,
    };

    const response = http.get(url, { headers: headers });

    if (response.status === 200) {
        console.log(response.body);
    } else {
        console.error(`Error getting data banners: ${response.status}`);
    }
}
