import { create } from "zustand"

export interface User {
  id: string
  email: string
  full_name: string
  role: "admin" | "qp" | "warehouse_manager" | "storekeeper" | "pharmacist" | "unverified"
  ns_pv_access: boolean
  ukep_bound: boolean
}

interface AuthState {
  accessToken: string | null
  role: string | null
  user: User | null
  setAuth: (accessToken: string, role: string, user: User | null) => void
  logout: () => void
}

export const useAuthStore = create<AuthState>((set) => ({
  accessToken: null,
  role: null,
  user: null,
  setAuth: (accessToken, role, user) => set({ accessToken, role, user }),
  logout: () => set({ accessToken: null, role: null, user: null }),
}))
