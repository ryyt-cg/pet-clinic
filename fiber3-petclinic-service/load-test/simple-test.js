import http from 'k6/http';

export const options = {
    vus: 10, // Number of virtual users
    duration: '30s', // Duration of the test
}

export default () => {
    const random = Math.floor(Math.random() * 7 + 1);
    const url = 'https://localhost:8092/api/pet-clinic/v1/owners/' + random;
    const params = {
        headers: {
        'Content-Type': 'application/json',
        },
    };

    http.get(url);
}