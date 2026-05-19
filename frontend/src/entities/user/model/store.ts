import { create } from "zustand";
import { persist, createJSONStorage } from "zustand/middleware";

import type { UserDTO } from "./types";

interface UserState {
  accessToken: string | null;
  email: string | null;
  user: UserDTO | null;
  isAuthenticated: boolean;
  setSession: (accessToken: string, user: UserDTO) => void;
  clearSession: () => void;
  updateUser: (userUpdates: Partial<UserDTO>) => void;
  setEmail: (email: string) => void;
}

export const useUserStore = create<UserState>()(
  persist(
    (set) => ({
      accessToken: null, // JS memory only (secured from XSS)
      email: null,
      user: null,
      isAuthenticated: false,
      setSession: (accessToken, user) =>
        set({
          accessToken,
          user,
          email: user.email,
          isAuthenticated: true,
        }),
      clearSession: () =>
        set({
          accessToken: null,
          user: null,
          email: null,
          isAuthenticated: false,
        }),
      updateUser: (userUpdates) =>
        set((state) => ({
          user: state.user ? { ...state.user, ...userUpdates } : null,
          email: userUpdates.email ?? state.email,
        })),
      setEmail: (email) => set({ email }),
    }),
    {
      name: "pharma-hub-user-storage",
      storage: createJSONStorage(() => localStorage),
      // Partialize keeps only 'email' and 'user' in localStorage.
      // accessToken is kept purely in memory.
      partialize: (state) => ({
        email: state.email,
        user: state.user,
      }),
    }
  )
);

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