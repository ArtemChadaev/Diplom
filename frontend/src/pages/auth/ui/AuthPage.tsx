import { useState, useEffect } from "react"
import { useSearchParams } from "react-router-dom"
import { 
  Card, 
  CardContent, 
  CardDescription, 
  CardHeader, 
  CardTitle 
} from "@/shared/ui/card"
import { LoginForm } from "@/features/auth/login"
import { RegisterForm } from "@/features/auth/register"
import { OtpVerifyForm } from "@/features/auth/otp-verify"

type AuthStep = "email" | "register-request" | "otp"

import { Button } from "@/shared/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle
} from "@/shared/ui/card";
import { Separator } from "@/shared/ui/separator";
import { Tabs, TabsList, TabsTrigger } from "@/shared/ui/tabs";


export function AuthPage() {
  const [searchParams] = useSearchParams()
  const [step, setStep] = useState<AuthStep>("email")
  const [email, setEmail] = useState("")

  // Handle query params on mount (for token expiration redirect case)
  useEffect(() => {
    const emailParam = searchParams.get("email")
    const stepParam = searchParams.get("step")

    if (emailParam) {
      setEmail(emailParam)
      if (stepParam === "otp") {
        setStep("otp")
      }
    }
  }, [searchParams])

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
          description: "Введите одноразовый код для входа в систему"
        }
    }
  }

  const headerInfo = getStepHeader()

  return (
    <div className="w-full max-w-[440px] animate-in fade-in zoom-in-95 duration-500">
      <Card className="border-border/40 shadow-2xl bg-card/50 backdrop-blur-md overflow-hidden min-h-[500px] flex flex-col">
        <CardHeader className="space-y-1 text-center pt-8 pb-4">
          <CardTitle className="text-2xl font-bold tracking-tight">
            {headerInfo.title}
          </CardTitle>
          <CardDescription className="text-sm px-6">
            {headerInfo.description}
          </CardDescription>
        </CardHeader>
        
        <CardContent className="flex-1 px-8 pb-8 pt-4 flex flex-col justify-center">
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
            />
          )}

          {step === "register-request" && (
            <RegisterForm 
              email={email}
              onBack={() => setStep("email")}
              onSuccess={(emailVal) => {
                setEmail(emailVal)
                setStep("otp")
              }}
            />
          )}

          {step === "otp" && (
            <OtpVerifyForm 
              email={email}
              onBack={() => setStep("email")}
            />
          )}
        </CardContent>
      </Card>
    </div>
  )
}
