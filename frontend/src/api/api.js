import axios from "axios";

const API = axios.create({
  baseURL: "http://localhost:8080", // Ensure this is correct
  validateStatus: function (status) {
    return status >= 200 && status < 300; // Default range for successful status codes
  },
});

export const searchBill = (qrCode) => API.get(`/search-bill?qr_code=${qrCode}`);
export const registerBag = (data) => API.post("/register-bag", data);
export const linkChildBag = (data) => API.post("/link-child-bag", data);
export const linkBagToBill = (data) => API.post("/link-bag-to-bill", data);