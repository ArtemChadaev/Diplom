import { ChevronUp, ChevronDown, ChevronsUpDown } from "lucide-react"
import { Card, CardContent, CardHeader, CardTitle, CardFooter } from "@/components/ui/card"
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table"
import { Badge } from "@/components/ui/badge"
import { Pagination, PaginationContent, PaginationItem, PaginationLink, PaginationNext, PaginationPrevious } from "@/components/ui/pagination"
import { ToggleGroup, ToggleGroupItem } from "@/components/ui/toggle-group"
import { cn } from "@/lib/utils"
import { Medicament, SortOrder } from "../types"
import { TAGS } from "../constants"
import { ExpiryStatus } from "./expiry-status"

interface ResultsTableProps {
  data: Medicament[]
  sortBy: keyof Medicament | ""
  sortOrder: SortOrder
  onSort: (field: keyof Medicament) => void
  pageSize: string
  onPageSizeChange: (v: string) => void
}

export function ResultsTable({
  data,
  sortBy,
  sortOrder,
  onSort,
  pageSize,
  onPageSizeChange
}: ResultsTableProps) {
  const renderSortIcon = (field: keyof Medicament) => {
    if (sortBy !== field || sortOrder === "none") return <ChevronsUpDown className="ml-2 h-4 w-4 shrink-0 opacity-50" />
    if (sortOrder === "asc") return <ChevronUp className="ml-2 h-4 w-4 shrink-0 text-primary" />
    return <ChevronDown className="ml-2 h-4 w-4 shrink-0 text-primary" />
  }

  return (
    <Card className="overflow-hidden shadow-sm">
      <CardHeader className="flex flex-row items-center justify-between border-b border-border/50 pb-4">
        <div className="flex items-center gap-4">
          <CardTitle className="text-xl font-semibold">Список препаратов</CardTitle>
          <ToggleGroup type="single" value={pageSize} onValueChange={(v) => v && onPageSizeChange(v)} className="bg-muted p-1 rounded-full h-10">
            {["5", "10", "15", "20"].map(v => (
              <ToggleGroupItem key={v} value={v} className="px-4 text-xs h-8 rounded-full data-[state=on]:bg-background data-[state=on]:shadow-sm">{v}</ToggleGroupItem>
            ))}
          </ToggleGroup>
        </div>
        <span className="text-xs font-bold uppercase text-muted-foreground">Обновлено 5 мин. назад</span>
      </CardHeader>

      <CardContent className="p-0">
        <Table>
          <TableHeader className="bg-muted/30">
            <TableRow className="hover:bg-transparent">
              <TableHead className="px-6 py-4 text-xs font-bold uppercase tracking-wider text-muted-foreground h-auto cursor-pointer select-none" onClick={() => onSort("name")}>
                <div className="flex items-center">Название {renderSortIcon("name")}</div>
              </TableHead>
              <TableHead className="px-6 py-4 text-xs font-bold uppercase tracking-wider text-muted-foreground h-auto cursor-pointer select-none" onClick={() => onSort("category")}>
                <div className="flex items-center">Категория {renderSortIcon("category")}</div>
              </TableHead>
              <TableHead className="px-6 py-4 text-xs font-bold uppercase tracking-wider text-muted-foreground h-auto">Теги</TableHead>
              <TableHead className="px-6 py-4 text-xs font-bold uppercase tracking-wider text-muted-foreground h-auto cursor-pointer select-none" onClick={() => onSort("warehouse")}>
                <div className="flex items-center">Склад {renderSortIcon("warehouse")}</div>
              </TableHead>
              <TableHead className="px-6 py-4 text-xs font-bold uppercase tracking-wider text-muted-foreground h-auto cursor-pointer select-none" onClick={() => onSort("quantity")}>
                <div className="flex items-center">Количество {renderSortIcon("quantity")}</div>
              </TableHead>
              <TableHead className="px-6 py-4 text-xs font-bold uppercase tracking-wider text-muted-foreground h-auto cursor-pointer select-none" onClick={() => onSort("arrivalDate")}>
                <div className="flex items-center">Дата привоза {renderSortIcon("arrivalDate")}</div>
              </TableHead>
              <TableHead className="px-6 py-4 text-xs font-bold uppercase tracking-wider text-muted-foreground h-auto cursor-pointer select-none" onClick={() => onSort("expiryDate")}>
                <div className="flex items-center">Срок годности {renderSortIcon("expiryDate")}</div>
              </TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {data.map((item) => (
              <TableRow key={item.id} className="group">
                <TableCell className="px-6 py-5">
                  <div className="flex items-center gap-3">
                    <div className={cn("w-1.5 h-8 rounded-full shrink-0", item.tags.includes("critical") ? "bg-destructive" : item.tags.includes("warning") ? "bg-orange-400" : "bg-secondary")}></div>
                    <div>
                      <p className="text-sm font-semibold text-primary">{item.name}</p>
                      <p className="text-xs text-muted-foreground">Batch: {item.batch}</p>
                    </div>
                  </div>
                </TableCell>
                <TableCell className="px-6 py-5">
                  <Badge variant="outline" className="bg-secondary/10 text-secondary border-transparent font-bold tracking-tight">{item.category}</Badge>
                </TableCell>
                <TableCell className="px-6 py-5">
                  <div className="flex flex-wrap gap-1">
                    {item.tags.map(tagValue => {
                      const tag = TAGS.find(t => t.value === tagValue)
                      return <Badge key={tagValue} variant="secondary" className="bg-blue/10 text-blue hover:bg-blue/20 border-transparent text-[10px] uppercase tracking-wider">{tag?.label || tagValue}</Badge>
                    })}
                  </div>
                </TableCell>
                <TableCell className="px-6 py-5">
                  <p className="text-sm font-medium">{item.warehouse}</p>
                  <p className="text-xs text-muted-foreground">{item.sector}</p>
                </TableCell>
                <TableCell className="px-6 py-5">
                  <p className="text-sm font-bold">{item.quantity} <span className="font-normal text-muted-foreground">уп.</span></p>
                </TableCell>
                <TableCell className="px-6 py-5">
                  <p className="text-xs font-medium text-muted-foreground">{item.arrivalDate.toLocaleDateString("ru-RU")}</p>
                </TableCell>
                <TableCell className="px-6 py-5">
                  <ExpiryStatus productionDate={item.productionDate} expiryDate={item.expiryDate} />
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </CardContent>
      
      <CardFooter className="py-4 border-t border-border/50 flex items-center justify-between">
        <p className="text-sm text-muted-foreground">Показано 1-{data.length} из 1,284 результатов</p>
        <Pagination className="w-auto mx-0">
          <PaginationContent>
            <PaginationItem><PaginationPrevious href="#" onClick={(e) => e.preventDefault()} /></PaginationItem>
            <PaginationItem><PaginationLink href="#" isActive onClick={(e) => e.preventDefault()}>1</PaginationLink></PaginationItem>
            <PaginationItem><PaginationLink href="#" onClick={(e) => e.preventDefault()}>2</PaginationLink></PaginationItem>
            <PaginationItem><PaginationLink href="#" onClick={(e) => e.preventDefault()}>3</PaginationLink></PaginationItem>
            <PaginationItem><PaginationNext href="#" onClick={(e) => e.preventDefault()} /></PaginationItem>
          </PaginationContent>
        </Pagination>
      </CardFooter>
    </Card>
  )
}
