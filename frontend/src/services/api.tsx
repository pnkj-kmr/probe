
import axios from 'axios';

const api = axios.create({
  baseURL: process.env.REACT_APP_API_URL || 'https://api.example.com', // Replace with your base URL
  headers: {
    'Content-Type': 'application/json',
  },
});

// Optional: Request interceptor for adding auth tokens
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('authToken');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
}, (error) => Promise.reject(error));

// Optional: Response interceptor for handling errors
api.interceptors.response.use(
  (response) => response,
  (error) => {
    // Handle specific error codes (e.g., 401 logout)
    if (error.response?.status === 401) {
      console.warn('Unauthorized, redirecting to login...');
    }
    return Promise.reject(error);
  }
);

export default api;
