import { zodResolver } from "@hookform/resolvers/zod"
import { useEffect, useState } from "react"
import { useForm } from "react-hook-form"
import { z } from "zod"

import { api } from "@/shared/api"

export const emailSchema = z.object({
  email: z.email({ message: "Введите корректный адрес электронной почты" }),
})

export type EmailFormValues = z.infer<typeof emailSchema>

interface UseAuthLoginProps {
  onFound: (email: string) => void
  onNotFound: (email: string) => void
  defaultEmail?: string
}

export function useAuthLogin({ onFound, onNotFound, defaultEmail = "" }: UseAuthLoginProps) {
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const form = useForm<EmailFormValues>({
    resolver: zodResolver(emailSchema),
    defaultValues: {
      email: defaultEmail,
    },
  })

  // Update default value if it changes
  useEffect(() => {
    if (defaultEmail) {
      form.setValue("email", defaultEmail)
    }
  }, [defaultEmail, form])

  const onSubmit = async (values: EmailFormValues) => {
    setIsLoading(true)
    setError(null)
    try {
      await api.post("/auth/send-code", { email: values.email })
      onFound(values.email)
    } catch (err) {
      const apiErr = err as { status?: number; message?: string } | null
      if (apiErr?.status === 404) {
        onNotFound(values.email)
      } else {
        setError(apiErr?.message ?? "Произошла ошибка при отправке кода. Попробуйте позже.")
      }
    } finally {
      setIsLoading(false)
    }
  }

  return {
    form,
    onSubmit: form.handleSubmit(onSubmit),
    isLoading,
    error,
    setError,
  }
}
