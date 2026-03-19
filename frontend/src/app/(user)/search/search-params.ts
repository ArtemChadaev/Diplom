// src/app/search/search-params.ts
import {
  parseAsString,
  parseAsArrayOf,
  parseAsIsoDate,
  parseAsInteger,
} from 'nuqs' // Для клиента используем основной вход

export const searchParamsParsers = {
  q: parseAsString.withDefault(""),
  // Явно указываем пустой массив в дефолте
  categories: parseAsArrayOf(parseAsString).withDefault([]),
  warehouses: parseAsArrayOf(parseAsString).withDefault([]),
  tags: parseAsArrayOf(parseAsString).withDefault([]),

  aDate: parseAsIsoDate,
  date: parseAsIsoDate,
  days: parseAsString.withDefault(""),
  sortBy: parseAsString.withDefault("arrivalDate"),
  sortOrder: parseAsString.withDefault("desc"),
  limit: parseAsInteger.withDefault(10),
}