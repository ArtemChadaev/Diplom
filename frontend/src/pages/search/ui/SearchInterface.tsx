import { useState, useMemo, useCallback } from "react"
import { SearchHeader } from "./SearchHeader"
import { FilterSection } from "./FilterSection"
import { ResultsTable } from "./ResultsTable"
import type { Medicament, SortOrder } from "../types"
import { MOCK_DATA } from "../data"

export function SearchInterface() {
  const [params, setParams] = useState({
    q: "",
    categories: [] as string[],
    warehouses: [] as string[],
    tags: [] as string[],
    aDate: null as Date | null,
    days: "",
    sortBy: "arrivalDate",
    sortOrder: "desc" as SortOrder,
    limit: 10
  })

  const processedData = useMemo(() => {
    let result = [...MOCK_DATA]

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

    if (params.aDate) {
      result = result.filter(item =>
        new Date(item.arrivalDate).toDateString() === params.aDate!.toDateString()
      )
    }

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

    return result.slice(0, params.limit)
  }, [params])

  const handleSort = useCallback((field: keyof Medicament) => {
    setParams(prev => {
      const isCurrentField = prev.sortBy === field
      const nextOrder: SortOrder =
        isCurrentField && prev.sortOrder === "asc" ? "desc" :
          isCurrentField && prev.sortOrder === "desc" ? "none" : "asc"

      return {
        ...prev,
        sortBy: nextOrder === "none" ? "arrivalDate" : field,
        sortOrder: nextOrder
      }
    })
  }, [])

  const resetFilters = useCallback(() => {
    setParams({
      q: "",
      categories: [],
      warehouses: [],
      tags: [],
      aDate: null,
      days: "",
      sortBy: "arrivalDate",
      sortOrder: "desc",
      limit: 10
    })
  }, [])

  return (
    <div className="flex flex-col gap-8 w-full animate-in fade-in duration-500">
      <SearchHeader />

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
          onPageSizeChange={(val) => setParams(prev => ({ ...prev, limit: parseInt(val) }))}
        />
      </div>
    </div>
  )
}
