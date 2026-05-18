import React from "react";
import { Link } from "react-router-dom";
import { Button } from "@/shared/ui/button";
import { Badge } from "@/shared/ui/badge";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/shared/ui/dropdown-menu";
import { EllipsisVertical, UserCheck, UserX, Settings, Trash2 } from "lucide-react";

export type UserRow = {
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

export function UserRoleBadge({ role, isBlocked }: { role: string, isBlocked: boolean }) {
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

export function UserActionsDropdown({ user }: { user: UserRow }) {
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
          <DropdownMenuItem asChild>
            <Link to={`/admin/profile/${user.id}/settings`} className="cursor-pointer gap-2 flex items-center w-full">
              <Settings className="h-4 w-4" />
              <span>Изменить</span>
            </Link>
          </DropdownMenuItem>
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
