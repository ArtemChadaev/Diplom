import { useAuthStore } from "@/entities/user"

const BASE_URL = "/api/v1"

interface RequestOptions extends RequestInit {
  json?: any
}

class APIError extends Error {
  status: number
  data: any

  constructor(message: string, status: number, data: any) {
    super(message)
    this.name = "APIError"
    this.status = status
    this.data = data
  }
}

async function request(path: string, options: RequestOptions = {}) {
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

  const response = await fetch(`${BASE_URL}${path}`, {
    ...init,
    headers,
  })

  let data: any = null
  const contentType = response.headers.get("content-type")
  if (contentType && contentType.includes("application/json")) {
    data = await response.json()
  } else {
    data = await response.text()
  }

  if (!response.ok) {
    if (response.status === 401 && !path.includes("/auth/refresh")) {
      const email = useAuthStore.getState().user?.email
      useAuthStore.getState().logout()
      if (email) {
        window.location.href = `/auth?step=otp&email=${encodeURIComponent(email)}`
      } else {
        window.location.href = "/auth"
      }
    }
    throw new APIError(
      data?.error || data?.message || `Request failed with status ${response.status}`,
      response.status,
      data
    )
  }

  return data
}

export const api = {
  get: (path: string, options?: RequestOptions) => request(path, { ...options, method: "GET" }),
  post: (path: string, json?: any, options?: RequestOptions) => request(path, { ...options, json, method: "POST" }),
  put: (path: string, json?: any, options?: RequestOptions) => request(path, { ...options, json, method: "PUT" }),
  patch: (path: string, json?: any, options?: RequestOptions) => request(path, { ...options, json, method: "PATCH" }),
  delete: (path: string, options?: RequestOptions) => request(path, { ...options, method: "DELETE" }),
}
