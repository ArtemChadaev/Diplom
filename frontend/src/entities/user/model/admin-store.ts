import { create } from "zustand"
import { z } from "zod"

import { api } from "@/shared/api"

// Схема Zod для отдельного профиля пользователя
export const userProfileSchema = z.object({
  id: z.number(),
  email: z.string().email(),
  role: z.string(), // "admin" | "qp" | "warehouse_manager" | "storekeeper" | "pharmacist"
  ns_pv_access: z.boolean(),
  ukep_bound: z.boolean(),
  is_blocked: z.boolean(),
  employee_code: z.string().nullable().optional(),
  full_name: z.string().nullable().optional(),
  position: z.string().nullable().optional(),
  department: z.string().nullable().optional(),
  avatar_url: z.string().nullable().optional(),
  telegram_handle: z.string().nullable().optional(),
})

// Схема Zod для списка пользователей
export const userListSchema = z.array(userProfileSchema)

export type UserProfile = z.infer<typeof userProfileSchema>

interface AdminUsersState {
  users: UserProfile[]
  isLoading: boolean
  error: string | null
  fetchUsers: () => Promise<void>
  toggleBlockUser: (id: number, blocked: boolean) => Promise<void>
}

export const useAdminUsersStore = create<AdminUsersState>((set) => ({
  users: [],
  isLoading: false,
  error: null,
  fetchUsers: async () => {
    set({ isLoading: true, error: null })
    try {
      const response = await api.get<unknown>("/admin/users")
      // Валидация ответа от сервера с помощью Zod
      const parsedUsers = userListSchema.parse(response)
      set({ users: parsedUsers, isLoading: false })
    } catch (err) {
      console.error("Failed to fetch users:", err)
      const message = err instanceof Error ? err.message : "Произошла ошибка при загрузке пользователей"
      set({ error: message, isLoading: false })
    }
  },
  toggleBlockUser: async (id: number, blocked: boolean) => {
    try {
      await api.patch(`/admin/users/${id}/blocked`, { blocked })
      // Локально обновляем статус блокировки
      set((state) => ({
        users: state.users.map((u) => (u.id === id ? { ...u, is_blocked: blocked } : u)),
      }))
    } catch (err) {
      console.error(`Failed to toggle block status for user ${id}:`, err)
      throw err;
    }
  },
}))
