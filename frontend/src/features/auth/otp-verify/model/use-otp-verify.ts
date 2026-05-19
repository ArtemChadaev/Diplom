import { useState, useEffect } from "react"
import { jwtDecode } from "jwt-decode"
import { useNavigate } from "react-router-dom"
import { api } from "@/shared/api"
import { useAuthStore } from "@/entities/user"

interface UseOtpVerifyProps {
  email: string
  onSuccess?: () => void
}

interface DecodedToken {
  sub: number
  jti: string
  role: "admin" | "qp" | "warehouse_manager" | "storekeeper" | "pharmacist" | "unverified"
  email: string
  exp: number
}

const REDIRECT_MAP: Record<string, string> = {
  admin: "/admin/users",
  qp: "/receiving",
  warehouse_manager: "/warehouse",
  storekeeper: "/search",
  pharmacist: "/search",
}

export function useOtpVerify({ email, onSuccess }: UseOtpVerifyProps) {
  const [code, setCode] = useState("")
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [secondsLeft, setSecondsLeft] = useState(600) // 10 minutes (600s)
  const [resendLoading, setResendLoading] = useState(false)
  const [resendMessage, setResendMessage] = useState<string | null>(null)

  const navigate = useNavigate()
  const setAuth = useAuthStore((state) => state.setAuth)

  // Countdown timer effect
  useEffect(() => {
    if (secondsLeft <= 0) return

    const interval = setInterval(() => {
      setSecondsLeft((s) => Math.max(0, s - 1))
    }, 1000)

    return () => clearInterval(interval)
  }, [secondsLeft])

  const verify = async (otpCode: string) => {
    setIsLoading(true)
    setError(null)
    try {
      const response = await api.post("/auth/verify-code", { email, code: otpCode })
      const { access_token } = response

      // Decode token to extract role, email, sub
      const decoded = jwtDecode<DecodedToken>(access_token)

      // Update Zustand store
      setAuth(access_token, decoded.role, {
        id: String(decoded.sub),
        email: decoded.email || email,
        full_name: "", // Will be fetched later from profile if needed
        role: decoded.role,
        ns_pv_access: false, // Default placeholders
        ukep_bound: false,
      })

      if (onSuccess) {
        onSuccess()
      }

      // Redirect by role
      const redirectPath = REDIRECT_MAP[decoded.role] || "/search"
      navigate(redirectPath)
    } catch (err: any) {
      if (err.status === 400) {
        setError(err.message || "Неверный код.")
      } else if (err.status === 410) {
        setError("Срок действия кода истек. Пожалуйста, запросите новый код.")
      } else if (err.status === 429) {
        setError("Превышено максимальное количество попыток ввода. Запросите новый код.")
      } else {
        setError(err.message || "Ошибка верификации кода. Попробуйте еще раз.")
      }
    } finally {
      setIsLoading(false)
    }
  }

  const resend = async () => {
    setResendLoading(true)
    setError(null)
    setResendMessage(null)
    try {
      await api.post("/auth/send-code", { email })
      setSecondsLeft(600) // Reset to 10 mins
      setResendMessage("Новый код отправлен на ваш email")
      setCode("") // Clear input code
    } catch (err: any) {
      setError(err.message || "Не удалось отправить повторный код.")
    } finally {
      setResendLoading(false)
    }
  }

  const formatTime = (seconds: number) => {
    const mins = Math.floor(seconds / 60)
    const secs = seconds % 60
    return `${mins}:${secs.toString().padStart(2, "0")}`
  }

  return {
    code,
    setCode,
    isLoading,
    error,
    secondsLeft,
    resendLoading,
    resendMessage,
    verify,
    resend,
    formatTime,
  }
}
