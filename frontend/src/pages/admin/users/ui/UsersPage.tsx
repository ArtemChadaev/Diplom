import React from "react";
import { 
  Table, 
  TableBody, 
  TableCell, 
  TableHead, 
  TableHeader, 
  TableRow 
} from "@/shared/ui/table";
import { Avatar, AvatarFallback } from "@/shared/ui/avatar";
import { UserActionsDropdown, UserRoleBadge, type UserRow } from "@/features/admin/user-actions";

const mockUsers: UserRow[] = [
  {
    id: 1,
    login: "admin_test",
    email: "admin@example.com",
    role: "admin",
    is_blocked: false,
    status: "active",
    full_name: "Иван Иванов",
    position: "Главный администратор",
    department: "IT Департамент",
    avatar_url: null,
    employee_code: "EMP-001",
  },
  {
    id: 2,
    login: "p_petrov",
    email: "petrov@example.com",
    role: "employee",
    is_blocked: false,
    status: "active",
    full_name: "Петр Петров",
    position: "Фармацевт",
    department: "Аптека №1",
    avatar_url: null,
    employee_code: "EMP-002",
  },
  {
    id: 3,
    login: "new_user",
    email: "new@example.com",
    role: "unverified",
    is_blocked: false,
    status: "unverified",
    full_name: "Алексей Сидоров",
    position: "Стажер",
    department: "Приемный покой",
    avatar_url: null,
    employee_code: "EMP-003",
  },
  {
    id: 4,
    login: "blocked_user",
    email: "blocked@example.com",
    role: "employee",
    is_blocked: true,
    status: "blocked",
    full_name: "Мария Смирнова",
    position: "Медсестра",
    department: "Терапия",
    avatar_url: null,
    employee_code: "EMP-004",
  },
];

export function UsersPage() {
  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Управление пользователями</h1>
          <p className="text-muted-foreground mt-1">
            Просмотр и управление учетными записями сотрудников
          </p>
        </div>
      </div>

      <div className="rounded-lg border border-border/50 bg-card overflow-hidden shadow-sm">
        <Table>
          <TableHeader className="bg-muted/50">
            <TableRow>
              <TableHead className="w-[80px]">Профиль</TableHead>
              <TableHead>Имя</TableHead>
              <TableHead>Должность</TableHead>
              <TableHead className="text-right">Действия</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {mockUsers.map((user) => (
              <TableRow key={user.id} className="hover:bg-muted/30 transition-colors">
                <TableCell>
                  <Avatar className="h-10 w-10 border border-border/10 shadow-sm">
                    <AvatarFallback className="bg-primary text-primary-foreground font-medium">
                      {user.full_name.split(" ").map(n => n[0]).join("")}
                    </AvatarFallback>
                  </Avatar>
                </TableCell>
                <TableCell>
                  <div className="flex flex-col gap-1">
                    <span className="font-semibold">{user.full_name}</span>
                    <div className="flex gap-2 items-center">
                      <span className="text-xs text-muted-foreground">@{user.login}</span>
                      <UserRoleBadge role={user.role} isBlocked={user.is_blocked} />
                    </div>
                  </div>
                </TableCell>
                <TableCell>
                  <div className="flex flex-col">
                    <span className="text-sm font-medium">{user.position || "—"}</span>
                    <span className="text-xs text-muted-foreground">{user.department || "—"}</span>
                  </div>
                </TableCell>
                <TableCell className="text-right">
                  <UserActionsDropdown user={user} />
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </div>
    </div>
  );
}
