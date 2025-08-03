import axios from "axios";

const api = axios.create({
  baseURL: process.env.NEXT_PUBLIC_BACKEND_URL,
  timeout: 30000,
});

api.interceptors.response.use(
  (response) => response,
  (error) => {
    throw new Error(error?.response?.data?.message || error.message);
  }
);

export default api;
