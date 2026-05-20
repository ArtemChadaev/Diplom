import { jwtDecode } from "jwt-decode"
import { useEffect, useState } from "react"
import { Outlet, useNavigate } from "react-router-dom"

import { Footer } from "@/widgets/footer"
import { Header } from "@/widgets/header"
import { AppSidebar } from "@/widgets/sidebar"

import { useAuthStore } from "@/entities/user"

import { api } from "@/shared/api"
import { SidebarProvider } from "@/shared/ui/sidebar"
import { TooltipProvider } from "@/shared/ui/tooltip"

interface DecodedToken {
  sub: number
  jti: string
  role: "admin" | "qp" | "warehouse_manager" | "storekeeper" | "pharmacist" | "unverified"
  email?: string
  exp: number
}

export function AppLayout() {
  const { accessToken, setAuth, logout, user } = useAuthStore()
  const navigate = useNavigate()
  const [isInitializing, setIsInitializing] = useState(true)

  useEffect(() => {
    const checkAuth = async () => {
      if (accessToken) {
        setIsInitializing(false)
        return
      }

      try {
        // Try to refresh token using HTTPOnly cookie
        const response = await api.post<{ access_token: string }>("/auth/refresh")
        const { access_token } = response
        
        const decoded = jwtDecode<DecodedToken>(access_token)
        
        setAuth(access_token, decoded.role, {
          id: String(decoded.sub),
          email: decoded.email ?? "",
          full_name: "",
          role: decoded.role,
          ns_pv_access: false,
          ukep_bound: false,
        })
      } catch {
        await logout()
      } finally {
        setIsInitializing(false)
      }
    }

    void checkAuth()
  }, [accessToken, setAuth, logout, navigate, user?.email])

  if (isInitializing) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-background">
        <div className="flex flex-col items-center space-y-4">
          <div className="h-8 w-8 animate-spin rounded-full border-4 border-primary border-t-transparent" />
          <p className="text-sm text-muted-foreground font-medium animate-pulse">Загрузка сессии...</p>
        </div>
      </div>
    )
  }

  if (!accessToken) {
    return null
  }

  return (
    <TooltipProvider>
      <SidebarProvider>
        <div className="min-h-screen flex w-full bg-background selection:bg-secondary/20">
          <AppSidebar />
          <div className="flex-1 flex flex-col min-w-0">
            <Header />
            <main className="flex-1 w-full mx-auto p-6 md:p-10">
              <Outlet />
            </main>
            <Footer />
          </div>
        </div>
      </SidebarProvider>
    </TooltipProvider>
  )
}
