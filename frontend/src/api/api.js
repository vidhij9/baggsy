import axios from "axios";

const API = axios.create({
  baseURL: "http://localhost:8080", // Ensure this is correct
  validateStatus: function (status) {
    return status >= 200 && status < 300; // Default range for successful status codes
  },
});

export const registerBag = (data) => API.post("/register-bag", data);
export const linkBags = (data) => API.post("/link-bags", data);
export const searchBill = (qrCode) => API.get(`/search-bill?qr_code=${qrCode}`);