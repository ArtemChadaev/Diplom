import { jwtDecode } from "jwt-decode"
import { useState, useEffect } from "react"
import { useNavigate } from "react-router-dom"

import { useAuthStore } from "@/entities/user"

import { api } from "@/shared/api"

interface UseOtpVerifyProps {
  email: string
  onSuccess?: () => void
}

interface DecodedToken {
  sub: number
  jti: string
  role: "admin" | "qp" | "warehouse_manager" | "storekeeper" | "pharmacist" | "unverified"
  email?: string
  exp: number
}


export function useOtpVerify({ email, onSuccess }: UseOtpVerifyProps) {
  const [code, setCode] = useState("")
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [secondsLeft, setSecondsLeft] = useState(60) // 1 minute (60s)
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

    return () => { clearInterval(interval); }
  }, [secondsLeft])

  const verify = async (otpCode: string) => {
    setIsLoading(true)
    setError(null)
    try {
      const response = await api.post<{ access_token: string }>("/auth/verify-code", { email, code: otpCode })
      const { access_token } = response

      // Decode token to extract role, email, sub
      const decoded = jwtDecode<DecodedToken>(access_token)

      // Update Zustand store
      setAuth(access_token, decoded.role, {
        id: String(decoded.sub),
        email: decoded.email ?? email,
        full_name: "", // Will be fetched later from profile if needed
        role: decoded.role,
        ns_pv_access: false, // Default placeholders
        ukep_bound: false,
      })

      if (onSuccess) {
        onSuccess()
      }

      void navigate("/")
    } catch (err) {
      setCode("") // Clear input code on error
      const apiErr = err as { status?: number; message?: string } | null
      if (apiErr?.status === 400 || apiErr?.status === 401) {
        setError(apiErr.message ?? "Неверный код.")
      } else if (apiErr?.status === 410) {
        setError("Срок действия кода истек. Пожалуйста, запросите новый код.")
      } else if (apiErr?.status === 429) {
        setError("Превышено максимальное количество попыток ввода. Запросите новый код.")
      } else {
        setError(apiErr?.message ?? "Ошибка верификации кода. Попробуйте еще раз.")
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
      setSecondsLeft(60) // Reset to 1 min (60s)
      setResendMessage("Новый код отправлен на ваш email")
      setCode("") // Clear input code
    } catch (err) {
      const apiErr = err as { message?: string } | null
      setError(apiErr?.message ?? "Не удалось отправить повторный код.")
    } finally {
      setResendLoading(false)
    }
  }

  const formatTime = (seconds: number) => {
    const mins = Math.floor(seconds / 60)
    const secs = seconds % 60
    return `${String(mins)}:${secs.toString().padStart(2, "0")}`
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
