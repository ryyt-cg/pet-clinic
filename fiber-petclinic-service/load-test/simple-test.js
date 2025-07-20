import http from 'k6/http';

export const options = {
    vus: 10, // Number of virtual users
    duration: '30s', // Duration of the test
}

export default () => {
    const random = Math.floor(Math.random() * 3 + 1);
    const url = 'http://localhost:8090/api/gof/v1/authors/' + random;
    const params = {
        headers: {
        'Content-Type': 'application/json',
        },
    };

    http.get(url);
}