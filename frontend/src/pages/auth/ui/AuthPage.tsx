import { jwtDecode } from "jwt-decode"
import { useState, useEffect } from "react"
import { useNavigate, useSearchParams } from "react-router-dom"

import { LoginForm } from "@/features/auth/login"
import { OtpVerifyForm } from "@/features/auth/otp-verify"
import { RegisterForm } from "@/features/auth/register"

import { useAuthStore } from "@/entities/user"

import { api } from "@/shared/api"
import { 
  Card, 
  CardContent, 
  CardDescription, 
  CardHeader, 
  CardTitle 
} from "@/shared/ui/card"

type AuthStep = "email" | "register-request" | "otp"

interface DecodedToken {
  sub: number
  jti: string
  role: "admin" | "qp" | "warehouse_manager" | "storekeeper" | "pharmacist" | "unverified"
  email?: string
  exp: number
}

export function AuthPage() {
  const [searchParams] = useSearchParams()
  const [email, setEmail] = useState(() => searchParams.get("email") ?? "")
  const [step, setStep] = useState<AuthStep>(() => {
    const emailParam = searchParams.get("email")
    const stepParam = searchParams.get("step")
    return (emailParam && stepParam === "otp") ? "otp" : "email"
  })

  const [googleLoading, setGoogleLoading] = useState(false)
  const [googleError, setGoogleError] = useState<string | null>(null)

  const navigate = useNavigate()
  const setAuth = useAuthStore((state) => state.setAuth)

  useEffect(() => {
    const handleGoogleCallback = async () => {
      const hash = window.location.hash
      if (!hash) return

      const params = new URLSearchParams(hash.substring(1)) // Remove '#'
      const idToken = params.get("id_token")
      if (!idToken) return

      // Clear hash to prevent double verification
      window.history.replaceState({}, document.title, window.location.pathname + window.location.search)

      setGoogleLoading(true)
      setGoogleError(null)
      try {
        const response = await api.post<{ access_token: string }>("/auth/google", { id_token: idToken })
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

        void navigate("/")
      } catch (err) {
        const apiErr = err as { message?: string } | null
        setGoogleError(apiErr?.message ?? "Ошибка авторизации через Google.")
      } finally {
        setGoogleLoading(false)
      }
    }

    void handleGoogleCallback()
  }, [navigate, setAuth])

  const getStepHeader = () => {
    switch (step) {
      case "email":
        return {
          title: "Авторизация",
          description: "Введите электронную почту для входа или регистрации"
        }
      case "register-request":
        return {
          title: "Регистрация",
          description: "Создание новой учетной записи сотрудника"
        }
      case "otp":
        return {
          title: "Подтверждение",
          description: (
            <span>
              Введите одноразовый код для входа в систему, отправленный на адрес{" "}
              <span className="font-semibold text-muted-foreground/90 whitespace-nowrap">
                {email}
              </span>
            </span>
          )
        }
    }
  }

  const headerInfo = getStepHeader()

  if (googleLoading) {
    return (
      <div className="w-full max-w-[440px] animate-in fade-in zoom-in-95 duration-500">
        <Card className="border-border/40 shadow-2xl bg-card/50 backdrop-blur-md overflow-hidden min-h-[380px] flex flex-col items-center justify-center p-8">
          <div className="flex flex-col items-center space-y-4">
            <div className="h-8 w-8 animate-spin rounded-full border-4 border-primary border-t-transparent" />
            <p className="text-sm text-muted-foreground font-medium animate-pulse">
              Авторизация через Google...
            </p>
          </div>
        </Card>
      </div>
    )
  }

  return (
    <div className="w-full max-w-[440px] animate-in fade-in zoom-in-95 duration-500">
      <Card className="border-border/40 shadow-2xl bg-card/50 backdrop-blur-md overflow-hidden min-h-[380px] flex flex-col">
        <CardHeader className="space-y-1 text-center pt-8 pb-1">
          <CardTitle className="text-2xl font-bold tracking-tight">
            {headerInfo.title}
          </CardTitle>
          <CardDescription className="text-sm px-6">
            {headerInfo.description}
          </CardDescription>
        </CardHeader>
        
        <CardContent className="flex-1 px-8 pb-6 pt-1 flex flex-col justify-center">
          {step === "email" && (
            <LoginForm 
              onFound={(emailVal) => {
                setEmail(emailVal)
                setStep("otp")
              }}
              onNotFound={(emailVal) => {
                setEmail(emailVal)
                setStep("register-request")
              }}
              defaultEmail={email}
              externalError={googleError}
            />
          )}

          {step === "register-request" && (
            <RegisterForm 
              email={email}
              onBack={() => { setStep("email"); }}
              onSuccess={(emailVal) => {
                setEmail(emailVal)
                setStep("otp")
              }}
            />
          )}

          {step === "otp" && (
            <OtpVerifyForm 
              email={email}
              onBack={() => { setStep("email"); }}
            />
          )}
        </CardContent>
      </Card>
    </div>
  )
}
