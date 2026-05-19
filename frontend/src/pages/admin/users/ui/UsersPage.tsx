import { useEffect, useState } from "react"

import { UserActionsDropdown, UserRoleBadge } from "@/features/admin/user-actions"

import { useAdminUsersStore, type UserProfile } from "@/entities/user"

import { Avatar, AvatarFallback, AvatarImage } from "@/shared/ui/avatar"
import { Badge } from "@/shared/ui/badge"
import {
  Pagination,
  PaginationContent,
  PaginationItem,
  PaginationLink,
  PaginationNext,
  PaginationPrevious,
} from "@/shared/ui/pagination"
import { Skeleton } from "@/shared/ui/skeleton"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/shared/ui/table"

export function UsersPage() {
  const { users, isLoading, error, fetchUsers } = useAdminUsersStore()

  // Состояние страниц для каждой из трех таблиц
  const [unregisteredPage, setUnregisteredPage] = useState(1)
  const [activePage, setActivePage] = useState(1)
  const [blockedPage, setBlockedPage] = useState(1)

  useEffect(() => {
    void fetchUsers()
  }, [fetchUsers])

  // Фильтрация пользователей по группам
  // 1. Пользователи без профиля (незарегистрированные и не заблокированные)
  const unregisteredUsers = users.filter(
    (user) => (!user.full_name || !user.employee_code) && !user.is_blocked
  )
  
  // 2. Активные пользователи (с профилем и не заблокированные)
  const activeUsers = users.filter(
    (user) => user.full_name && user.employee_code && !user.is_blocked
  )
  
  // 3. Заблокированные пользователи
  const blockedUsers = users.filter(
    (user) => user.is_blocked
  )

  const itemsPerPage = 10

  const formatPageRange = (total: number, currentPage: number) => {
    if (total === 0) return "0/0"
    const start = (currentPage - 1) * itemsPerPage + 1
    const end = Math.min(currentPage * itemsPerPage, total)
    return `${start}-${end}/${total}`
  }

  // Получить данные для текущей страницы
  const getPagedData = (data: UserProfile[], page: number) => {
    return data.slice((page - 1) * itemsPerPage, page * itemsPerPage)
  }

  // Общий рендер пагинации для таблиц
  const renderPagination = (
    totalItems: number,
    currentPage: number,
    setCurrentPage: (page: number) => void
  ) => {
    if (totalItems <= itemsPerPage) return null
    const totalPages = Math.ceil(totalItems / itemsPerPage)
    
    const pages = []
    for (let i = 1; i <= totalPages; i++) {
      pages.push(i)
    }

    return (
      <div className="py-3 px-4 border-t border-border/50 bg-muted/10 flex justify-center">
        <Pagination>
          <PaginationContent>
            <PaginationItem>
              <PaginationPrevious
                href="#"
                onClick={(e) => {
                  e.preventDefault()
                  if (currentPage > 1) setCurrentPage(currentPage - 1)
                }}
                className={currentPage === 1 ? "pointer-events-none opacity-50" : "cursor-pointer"}
                text="Назад"
              />
            </PaginationItem>
            {pages.map((p) => (
              <PaginationItem key={p}>
                <PaginationLink
                  href="#"
                  isActive={p === currentPage}
                  onClick={(e) => {
                    e.preventDefault()
                    setCurrentPage(p)
                  }}
                  className="cursor-pointer"
                >
                  {p}
                </PaginationLink>
              </PaginationItem>
            ))}
            <PaginationItem>
              <PaginationNext
                href="#"
                onClick={(e) => {
                  e.preventDefault()
                  if (currentPage < totalPages) setCurrentPage(currentPage + 1)
                }}
                className={currentPage === totalPages ? "pointer-events-none opacity-50" : "cursor-pointer"}
                text="Вперед"
              />
            </PaginationItem>
          </PaginationContent>
        </Pagination>
      </div>
    )
  }

  // Рендер ячейки профиля с центрированием
  const renderAvatarCell = (user: UserProfile) => {
    const initials = user.full_name
      ? user.full_name
          .split(" ")
          .filter(Boolean)
          .map((n) => n[0])
          .join("")
      : "?"

    return (
      <TableCell className="w-[80px]">
        <div className="flex justify-center items-center">
          <Avatar className="h-10 w-10 border border-border/10 shadow-sm" size="lg">
            {user.avatar_url && (
              <AvatarImage src={user.avatar_url} alt={user.full_name ?? ""} />
            )}
            <AvatarFallback className="bg-primary text-primary-foreground font-medium flex items-center justify-center">
              {initials}
            </AvatarFallback>
          </Avatar>
        </div>
      </TableCell>
    )
  }

  // Рендер информации о пользователе и его бейджей
  const renderInfoCell = (user: UserProfile) => {
    const isUnregistered = !user.full_name || !user.employee_code
    
    return (
      <TableCell>
        <div className="flex flex-col gap-1.5">
          {isUnregistered ? (
            <span className="text-foreground">—</span>
          ) : (
            <div className="flex items-center gap-2">
              <span className="font-semibold text-foreground text-sm">
                {user.full_name}
              </span>
              {user.employee_code && (
                <span className="text-[10px] font-mono bg-muted text-muted-foreground px-1.5 py-0.5 border border-border/50">
                  {user.employee_code}
                </span>
              )}
            </div>
          )}
          
          <div className="flex gap-2 items-center flex-wrap">
            {user.telegram_handle && (
              <span className="text-xs text-muted-foreground">{user.telegram_handle}</span>
            )}
            <span className="text-xs text-muted-foreground/60">({user.email})</span>
            <UserRoleBadge role={user.role} isBlocked={user.is_blocked} />
            {user.ns_pv_access && (
              <Badge variant="outline" className="text-[9px] h-4.5 px-1.5 border-destructive/20 text-destructive bg-destructive/5 font-bold uppercase tracking-wider">
                НС/ПВ
              </Badge>
            )}
            {user.ukep_bound && (
              <Badge variant="outline" className="text-[9px] h-4.5 px-1.5 border-primary/20 text-primary bg-primary/5 font-bold uppercase tracking-wider">
                УКЭП
              </Badge>
            )}
          </div>
        </div>
      </TableCell>
    )
  }

  if (error) {
    return (
      <div className="p-6 text-center space-y-4">
        <div className="text-destructive font-semibold text-lg">{error}</div>
        <button
          onClick={() => fetchUsers()}
          className="bg-primary text-primary-foreground px-4 py-2 hover:bg-primary/90 text-xs font-semibold cursor-pointer"
        >
          Попробовать снова
        </button>
      </div>
    )
  }

  return (
    <div className="space-y-8 pb-10">
      <div>
        <h1 className="text-3xl font-bold tracking-tight">Управление пользователями</h1>
        <p className="text-muted-foreground mt-1">
          Просмотр и управление учетными записями сотрудников
        </p>
      </div>

      {isLoading ? (
        <div className="space-y-6">
          {[1, 2].map((i) => (
            <div key={i} className="space-y-3">
              <div className="flex justify-between items-center">
                <Skeleton className="h-6 w-48" />
                <Skeleton className="h-4 w-32" />
              </div>
              <div className="border border-border/50 bg-card rounded-none divide-y divide-border/50">
                {[1, 2, 3].map((j) => (
                  <div key={j} className="p-4 flex items-center justify-between">
                    <div className="flex items-center gap-4">
                      <Skeleton className="h-10 w-10 rounded-full" />
                      <div className="space-y-2">
                        <Skeleton className="h-4 w-32" />
                        <Skeleton className="h-3 w-48" />
                      </div>
                    </div>
                    <Skeleton className="h-8 w-24" />
                  </div>
                ))}
              </div>
            </div>
          ))}
        </div>
      ) : (
        <div className="space-y-8">
          
          {/* ТАБЛИЦА 1: ПОЛЬЗОВАТЕЛИ БЕЗ ПРОФИЛЯ */}
          {unregisteredUsers.length > 0 && (
            <div className="space-y-3 animate-in fade-in duration-300">
              <div className="flex items-center justify-between">
                <h2 className="text-lg font-semibold text-foreground">
                  Пользователи без профиля
                </h2>
                <span className="text-xs text-muted-foreground">
                  {formatPageRange(unregisteredUsers.length, unregisteredPage)}
                </span>
              </div>
              <div className="rounded-none border border-border/50 bg-card overflow-hidden shadow-sm">
                <Table>
                  <TableHeader className="bg-muted/30">
                    <TableRow>
                      <TableHead className="w-[80px] text-center">Профиль</TableHead>
                      <TableHead>Данные учетной записи</TableHead>
                      <TableHead>Должность в системе / Отдел</TableHead>
                      <TableHead className="text-right w-[120px]">Действия</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {getPagedData(unregisteredUsers, unregisteredPage).map((user) => (
                      <TableRow key={user.id} className="hover:bg-muted/20 transition-colors">
                        {renderAvatarCell(user)}
                        {renderInfoCell(user)}
                        <TableCell>
                          <div className="flex flex-col gap-0.5">
                            <span className="text-xs font-semibold text-muted-foreground">
                              Нет профиля сотрудника
                            </span>
                            <span className="text-[10px] text-muted-foreground/60 italic">
                              Необходимо ввести информацию
                            </span>
                          </div>
                        </TableCell>
                        <TableCell className="text-right">
                          <UserActionsDropdown user={user} />
                        </TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
                {renderPagination(unregisteredUsers.length, unregisteredPage, setUnregisteredPage)}
              </div>
            </div>
          )}

          {/* ТАБЛИЦА 2: АКТИВНЫЕ СОТРУДНИКИ */}
          <div className="space-y-3 animate-in fade-in duration-300">
            <div className="flex items-center justify-between">
              <h2 className="text-lg font-semibold text-foreground">Активные сотрудники</h2>
              <span className="text-xs text-muted-foreground">
                {formatPageRange(activeUsers.length, activePage)}
              </span>
            </div>
            <div className="rounded-none border border-border/50 bg-card overflow-hidden shadow-sm">
              {activeUsers.length === 0 ? (
                <div className="p-8 text-center text-muted-foreground">
                  Нет активных сотрудников в системе.
                </div>
              ) : (
                <>
                  <Table>
                    <TableHeader className="bg-muted/30">
                      <TableRow>
                        <TableHead className="w-[80px] text-center">Профиль</TableHead>
                        <TableHead>Сотрудник</TableHead>
                        <TableHead>Должность / Отдел</TableHead>
                        <TableHead className="text-right w-[120px]">Действия</TableHead>
                      </TableRow>
                    </TableHeader>
                    <TableBody>
                      {getPagedData(activeUsers, activePage).map((user) => (
                        <TableRow key={user.id} className="hover:bg-muted/20 transition-colors">
                          {renderAvatarCell(user)}
                          {renderInfoCell(user)}
                          <TableCell>
                            <div className="flex flex-col">
                              <span className="text-sm font-medium">{user.position ?? "—"}</span>
                              <span className="text-xs text-muted-foreground">{user.department ?? "—"}</span>
                            </div>
                          </TableCell>
                          <TableCell className="text-right">
                            <UserActionsDropdown user={user} />
                          </TableCell>
                        </TableRow>
                      ))}
                    </TableBody>
                  </Table>
                  {renderPagination(activeUsers.length, activePage, setActivePage)}
                </>
              )}
            </div>
          </div>

          {/* ТАБЛИЦА 3: ЗАБЛОКИРОВАННЫЕ СОТРУДНИКИ */}
          {blockedUsers.length > 0 && (
            <div className="space-y-3 animate-in fade-in duration-300">
              <div className="flex items-center justify-between">
                <h2 className="text-lg font-semibold text-foreground">
                  Заблокированные сотрудники
                </h2>
                <span className="text-xs text-muted-foreground">
                  {formatPageRange(blockedUsers.length, blockedPage)}
                </span>
              </div>
              <div className="rounded-none border border-border/50 bg-card overflow-hidden shadow-sm">
                <Table>
                  <TableHeader className="bg-muted/30">
                    <TableRow>
                      <TableHead className="w-[80px] text-center">Профиль</TableHead>
                      <TableHead>Сотрудник</TableHead>
                      <TableHead>Должность / Отдел</TableHead>
                      <TableHead className="text-right w-[120px]">Действия</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {getPagedData(blockedUsers, blockedPage).map((user) => (
                      <TableRow key={user.id} className="hover:bg-muted/20 transition-colors">
                        {renderAvatarCell(user)}
                        {renderInfoCell(user)}
                        <TableCell>
                          <div className="flex flex-col">
                            <span className="text-sm font-medium text-foreground">{user.position ?? "—"}</span>
                            <span className="text-xs text-muted-foreground">{user.department ?? "—"}</span>
                          </div>
                        </TableCell>
                        <TableCell className="text-right">
                          <UserActionsDropdown user={user} />
                        </TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
                {renderPagination(blockedUsers.length, blockedPage, setBlockedPage)}
              </div>
            </div>
          )}

        </div>
      )}
    </div>
  )
}
