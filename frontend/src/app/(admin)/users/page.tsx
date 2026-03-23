import React from "react";
import { 
  Table, 
  TableBody, 
  TableCell, 
  TableHead, 
  TableHeader, 
  TableRow 
} from "@/components/ui/table";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { EllipsisVertical, UserCheck, ShieldAlert, UserX, Settings, Trash2 } from "lucide-react";
import Link from "next/link";

type UserRow = {
  id: number;
  login: string;
  email: string | null;
  role: 'admin' | 'employee' | 'unverified';
  is_blocked: boolean;
  status: 'unverified' | 'active' | 'blocked';
  full_name: string;
  position: string | null;
  department: string | null;
  avatar_url: string | null;
  employee_code: string;
};

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

export default function UsersPage() {
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
                    {user.avatar_url && <AvatarImage src={user.avatar_url} alt={user.full_name} />}
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

function UserRoleBadge({ role, isBlocked }: { role: string, isBlocked: boolean }) {
  if (isBlocked) {
    return (
      <Badge variant="destructive" className="text-[10px] h-4 px-1.5 uppercase font-bold tracking-wider">
        Заблокирован
      </Badge>
    );
  }

  const variants: Record<string, { label: string, className: string }> = {
    admin: { label: "Админ", className: "bg-blue text-blue-foreground hover:bg-blue/90" },
    employee: { label: "Сотрудник", className: "bg-secondary text-secondary-foreground hover:bg-secondary/90" },
    unverified: { label: "Ожидает", className: "bg-accent text-accent-foreground hover:bg-accent/90" },
  };

  const config = variants[role] || variants.employee;

  return (
    <Badge className={`text-[10px] h-4 px-1.5 uppercase font-bold tracking-wider ${config.className}`}>
      {config.label}
    </Badge>
  );
}

function UserActionsDropdown({ user }: { user: UserRow }) {
  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button variant="ghost" size="icon" className="h-8 w-8 hover:bg-muted">
          <EllipsisVertical className="h-4 w-4" />
          <span className="sr-only">Открыть меню</span>
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end" className="w-[180px] animate-in fade-in-0 zoom-in-95 duration-200">
        <DropdownMenuLabel>Действия</DropdownMenuLabel>
        <DropdownMenuSeparator />
        
        {user.status === 'unverified' && (
          <DropdownMenuItem className="cursor-pointer gap-2 focus:bg-secondary focus:text-secondary-foreground">
            <UserCheck className="h-4 w-4" />
            <span>Подтвердить</span>
          </DropdownMenuItem>
        )}

        {(user.role === 'employee' || user.role === 'admin') && !user.is_blocked && (
          <>
            <DropdownMenuItem asChild>
              <Link href={`/profile/${user.id}/admin/settings`} className="cursor-pointer gap-2 flex items-center">
                <Settings className="h-4 w-4" />
                <span>Изменить</span>
              </Link>
            </DropdownMenuItem>
          </>
        )}

        {user.is_blocked && (
          <DropdownMenuItem className="cursor-pointer gap-2 focus:bg-blue focus:text-blue-foreground">
            <UserCheck className="h-4 w-4" />
            <span>Разблокировать</span>
          </DropdownMenuItem>
        )}

        {!user.is_blocked && user.status !== 'unverified' && (
          <DropdownMenuItem className="cursor-pointer gap-2 text-destructive focus:bg-destructive focus:text-destructive-foreground">
            <UserX className="h-4 w-4" />
            <span>Заблокировать</span>
          </DropdownMenuItem>
        )}

        <DropdownMenuSeparator />
        <DropdownMenuItem className="cursor-pointer gap-2 text-destructive focus:bg-destructive focus:text-destructive-foreground">
          <Trash2 className="h-4 w-4" />
          <span>Удалить</span>
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
