import http from 'k6/http';
import { check } from 'k6';

export const options = {
    vus: 10,
    iterations: 10,
};

export default function () {
    const url = 'http://localhost:8080/api/v1/pending-order/create';

    const payload = JSON.stringify({
        payment_id: 1,
        address: "DN",
        name: "Doan Vu Duy",
        email: "doanvuduydtu@gmail.com",
        phone: "0384396351",
        total: 120,
        subtotal: 120,
        shipping: 120,
        tax: 1,
        items: [
            {
                variant_id: 84,
                quantity: 1,
                price: 24
            }
        ]
    });

    const params = {
        headers: {
            'Content-Type': 'application/json',
            'Authorization': 'Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJkYXRhIjoiQWRDb1pUaTBfcmJmR1lHVmQxbk9hbHhlUnlLNWtNcVQ5UERjT3NKdTh0LWdkdmxPNmxjUGlLY0xwV0dZVkRObzNvbko1Nk1aczczd2hMdmJRNkxzWnlzVjd6OG9mTUFhMlFxZS1jMHRmTkpmUldveHdKb1d0NWNSdFlpeVMwMzBoajRXYzcyX21XRDVpbnpmdENCSmdPUUxWYm1tZjBRSTVCTkY4b1hSMWVyMTV2R1hjYkFJcXJWcm1yX0pzbWt4dzRkN3Y1TkgweFpfV3Y5NXFnPT0iLCJleHAiOjE3OTE1NjM1NTUsImlhdCI6MTc5MTU2MjY1NSwianRpIjoiYTI1YzhjNzYtMjMzZC00MDBjLWIyMzEtODI0NzdmNjZhZjJkIn0.fbPN6nhDqq3wX13Soi6j5IRzD6RgUNIqt5T_Eby9bms',
        },
    };

    const res = http.post(url, payload, params);

    check(res, {
        'status is 200 or 400': (r) => r.status === 200 || r.status === 400,
    });
}