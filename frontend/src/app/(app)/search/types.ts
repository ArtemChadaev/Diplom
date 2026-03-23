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
