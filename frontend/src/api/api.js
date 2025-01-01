import axios from 'axios';

const API = axios.create({
  baseURL: 'http://localhost:8080', // Replace with your backend base URL
});

export const registerBag = (data) => API.post('/register-bag', data);
export const linkBags = (data) => API.post('/link-bags', data);
export const linkBagToBill = (data) => API.post('/link-bag-to-bill', data);
export const searchBill = (qrCode) => API.get(`/search-bill?qr_code=${qrCode}`);