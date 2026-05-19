import axios from "axios"

import { useUserStore } from "@/entities/user"

import type { AxiosError, InternalAxiosRequestConfig } from "axios"

// A clean instance for standard operations (like login, refresh, register)
// to avoid interceptor recursive loops.
export const authClient = axios.create({
  baseURL: "", // relative URL to current host
  withCredentials: true,
  headers: {
    "Content-Type": "application/json",
  },
})

// The automated API client that injects Bearer tokens and handles silent 401 refresh
export const apiClient = axios.create({
  baseURL: "", // relative URL to current host
  withCredentials: true,
  headers: {
    "Content-Type": "application/json",
  },
})

// State for concurrent refresh token requests
let isRefreshing = false
let failedQueue: {
  resolve: (token: string) => void
  reject: (err: unknown) => void
}[] = []

const processQueue = (error: unknown, token: string | null = null) => {
  failedQueue.forEach((promise) => {
    if (token !== null) {
      promise.resolve(token)
    } else {
      promise.reject(error)
    }
  })
  failedQueue = []
}

// Request Interceptor: Attach the Bearer Access Token to every request if present in store
apiClient.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    const token = useUserStore.getState().accessToken
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error: unknown) => {
    return Promise.reject(error instanceof Error ? error : new Error(String(error)))
  }
)

// Response Interceptor: Catch 401 and perform a silent token rotation
apiClient.interceptors.response.use(
  (response) => response,
  async (error: AxiosError) => {
    const originalRequest = error.config as InternalAxiosRequestConfig & { _retry?: boolean }

    // If unauthorized (401) and we haven't already retried this request
    if (error.response?.status === 401 && !originalRequest._retry) {
      // If we are already refreshing, wait for it to complete
      if (isRefreshing) {
        return new Promise<string>((resolve, reject) => {
          failedQueue.push({ resolve, reject })
        })
          .then((token) => {
            originalRequest.headers.Authorization = `Bearer ${token}`
            return apiClient(originalRequest)
          })
          .catch((err: unknown) => {
            return Promise.reject(err instanceof Error ? err : new Error(String(err)))
          })
      }

      originalRequest._retry = true
      isRefreshing = true

      try {
        // Call refresh endpoint. The HttpOnly cookie with the refresh token is sent automatically.
        const response = await authClient.post<{ access_token: string; expires_in: number }>("/auth/refresh")
        const newAccessToken = response.data.access_token
        const currentUser = useUserStore.getState().user

        // If we have user metadata, update the store with the new access token.
        if (currentUser) {
          useUserStore.getState().setSession(newAccessToken, currentUser)
        } else {
          // Fallback if user profile wasn't retrieved yet (update token only)
          useUserStore.setState({ accessToken: newAccessToken, isAuthenticated: true })
        }

        // Retry the original request
        originalRequest.headers.Authorization = `Bearer ${newAccessToken}`

        processQueue(null, newAccessToken)
        isRefreshing = false

        return await apiClient(originalRequest)
      } catch (refreshError: unknown) {
        // Refresh failed (cookie expired, session revoked, etc.)
        processQueue(refreshError, null)
        isRefreshing = false

        // Clear session and let application redirect to login
        useUserStore.getState().clearSession()

        // Programmatic redirect to login page
        if (window.location.pathname !== "/auth") {
          window.location.href = "/auth"
        }

        return Promise.reject(refreshError instanceof Error ? refreshError : new Error(String(refreshError)))
      }
    }

    return Promise.reject(error)
  }
)
