import { useAuthLogin } from "../model/use-auth-login"
import { Button } from "@/shared/ui/button"
import { Input } from "@/shared/ui/input"
import { ButtonGroup } from "@/shared/ui/button-group"
import { Separator } from "@/shared/ui/separator"
import { ArrowRight, Mail } from "lucide-react"

interface LoginFormProps {
  onFound: (email: string) => void
  onNotFound: (email: string) => void
  defaultEmail?: string
}

export function LoginForm({ onFound, onNotFound, defaultEmail = "" }: LoginFormProps) {
  const { form, onSubmit, isLoading, error } = useAuthLogin({
    onFound,
    onNotFound,
    defaultEmail,
  })

  const handleGoogleLogin = () => {
    window.location.href = "/api/v1/auth/google"
  }

  return (
    <div className="space-y-6 animate-in fade-in slide-in-from-left-4 duration-300">
      <form onSubmit={onSubmit} className="space-y-4">
        <div className="space-y-2">
          <label htmlFor="email-input" className="text-sm font-medium text-foreground">
            Электронная почта
          </label>
          <div className="relative group">
            <ButtonGroup className="w-full">
              <div className="relative flex-1">
                <Mail className="absolute left-3.5 top-3.5 h-4 w-4 text-muted-foreground transition-colors group-focus-within:text-primary z-10" />
                <Input
                  id="email-input"
                  placeholder="name@example.com"
                  type="email"
                  className="pl-10 h-11 bg-background/50 border-border/40 focus-visible:border-primary/50 transition-all text-base rounded-none w-full"
                  disabled={isLoading}
                  {...form.register("email")}
                />
              </div>
              <Button 
                type="submit" 
                size="icon"
                className="h-11 w-12 rounded-none bg-primary hover:bg-primary/90 text-primary-foreground font-semibold shrink-0 transition-all active:scale-[0.98]"
                disabled={isLoading}
              >
                <ArrowRight className="h-5 w-5" />
              </Button>
            </ButtonGroup>
          </div>
          {form.formState.errors.email && (
            <p className="text-sm text-destructive font-medium animate-in fade-in slide-in-from-top-1">
              {form.formState.errors.email.message}
            </p>
          )}
          {error && (
            <p className="text-sm text-destructive font-medium animate-in fade-in slide-in-from-top-1 bg-destructive/10 p-3 rounded-lg border border-destructive/20">
              {error}
            </p>
          )}
        </div>
      </form>

      <div className="relative flex items-center justify-center my-4">
        <Separator className="w-full absolute" />
        <span className="relative bg-card px-3 text-xs text-muted-foreground font-medium uppercase tracking-wider">
          Или
        </span>
      </div>

      <Button
        variant="outline"
        type="button"
        className="w-full h-11 border-border/40 hover:bg-primary/5 font-medium transition-all group overflow-hidden relative active:scale-[0.98]"
        onClick={handleGoogleLogin}
        disabled={isLoading}
      >
        <div className="absolute inset-0 bg-gradient-to-r from-transparent via-primary/5 to-transparent -translate-x-full group-hover:translate-x-full transition-transform duration-1000" />
        <svg className="w-5 h-5 mr-2 inline-block transition-transform group-hover:scale-110" viewBox="0 0 24 24">
          <path d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z" fill="#4285F4" />
          <path d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z" fill="#34A853" />
          <path d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l3.66-2.84z" fill="#FBBC05" />
          <path d="M12 5.38c1.62 0 3.06.56 4.21 1.66l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z" fill="#EA4335" />
        </svg>
        Войти через Google
      </Button>
    </div>
  )
}
