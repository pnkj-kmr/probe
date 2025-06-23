import axios from "axios";
import environment from "@/Environement";

let baseURL = environment.API;
if (environment.SAME_ORIGIN) {
  baseURL = window.location.origin;
}

const api = axios.create({
  baseURL: baseURL,
  headers: { "Content-Type": "application/json" },
});

// Optional: Request interceptor for adding auth tokens
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem("authToken");
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => Promise.reject(error)
);

// Optional: Response interceptor for handling errors
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      console.warn("Unauthorized, redirecting to login...");
    }
    return Promise.reject(error);
  }
);

export default api;
