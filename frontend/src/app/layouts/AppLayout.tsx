import { useEffect, useState } from "react"
import { Outlet, useNavigate } from "react-router-dom"
import { Header } from "@/widgets/header"
import { Footer } from "@/widgets/footer"
import { useAuthStore } from "@/entities/user"
import { api } from "@/shared/api"
import { jwtDecode } from "jwt-decode"

interface DecodedToken {
  sub: number
  jti: string
  role: "admin" | "qp" | "warehouse_manager" | "storekeeper" | "pharmacist" | "unverified"
  email: string
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
        const response = await api.post("/auth/refresh")
        const { access_token } = response
        
        const decoded = jwtDecode<DecodedToken>(access_token)
        
        setAuth(access_token, decoded.role, {
          id: String(decoded.sub),
          email: decoded.email,
          full_name: "",
          role: decoded.role,
          ns_pv_access: false,
          ukep_bound: false,
        })
      } catch (err: any) {
        logout()
        const currentEmail = user?.email || ""
        if (currentEmail) {
          navigate(`/auth?step=otp&email=${encodeURIComponent(currentEmail)}`)
        } else {
          navigate("/auth")
        }
      } finally {
        setIsInitializing(false)
      }
    }

    checkAuth()
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
    <div className="min-h-screen flex flex-col bg-background selection:bg-secondary/20">
      <Header />
      <main className="flex-1 w-full mx-auto p-6 md:p-10">
        <Outlet />
      </main>
      <Footer />
    </div>
  )
}
