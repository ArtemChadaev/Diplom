import { create } from "zustand";
import { persist, createJSONStorage } from "zustand/middleware";

import type { UserDTO, UserRole } from "./types";

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
      setSession: (accessToken, user) => {
        set({
          accessToken,
          user,
          email: user.email,
          isAuthenticated: true,
        });
        // Sync with useAuthStore
        useAuthStore.setState({
          accessToken,
          role: user.role,
          user: {
            id: String(user.id),
            email: user.email,
            full_name: user.full_name,
            role: user.role,
            ns_pv_access: user.ns_pv_access,
            ukep_bound: user.ukep_bound,
          },
        });
      },
      clearSession: () => {
        set({
          accessToken: null,
          user: null,
          email: null,
          isAuthenticated: false,
        });
        // Sync with useAuthStore
        useAuthStore.setState({
          accessToken: null,
          role: null,
          user: null,
        });
      },
      updateUser: (userUpdates) =>
        set((state) => {
          const newUser = state.user ? { ...state.user, ...userUpdates } : null;
          if (newUser) {
            useAuthStore.setState((authState) => ({
              user: authState.user
                ? {
                    ...authState.user,
                    email: newUser.email,
                    full_name: newUser.full_name,
                    role: newUser.role,
                  }
                : null,
              role: newUser.role,
            }));
          }
          return {
            user: newUser,
            email: userUpdates.email ?? state.email,
          };
        }),
      setEmail: (email) => {
        set({ email });
        useAuthStore.setState((authState) => ({
          user: authState.user ? { ...authState.user, email } : null,
        }));
      },
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
  id: string;
  email: string;
  full_name: string;
  role: "admin" | "qp" | "warehouse_manager" | "storekeeper" | "pharmacist" | "unverified";
  ns_pv_access: boolean;
  ukep_bound: boolean;
}

interface AuthState {
  accessToken: string | null;
  role: string | null;
  user: User | null;
  setAuth: (accessToken: string, role: string, user: User | null) => void;
  logout: () => void;
}

export const useAuthStore = create<AuthState>((set) => ({
  accessToken: null,
  role: null,
  user: null,
  setAuth: (accessToken, role, user) => {
    set({ accessToken, role, user });
    if (user) {
      useUserStore.setState({
        accessToken,
        user: {
          id: Number(user.id),
          email: user.email,
          full_name: user.full_name,
          role: user.role as UserRole,
          ns_pv_access: user.ns_pv_access,
          ukep_bound: user.ukep_bound,
        },
        email: user.email,
        isAuthenticated: true,
      });
    }
  },
  logout: () => {
    set({ accessToken: null, role: null, user: null });
    useUserStore.setState({
      accessToken: null,
      user: null,
      email: null,
      isAuthenticated: false,
    });
  },
}));