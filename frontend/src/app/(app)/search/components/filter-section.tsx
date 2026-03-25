"use client"
import * as React from "react"
import { Input } from "@/components/ui/input"
import { Button } from "@/components/ui/button"
import { Card, CardContent } from "@/components/ui/card"
import { MultiSelect } from "@/components/ui/multi-select"
import { DatePicker } from "@/components/ui/date-picker"
import { ToggleGroup, ToggleGroupItem } from "@/components/ui/toggle-group"
import { CATEGORIES, WAREHOUSES, TAGS } from "../constants"

// Импортируем типы из нашего конфига параметров
import { searchParamsParsers } from "../search-params"

interface FilterSectionProps {
  // Используем Inference для типов, чтобы не переписывать их вручную
  params: any
  setParams: (state: any) => void
  resetFilters: () => void
}

export function FilterSection({
                                params,
                                setParams,
                                resetFilters
                              }: FilterSectionProps) {

  // Функция для удобного обновления одного поля
  const updateField = (field: string, value: any) => {
    // Если значение пустое, передаем null, чтобы nuqs удалил ключ из URL
    const finalValue = (Array.isArray(value) && value.length === 0) || value === "" ? null : value
    setParams({ [field]: finalValue })
  }

  return (
    <Card className="shadow-sm border-none bg-muted/30">
      <CardContent className="p-8">
        <div className="grid grid-cols-1 md:grid-cols-3 gap-y-8 gap-x-6">

          {/* Название */}
          <div>
            <label className="block text-xs font-bold uppercase text-muted-foreground mb-2">Название</label>
            <Input
              className="w-full bg-background min-h-12 border-muted-foreground/20 focus-visible:ring-secondary"
              placeholder="Введите название препарата..."
              value={params.q ?? ""} // Защита от null
              onChange={(e) => updateField("q", e.target.value)}
            />
          </div>

          {/* Категория */}
          <div>
            <label className="block text-xs font-bold uppercase text-muted-foreground mb-2">Категория</label>
            <MultiSelect
              options={CATEGORIES}
              selected={params.categories ?? []} // ГЛАВНАЯ ЗАЩИТА ОТ ОШИБКИ .length
              onChange={(vals) => updateField("categories", vals)}
              placeholder="Все категории"
            />
          </div>

          {/* Склад */}
          <div>
            <label className="block text-xs font-bold uppercase text-muted-foreground mb-2">Склад</label>
            <MultiSelect
              options={WAREHOUSES}
              selected={params.warehouses ?? []} // ЗАЩИТА
              onChange={(vals) => updateField("warehouses", vals)}
              placeholder="Все склады"
            />
          </div>

          {/* Теги */}
          <div>
            <label className="block text-xs font-bold uppercase text-muted-foreground mb-2">Теги</label>
            <MultiSelect
              options={TAGS}
              selected={params.tags ?? []} // ЗАЩИТА
              onChange={(vals) => updateField("tags", vals)}
              placeholder="Выберите теги"
              badgeClassName="bg-blue-500/10 text-blue-600 hover:bg-blue-500/20 border-blue-200"
            />
          </div>

          {/* Дата привоза */}
          <div>
            <label className="block text-xs font-bold uppercase text-muted-foreground mb-2">Дата привоза</label>
            <DatePicker
              date={params.aDate ?? undefined}
              setDate={(date) => updateField("aDate", date)}
              placeholder="Выберите дату"
            />
          </div>

          {/* Осталось дней */}
          <div>
            <label className="block text-xs font-bold uppercase text-muted-foreground mb-2">Осталось дней</label>
            <Input
              type="number"
              className="w-full bg-background min-h-12 border-muted-foreground/20 focus-visible:ring-secondary"
              placeholder="Количество дней..."
              value={params.days ?? ""}
              onChange={(e) => updateField("days", e.target.value)}
            />
          </div>
        </div>

        <div className="mt-8 flex flex-col md:flex-row md:items-center justify-between gap-4">
          <div className="flex items-center gap-4 grow">
            <ToggleGroup
              type="single"
              value={params.sortBy}
              onValueChange={(val) => val && updateField("sortBy", val)}
              variant="outline"
              className="justify-start h-10 bg-background border border-muted-foreground/20 p-1 rounded-full w-fit"
            >
              <ToggleGroupItem value="name" className="text-[10px] uppercase font-bold tracking-tight px-4 h-8 rounded-full data-[state=on]:bg-secondary data-[state=on]:text-secondary-foreground border-transparent">Имя</ToggleGroupItem>
              <ToggleGroupItem value="category" className="text-[10px] uppercase font-bold tracking-tight px-4 h-8 rounded-full data-[state=on]:bg-secondary data-[state=on]:text-secondary-foreground border-transparent">Категория</ToggleGroupItem>
              <ToggleGroupItem value="warehouse" className="text-[10px] uppercase font-bold tracking-tight px-4 h-8 rounded-full data-[state=on]:bg-secondary data-[state=on]:text-secondary-foreground border-transparent">Склад</ToggleGroupItem>
              <ToggleGroupItem value="quantity" className="text-[10px] uppercase font-bold tracking-tight px-4 h-8 rounded-full data-[state=on]:bg-secondary data-[state=on]:text-secondary-foreground border-transparent">Кол-во</ToggleGroupItem>
            </ToggleGroup>
          </div>

          <div className="flex justify-end gap-4 min-w-fit">
            <Button variant="ghost" className="text-muted-foreground" onClick={resetFilters}>
              Сбросить всё
            </Button>
            {/* Кнопка "Применить" теперь по сути декоративная, 
                так как nuqs обновляет всё на лету. Но оставим её для UX. */}
            <Button
              className="bg-secondary text-secondary-foreground hover:bg-secondary/90 px-8"
              onClick={() => {/* Можно добавить принудительный рефетч данных */}}
            >
              Применить фильтры
            </Button>
          </div>
        </div>
      </CardContent>
    </Card>
  )
}