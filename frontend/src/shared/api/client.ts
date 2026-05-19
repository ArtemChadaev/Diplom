/* eslint-disable @typescript-eslint/no-explicit-any */
import { jwtDecode } from "jwt-decode"

import { useAuthStore } from "@/entities/user"

export const BASE_URL = import.meta.env.DEV 
  ? "" 
  : ((import.meta.env.VITE_API_URL as string | undefined) ?? "https://backend.pharma-hub.ru")

interface RequestOptions extends RequestInit {
  json?: unknown
}

class APIError extends Error {
  status: number
  data: unknown

  constructor(message: string, status: number, data: unknown) {
    super(message)
    this.name = "APIError"
    this.status = status
    this.data = data
  }
}

let refreshPromise: Promise<string> | null = null

async function request<T = any>(path: string, options: RequestOptions = {}): Promise<T> {
  const { json, headers: customHeaders, ...init } = options
  
  const headers = new Headers(customHeaders)
  
  if (json !== undefined) {
    headers.set("Content-Type", "application/json")
    init.body = JSON.stringify(json)
  }

  const token = useAuthStore.getState().accessToken
  if (token) {
    headers.set("Authorization", `Bearer ${token}`)
  }

  const prefix = path.startsWith("/auth") ? "" : "/api/v1"
  const response = await fetch(`${BASE_URL}${prefix}${path}`, {
    ...init,
    headers,
    credentials: "include",
  })

  // If unauthorized and not an auth path, try to silently refresh token
  if (response.status === 401 && !path.startsWith("/auth")) {
    if (!refreshPromise) {
      refreshPromise = (async () => {
        try {
          const refreshRes = await fetch(`${BASE_URL}/auth/refresh`, {
            method: "POST",
            credentials: "include",
            headers: {
              "Content-Type": "application/json",
            },
          })

          if (!refreshRes.ok) {
            throw new Error("Session expired")
          }

          const refreshData = await refreshRes.json() as { access_token: string }
          const newAccessToken = refreshData.access_token

          const decoded = jwtDecode<{
            sub: number
            role: "admin" | "qp" | "warehouse_manager" | "storekeeper" | "pharmacist" | "unverified"
            email?: string
          }>(newAccessToken)

          useAuthStore.getState().setAuth(newAccessToken, decoded.role, {
            id: String(decoded.sub),
            email: decoded.email ?? "",
            full_name: "",
            role: decoded.role,
            ns_pv_access: false,
            ukep_bound: false,
          })

          return newAccessToken
        } catch (err) {
          await useAuthStore.getState().logout()
          throw err
        } finally {
          refreshPromise = null
        }
      })()
    }

    try {
      const newAccessToken = await refreshPromise
      headers.set("Authorization", `Bearer ${newAccessToken}`)
      const retryResponse = await fetch(`${BASE_URL}${prefix}${path}`, {
        ...init,
        headers,
        credentials: "include",
      })

      let retryData: unknown
      const contentType = retryResponse.headers.get("content-type")
      if (contentType?.includes("application/json")) {
        retryData = await retryResponse.json()
      } else {
        retryData = await retryResponse.text()
      }

      if (!retryResponse.ok) {
        const errData = retryData as Record<string, unknown> | null
        const errMessage = (errData && typeof errData === "object" && (errData.error ?? errData.message))
          ? String(errData.error ?? errData.message)
          : `Request failed with status ${String(retryResponse.status)}`
        throw new APIError(
          errMessage,
          retryResponse.status,
          retryData
        )
      }

      return retryData as T
    } catch (refreshErr) {
      throw refreshErr
    }
  }

  let data: unknown
  const contentType = response.headers.get("content-type")
  if (contentType?.includes("application/json")) {
    data = await response.json()
  } else {
    data = await response.text()
  }

  if (!response.ok) {
    const errData = data as Record<string, unknown> | null
    const errMessage = (errData && typeof errData === "object" && (errData.error ?? errData.message))
      ? String(errData.error ?? errData.message)
      : `Request failed with status ${String(response.status)}`
    throw new APIError(
      errMessage,
      response.status,
      data
    )
  }

  return data as T
}

export const api = {
  get: <T = any>(path: string, options?: RequestOptions) => request<T>(path, { ...options, method: "GET" }),
  post: <T = any>(path: string, json?: unknown, options?: RequestOptions) => request<T>(path, { ...options, json, method: "POST" }),
  put: <T = any>(path: string, json?: unknown, options?: RequestOptions) => request<T>(path, { ...options, json, method: "PUT" }),
  patch: <T = any>(path: string, json?: unknown, options?: RequestOptions) => request<T>(path, { ...options, json, method: "PATCH" }),
  delete: <T = any>(path: string, options?: RequestOptions) => request<T>(path, { ...options, method: "DELETE" }),
}
