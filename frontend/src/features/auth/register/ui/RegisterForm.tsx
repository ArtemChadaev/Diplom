import { useState } from "react"
import { Button } from "@/shared/ui/button"
import { api } from "@/shared/api"
import { ArrowLeft, UserPlus } from "lucide-react"

interface RegisterFormProps {
  email: string
  onBack: () => void
  onSuccess: (email: string) => void
}

export function RegisterForm({ email, onBack, onSuccess }: RegisterFormProps) {
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const handleRegister = async () => {
    setIsLoading(true)
    setError(null)
    try {
      await api.post("/auth/register", { email })
      onSuccess(email)
    } catch (err: any) {
      if (err.status === 409) {
        setError("Этот email уже зарегистрирован. Попробуйте войти.")
      } else {
        setError(err.message || "Не удалось отправить запрос на регистрацию. Пожалуйста, попробуйте позже.")
      }
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <div className="space-y-6 animate-in fade-in slide-in-from-right-4 duration-300">
      <div className="text-center space-y-3">
        <div className="inline-flex p-3 bg-primary/10 text-primary rounded-full ring-4 ring-primary/5">
          <UserPlus className="h-6 w-6 animate-pulse" />
        </div>
        <div className="space-y-1">
          <p className="text-sm font-medium text-muted-foreground">Создание нового аккаунта</p>
          <h3 className="text-lg font-semibold text-foreground tracking-tight px-4 leading-snug">
            Создать запрос на регистрацию аккаунта по email:
          </h3>
          <div className="inline-block mt-1 font-mono text-sm bg-accent/30 text-accent-foreground px-2.5 py-1 rounded border border-border/40 font-medium">
            {email}
          </div>
        </div>
      </div>

      {error && (
        <div className="text-sm text-destructive font-medium bg-destructive/10 p-3 rounded-lg border border-destructive/20 animate-in fade-in slide-in-from-top-1">
          {error}
        </div>
      )}

      <div className="flex gap-4 pt-2">
        <Button
          variant="outline"
          type="button"
          className="flex-1 h-11 border-border/40 hover:bg-primary/5 font-medium transition-all"
          onClick={onBack}
          disabled={isLoading}
        >
          <ArrowLeft className="h-4 w-4 mr-2" />
          Назад
        </Button>
        <Button
          type="button"
          className="flex-1 h-11 font-semibold text-base transition-all bg-primary text-primary-foreground hover:bg-primary/90 hover:shadow-[0_0_20px_rgba(var(--primary),0.3)] active:scale-[0.98]"
          onClick={handleRegister}
          disabled={isLoading}
        >
          {isLoading ? "Отправка..." : "Продолжить"}
        </Button>
      </div>
    </div>
  )
}
