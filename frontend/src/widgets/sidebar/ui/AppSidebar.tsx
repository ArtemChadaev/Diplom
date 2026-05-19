import {
  AlertTriangle,
  CheckSquare,
  ClipboardList,
  Database,
  Download,
  Home,
  Layers,
  LogOut,
  Package,
  Search,
  Settings,
  ShieldAlert,
  Thermometer,
  Truck,
  Users,
} from "lucide-react"
import { Link, useLocation } from "react-router-dom"

import { useAuthStore } from "@/entities/user"

import { Avatar, AvatarFallback } from "@/shared/ui/avatar"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/shared/ui/dropdown-menu"
import { Separator } from "@/shared/ui/separator"
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from "@/shared/ui/sidebar"

interface MenuItem {
  title: string
  path: string
  icon: React.ComponentType<{ className?: string }>
  roles?: string[]
}

const roleTranslations: Record<string, string> = {
  admin: "Администратор",
  qp: "Уполномоченное лицо",
  warehouse_manager: "Зав. складом",
  storekeeper: "Кладовщик",
  pharmacist: "Провизор",
  unverified: "Не верифицирован",
}

export function AppSidebar() {
  const { user, logout } = useAuthStore()
  const location = useLocation()

  const generalItems: MenuItem[] = [
    { title: "Панель управления", path: "/", icon: Home },
    { title: "Быстрый поиск", path: "/search", icon: Search },
    { title: "Каталог товаров", path: "/products", icon: Package },
    { title: "Поставщики", path: "/suppliers", icon: Truck },
    {
      title: "Приемка товара",
      path: "/receiving",
      icon: Download,
      roles: ["admin", "qp", "warehouse_manager", "storekeeper"],
    },
    { title: "Остатки / Партии", path: "/batches", icon: Database },
    { title: "Изъятия / Брак", path: "/warehouse/recalled", icon: ShieldAlert },
    { title: "Микроклимат", path: "/microclimate", icon: Thermometer },
    { title: "Заказы (Отгрузка)", path: "/orders", icon: ClipboardList },
    { title: "Претензии (Брак)", path: "/claims", icon: AlertTriangle },
  ]

  const controlItems: MenuItem[] = [
    {
      title: "Зоны склада",
      path: "/warehouse/zones",
      icon: Layers,
      roles: ["admin", "warehouse_manager", "qp"],
    },
    {
      title: "Инвентаризация",
      path: "/inventory",
      icon: CheckSquare,
      roles: ["admin", "warehouse_manager", "storekeeper"],
    },
  ]

  const adminItems: MenuItem[] = [
    {
      title: "Пользователи",
      path: "/admin/users",
      icon: Users,
      roles: ["admin"],
    },
  ]

  const filterByRole = (items: MenuItem[]) => {
    return items.filter((item) => {
      if (!item.roles) return true
      if (!user?.role) return false
      return item.roles.includes(user.role)
    })
  }

  const visibleGeneral = filterByRole(generalItems)
  const visibleControl = filterByRole(controlItems)
  const visibleAdmin = filterByRole(adminItems)

  const renderItem = (item: MenuItem) => {
    const Icon = item.icon
    const isActive = location.pathname === item.path
    return (
      <SidebarMenuItem key={item.path}>
        <SidebarMenuButton
          asChild
          isActive={isActive}
          tooltip={item.title}
          className="transition-colors hover:bg-sidebar-accent/50"
        >
          <Link to={item.path} className="flex items-center w-full gap-3 px-3 py-2 text-sm">
            <Icon className="h-4 w-4 shrink-0" />
            <span className="group-data-[collapsible=icon]:hidden">{item.title}</span>
          </Link>
        </SidebarMenuButton>
      </SidebarMenuItem>
    )
  }

  // Format avatar initials
  const fallback = (
    user?.email ? user.email.substring(0, 2) : "ph"
  ).toLowerCase()

  // Format display name
  const displayName = user?.email ? user.email.split("@")[0] : "Пользователь"

  return (
    <Sidebar collapsible="icon" className="border-r border-sidebar-border bg-sidebar/50 backdrop-blur-md">
      {/* Sidebar Header */}
      <SidebarHeader className="p-3">
        <SidebarMenu>
          <SidebarMenuItem>
            <SidebarMenuButton size="lg" className="w-full justify-start hover:bg-transparent active:bg-transparent cursor-default">
              <div className="flex h-7 w-7 shrink-0 items-center justify-center rounded-none bg-emerald-600 text-white font-bold text-[13px] shadow-none">
                ph
              </div>
              <div className="flex flex-col group-data-[collapsible=icon]:hidden">
                <span className="font-sans font-bold text-sm leading-tight text-sidebar-foreground">
                  pharma-hub
                </span>
                <span className="text-[10px] text-muted-foreground leading-none">
                  ERP System
                </span>
              </div>
            </SidebarMenuButton>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarHeader>
      <Separator className="bg-sidebar-border" />

      {/* Sidebar Content */}
      <SidebarContent>
        {visibleGeneral.length > 0 && (
          <SidebarGroup>
            <SidebarGroupLabel className="group-data-[collapsible=icon]:hidden">
              Операции
            </SidebarGroupLabel>
            <SidebarGroupContent>
              <SidebarMenu>
                {visibleGeneral.map(renderItem)}
              </SidebarMenu>
            </SidebarGroupContent>
          </SidebarGroup>
        )}

        {visibleControl.length > 0 && (
          <SidebarGroup>
            <SidebarGroupLabel className="group-data-[collapsible=icon]:hidden">
              Контроль и зоны
            </SidebarGroupLabel>
            <SidebarGroupContent>
              <SidebarMenu>
                {visibleControl.map(renderItem)}
              </SidebarMenu>
            </SidebarGroupContent>
          </SidebarGroup>
        )}

        {visibleAdmin.length > 0 && (
          <SidebarGroup>
            <SidebarGroupLabel className="group-data-[collapsible=icon]:hidden">
              Администрирование
            </SidebarGroupLabel>
            <SidebarGroupContent>
              <SidebarMenu>
                {visibleAdmin.map(renderItem)}
              </SidebarMenu>
            </SidebarGroupContent>
          </SidebarGroup>
        )}
      </SidebarContent>

      <Separator className="bg-sidebar-border" />

      {/* Sidebar Footer */}
      <SidebarFooter className="p-3 gap-3">
        {user?.role === "admin" && (
          <SidebarMenu>
            <SidebarMenuItem>
              <SidebarMenuButton
                asChild
                isActive={location.pathname === "/settings"}
                tooltip="Настройки системы"
                className="transition-colors hover:bg-sidebar-accent/50"
              >
                <Link to="/settings" className="flex items-center w-full gap-3 px-3 py-2 text-sm">
                  <Settings className="h-4 w-4 shrink-0" />
                  <span className="group-data-[collapsible=icon]:hidden">Настройки системы</span>
                </Link>
              </SidebarMenuButton>
            </SidebarMenuItem>
          </SidebarMenu>
        )}

        <SidebarMenu>
          <SidebarMenuItem>
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <SidebarMenuButton
                  size="lg"
                  className="w-full flex items-center gap-3 px-2 py-1.5 rounded-lg hover:bg-sidebar-accent transition-colors text-left"
                >
                  <Avatar className="h-8 w-8 rounded-full border border-border/40 shrink-0">
                    <AvatarFallback className="bg-emerald-600/10 text-emerald-600 font-semibold text-xs rounded-full flex items-center justify-center uppercase">
                      {fallback}
                    </AvatarFallback>
                  </Avatar>
                  <div className="flex-1 min-w-0 flex flex-col justify-center group-data-[collapsible=icon]:hidden">
                    <span className="text-xs font-semibold text-sidebar-foreground truncate leading-tight">
                      {displayName}
                    </span>
                    <span className="text-[10px] text-muted-foreground truncate leading-none mt-0.5">
                      {roleTranslations[user?.role ?? ""] ?? "Сотрудник"}
                    </span>
                  </div>
                </SidebarMenuButton>
              </DropdownMenuTrigger>
              <DropdownMenuContent
                side="right"
                align="end"
                sideOffset={8}
                className="w-56 p-1 border border-border/40 bg-popover/95 backdrop-blur-md shadow-xl"
              >
                {user && (
                  <DropdownMenuItem asChild>
                    <Link
                      to={`/admin/profile/${user.id}/settings`}
                      className="w-full flex items-center gap-2 px-2.5 py-2 text-xs cursor-pointer rounded hover:bg-accent hover:text-accent-foreground"
                    >
                      <Settings className="h-4 w-4 text-muted-foreground" />
                      <span>Настройки профиля</span>
                    </Link>
                  </DropdownMenuItem>
                )}
                <DropdownMenuSeparator className="bg-border/40" />
                <DropdownMenuItem
                  onClick={logout}
                  className="w-full flex items-center gap-2 px-2.5 py-2 text-xs text-destructive cursor-pointer rounded hover:bg-destructive/10 hover:text-destructive focus:bg-destructive/10 focus:text-destructive"
                >
                  <LogOut className="h-4 w-4" />
                  <span>Выйти из аккаунта</span>
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarFooter>
    </Sidebar>
  )
}
