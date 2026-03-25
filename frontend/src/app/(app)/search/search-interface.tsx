"use client"

import { useMemo, useCallback } from "react"
import { useQueryStates } from 'nuqs'
import { searchParamsParsers } from "./search-params"
import { SearchHeader } from "./components/search-header"
import { FilterSection } from "./components/filter-section"
import { ResultsTable } from "./components/results-table"
import { Medicament, SortOrder } from "./types"
import { MOCK_DATA } from "./data"

export function SearchInterface() {
  // 1. Синхронизация с URL через nuqs
  const [params, setParams] = useQueryStates(searchParamsParsers, {
    shallow: true, // Обновляет URL без перезагрузки страницы (SPA style)
    history: 'replace' // Не забивает историю браузера каждым символом поиска
  })

  // 2. Логика обработки данных (Фильтрация + Сортировка)
  const processedData = useMemo(() => {
    let result = [...MOCK_DATA]

    // --- ФИЛЬТРАЦИЯ ---
    if (params.q) {
      const query = params.q.toLowerCase()
      result = result.filter(item => item.name.toLowerCase().includes(query))
    }

    if (params.categories.length > 0) {
      result = result.filter(item => params.categories.includes(item.category))
    }

    if (params.warehouses.length > 0) {
      result = result.filter(item => params.warehouses.includes(item.warehouse))
    }

    // Фильтр по дате прихода (если есть)
    if (params.aDate) {
      result = result.filter(item =>
        new Date(item.arrivalDate).toDateString() === params.aDate!.toDateString()
      )
    }

    // --- СОРТИРОВКА ---
    if (params.sortOrder !== "none") {
      result.sort((a, b) => {
        const aValue = a[params.sortBy as keyof Medicament]
        const bValue = b[params.sortBy as keyof Medicament]

        const modifier = params.sortOrder === "asc" ? 1 : -1

        if (aValue instanceof Date && bValue instanceof Date) {
          return (aValue.getTime() - bValue.getTime()) * modifier
        }
        return (aValue > bValue ? 1 : -1) * modifier
      })
    }

    // --- ПАГИНАЦИЯ (лимит) ---
    return result.slice(0, params.limit)
  }, [params])

  // 3. Обработчики действий
  const handleSort = useCallback((field: keyof Medicament) => {
    setParams(prev => {
      const isCurrentField = prev.sortBy === field
      const nextOrder: SortOrder =
        isCurrentField && prev.sortOrder === "asc" ? "desc" :
          isCurrentField && prev.sortOrder === "desc" ? "none" : "asc"

      return {
        sortBy: nextOrder === "none" ? "arrivalDate" : field,
        sortOrder: nextOrder
      }
    })
  }, [setParams])

  const resetFilters = useCallback(() => {
    setParams(null) // nuqs сбросит всё в default значения из конфига
  }, [setParams])

  return (
    <div className="flex flex-col gap-8 w-full animate-in fade-in duration-500">
      <SearchHeader />

      {/* Передаем params целиком, чтобы не плодить пропсы */}
      <FilterSection
        params={params}
        setParams={setParams}
        resetFilters={resetFilters}
      />

      <div className="rounded-xl border bg-card shadow-sm">
        <ResultsTable
          data={processedData}
          sortBy={params.sortBy as keyof Medicament}
          sortOrder={params.sortOrder as SortOrder}
          onSort={handleSort}
          pageSize={params.limit.toString()}
          onPageSizeChange={(val) => setParams({ limit: parseInt(val) })}
        />
      </div>
    </div>
  )
}