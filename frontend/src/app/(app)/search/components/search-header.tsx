import { Button } from "@/components/ui/button"
import { Download } from "lucide-react"

export function SearchHeader() {
  return (
    <header className="flex flex-col md:flex-row md:items-end justify-between gap-6">
      <div>
        <p className="text-xs font-bold uppercase tracking-widest text-secondary mb-2">Логистика и учет</p>
        <h1 className="text-4xl font-bold text-primary tracking-tight">Склад: Мониторинг и Поиск</h1>
      </div>
      <div className="flex gap-3">
        <Button variant="secondary" className="px-6 py-5 rounded-xl text-sm font-medium h-auto flex items-center gap-2">
          <Download className="h-5 w-5" />
          Экспорт отчета
        </Button>
      </div>
    </header>
  )
}
