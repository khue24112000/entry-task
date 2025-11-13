import axios, { AxiosError } from "axios";

const baseURL = "http://localhost:8080/api/v1";

export const publicRoute = axios.create({
  baseURL,
  withCredentials: false,
  headers: {
    "Content-Type": "application/json",
  },
});

export const privateRoute = axios.create({
  baseURL,
  withCredentials: true,
  headers: {
    "Content-Type": "application/json",
  },
});

privateRoute.interceptors.response.use(
  (response) => response,
  async (error: AxiosError) => {
    const originalRequest: any = error.config;

    if (
      error.response?.status === 401 &&
      !originalRequest._retry &&
      originalRequest.url !== "/refresh" &&
      originalRequest.url !== "/login"
    ) {
      originalRequest._retry = true;
      try {
        await privateRoute.post("/refresh");
        return privateRoute(originalRequest);
      } catch (refreshErr) {
        return Promise.reject(refreshErr);
      }
    }

    return Promise.reject(error);
  }
);
