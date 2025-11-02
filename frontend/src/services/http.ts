import axios from 'axios';

export interface AppError {
  code?: number;
  message: string;
}

const http = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL ?? '/api/v1',
  timeout: 15000
});

http.interceptors.response.use(
  (response) => response,
  (error) => {
    const appError: AppError = {
      code: error.response?.data?.code,
      message:
        error.response?.data?.message ||
        error.message ||
        '请求失败，请稍后再试'
    };
    return Promise.reject(appError);
  }
);

export default http;

