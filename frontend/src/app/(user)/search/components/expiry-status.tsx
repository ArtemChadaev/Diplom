"use client"

import { differenceInDays, formatDistanceStrict, isBefore, addDays } from "date-fns"
import { ru } from "date-fns/locale"
import { Badge } from "@/components/ui/badge"
import { cn } from "@/lib/utils"

interface ExpiryStatusProps {
  productionDate: Date
  expiryDate: Date
  className?: string
}

export function ExpiryStatus({ productionDate, expiryDate, className }: ExpiryStatusProps) {
  const today = new Date()
  const daysLeft = differenceInDays(expiryDate, today)
  const isExpired = isBefore(expiryDate, today)

  let statusColor = "text-secondary"
  let label = "Истек: Нет"

  if (isExpired) {
    statusColor = "text-destructive"
    label = "Истек: Да"
  } else if (daysLeft <= 30) {
    statusColor = "text-destructive"
    label = "Истек: Критично"
  } else if (daysLeft <= 90) {
    statusColor = "text-orange-600"
    label = "Истек: Скоро"
  }

  const formatDate = (date: Date) => date.toLocaleDateString("ru-RU", { day: '2-digit', month: '2-digit', year: 'numeric' })

  return (
    <div className={cn("flex flex-col gap-0.5", className)}>
      <p className="text-[10px] font-medium text-muted-foreground whitespace-nowrap">
        Изг: <span className="text-primary">{formatDate(productionDate)}</span>
      </p>
      <p className="text-[10px] font-medium text-muted-foreground whitespace-nowrap">
        Срок: <span className="text-primary">{formatDate(expiryDate)}</span>
      </p>
      <p className={cn("text-[10px] font-bold uppercase tracking-tight", statusColor)}>
        {label}
      </p>
    </div>
  )
}
