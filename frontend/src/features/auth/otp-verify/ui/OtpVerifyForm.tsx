import { useOtpVerify } from "../model/use-otp-verify"
import { Button } from "@/shared/ui/button"
import { 
  InputOTP, 
  InputOTPGroup, 
  InputOTPSlot, 
  InputOTPSeparator 
} from "@/shared/ui/input-otp"
import { ArrowLeft, RefreshCw, KeyRound } from "lucide-react"

interface OtpVerifyFormProps {
  email: string
  onBack: () => void
}

export function OtpVerifyForm({ email, onBack }: OtpVerifyFormProps) {
  const {
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
  } = useOtpVerify({ email })

  const handleOTPChange = (value: string) => {
    setCode(value)
    if (value.length === 6) {
      verify(value)
    }
  }

  return (
    <div className="space-y-6 animate-in fade-in slide-in-from-right-4 duration-300">
      <div className="text-center space-y-3">
        <div className="inline-flex p-3 bg-primary/10 text-primary rounded-full ring-4 ring-primary/5">
          <KeyRound className="h-6 w-6" />
        </div>
        <div className="space-y-1">
          <h3 className="text-lg font-semibold text-foreground tracking-tight">Подтверждение почты</h3>
          <p className="text-sm text-muted-foreground px-4 leading-relaxed">
            Мы отправили 6-значный код подтверждения на адрес:
          </p>
          <div className="inline-block mt-1 font-mono text-sm bg-accent/30 text-accent-foreground px-2.5 py-1 rounded border border-border/40 font-medium">
            {email}
          </div>
        </div>
      </div>

      <div className="flex flex-col items-center justify-center space-y-4">
        <InputOTP
          maxLength={6}
          value={code}
          onChange={handleOTPChange}
          disabled={isLoading}
          containerClassName="justify-center"
        >
          <InputOTPGroup>
            <InputOTPSlot index={0} />
            <InputOTPSlot index={1} />
            <InputOTPSlot index={2} />
          </InputOTPGroup>
          <InputOTPSeparator />
          <InputOTPGroup>
            <InputOTPSlot index={3} />
            <InputOTPSlot index={4} />
            <InputOTPSlot index={5} />
          </InputOTPGroup>
        </InputOTP>
        
        {isLoading && (
          <p className="text-sm text-muted-foreground animate-pulse font-medium">
            Проверка кода...
          </p>
        )}
      </div>

      {error && (
        <div className="text-sm text-destructive font-medium bg-destructive/10 p-3 rounded-lg border border-destructive/20 animate-in fade-in slide-in-from-top-1">
          {error}
        </div>
      )}

      {resendMessage && !error && (
        <div className="text-sm text-emerald-600 dark:text-emerald-400 font-medium bg-emerald-500/10 p-3 rounded-lg border border-emerald-500/20 animate-in fade-in slide-in-from-top-1">
          {resendMessage}
        </div>
      )}

      <div className="space-y-4 pt-2">
        <div className="text-center text-sm">
          {secondsLeft > 0 ? (
            <p className="text-muted-foreground">
              Запросить код повторно через <span className="font-mono font-semibold text-foreground bg-accent/30 px-1.5 py-0.5 rounded border border-border/40">{formatTime(secondsLeft)}</span>
            </p>
          ) : (
            <Button
              variant="link"
              type="button"
              className="text-primary hover:text-primary/80 font-medium p-0 h-auto"
              onClick={resend}
              disabled={resendLoading || isLoading}
            >
              <RefreshCw className={`h-4 w-4 mr-2 ${resendLoading ? "animate-spin" : ""}`} />
              Отправить код повторно
            </Button>
          )}
        </div>

        <Button
          variant="outline"
          type="button"
          className="w-full h-11 border-border/40 hover:bg-primary/5 font-medium transition-all"
          onClick={onBack}
          disabled={isLoading}
        >
          <ArrowLeft className="h-4 w-4 mr-2" />
          Назад к вводу email
        </Button>
      </div>
    </div>
  )
}
