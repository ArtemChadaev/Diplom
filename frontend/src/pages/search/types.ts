export type SortOrder = "asc" | "desc" | "none"

export interface Medicament {
  id: string
  name: string
  batch: string
  category: string
  warehouse: string
  sector: string
  quantity: number
  arrivalDate: Date
  productionDate: Date
  expiryDate: Date
  tags: string[]
}

export interface SearchParams {
  q: string
  categories: string[]
  warehouses: string[]
  tags: string[]
  aDate: Date | null
  days: string
  sortBy: string
  sortOrder: SortOrder
  limit: number
}
