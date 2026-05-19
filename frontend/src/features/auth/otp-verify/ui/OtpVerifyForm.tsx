import { ArrowLeft, RefreshCw } from "lucide-react"

import { Button } from "@/shared/ui/button"
import { 
  InputOTP, 
  InputOTPGroup, 
  InputOTPSeparator, 
  InputOTPSlot, 
} from "@/shared/ui/input-otp"

import { useOtpVerify } from "../model/use-otp-verify"

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
    verify,
    resend,
    formatTime,
  } = useOtpVerify({ email })

  const handleOTPChange = (value: string) => {
    setCode(value)
    if (value.length === 6) {
      void verify(value)
    }
  }

  return (
    <div className="space-y-6 animate-in fade-in slide-in-from-right-4 duration-300">

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
        <p className="text-sm text-destructive font-medium animate-in fade-in slide-in-from-top-1 text-center">
          {error}
        </p>
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
              onClick={() => { void resend(); }}
              disabled={resendLoading || isLoading}
            >
              <RefreshCw className={`h-4 w-4 mr-2 ${resendLoading ? "animate-spin" : ""}`} />
              Отправить код повторно
            </Button>
          )}
        </div>

        <p className="text-xs text-muted-foreground/70 text-center">
          Если письмо не пришло, проверьте папку «Спам»
        </p>

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
