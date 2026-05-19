import { UserCheck, UserX, Settings } from "lucide-react"
import { Link } from "react-router-dom"

import { useAdminUsersStore, type UserProfile } from "@/entities/user"

import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from "@/shared/ui/alert-dialog"
import { Badge } from "@/shared/ui/badge"
import { Button } from "@/shared/ui/button"


// Экспортируем тип UserRow как синоним для UserProfile для сохранения совместимости
export type UserRow = UserProfile;

export function UserRoleBadge({ role }: { role: string; isBlocked?: boolean }) {
  const variants: Record<string, { label: string; className: string }> = {
    admin: { label: "Администратор", className: "bg-primary text-primary-foreground" },
    qp: { label: "Уполн. лицо", className: "bg-info text-info-foreground" },
    warehouse_manager: { label: "Зав. складом", className: "bg-secondary text-secondary-foreground" },
    storekeeper: { label: "Кладовщик", className: "bg-accent text-accent-foreground" },
    pharmacist: { label: "Фармацевт", className: "bg-muted text-muted-foreground border border-border" },
  }

  const config = variants[role] ?? { label: role, className: "bg-secondary text-secondary-foreground" }

  return (
    <Badge className={`text-[10px] h-4 px-1.5 uppercase font-bold tracking-wider ${config.className}`}>
      {config.label}
    </Badge>
  )
}

export function UserActionsDropdown({ user }: { user: UserProfile }) {
  const toggleBlockUser = useAdminUsersStore((state) => state.toggleBlockUser)

  const handleUnblock = async () => {
    try {
      await toggleBlockUser(user.id, false)
    } catch (err) {
      console.error("Failed to unblock user:", err)
    }
  }

  const handleBlock = async () => {
    try {
      await toggleBlockUser(user.id, true)
    } catch (err) {
      console.error("Failed to block user:", err)
    }
  }

  return (
    <div className="flex items-center justify-end gap-2">
      {/* Кнопка изменения настроек (шестеренка) - показываем только для незаблокированных */}
      {!user.is_blocked ? (
        <Button variant="ghost" size="icon" className="h-8 w-8 hover:bg-muted" asChild>
          <Link to={`/admin/profile/${String(user.id)}/settings`} title="Изменить настройки">
            <Settings className="h-4 w-4 text-muted-foreground" />
            <span className="sr-only">Изменить</span>
          </Link>
        </Button>
      ) : (
        <div className="w-8 h-8" />
      )}

      {/* Кнопка блокировки / разблокировки */}
      {user.is_blocked ? (
        <Button
          variant="ghost"
          size="icon"
          className="h-8 w-8 hover:bg-muted text-emerald-600 dark:text-emerald-400"
          onClick={handleUnblock}
          title="Разблокировать"
        >
          <UserCheck className="h-4 w-4" />
          <span className="sr-only">Разблокировать</span>
        </Button>
      ) : (
        <AlertDialog>
          <AlertDialogTrigger asChild>
            <Button
              variant="ghost"
              size="icon"
              className="h-8 w-8 hover:bg-muted text-destructive"
              title="Заблокировать"
            >
              <UserX className="h-4 w-4" />
              <span className="sr-only">Заблокировать</span>
            </Button>
          </AlertDialogTrigger>
          <AlertDialogContent>
            <AlertDialogHeader>
              <AlertDialogTitle>Блокировка пользователя</AlertDialogTitle>
              <AlertDialogDescription className="space-y-3 text-muted-foreground">
                <span>Вы уверены, что хотите заблокировать данного пользователя? Он потеряет доступ к системе.</span>
                <div className="rounded-none border border-border/50 bg-muted/30 p-3 text-left space-y-1.5 text-foreground mt-2">
                  {(!user.full_name || user.full_name === "") && (
                    <div className="text-destructive font-semibold text-xs pb-0.5">
                      Профиль не заполнен
                    </div>
                  )}
                  <div>
                    <span className="text-muted-foreground text-xs">Имя: </span>
                    <span className="font-semibold text-xs">{user.full_name && user.full_name !== "" ? user.full_name : "—"}</span>
                  </div>
                  <div>
                    <span className="text-muted-foreground text-xs">Должность: </span>
                    <span className="font-semibold text-xs">{user.position && user.position !== "" ? user.position : "—"}</span>
                  </div>
                  <div>
                    <span className="text-muted-foreground text-xs">Email: </span>
                    <span className="font-semibold text-xs">{user.email}</span>
                  </div>
                </div>
              </AlertDialogDescription>
            </AlertDialogHeader>
            <AlertDialogFooter>
              <AlertDialogCancel>Отмена</AlertDialogCancel>
              <AlertDialogAction
                onClick={handleBlock}
                className="bg-destructive hover:bg-destructive/90 text-destructive-foreground"
              >
                Заблокировать
              </AlertDialogAction>
            </AlertDialogFooter>
          </AlertDialogContent>
        </AlertDialog>
      )}
    </div>
  )
}
