"use client";

import React from "react";
import { Card, CardContent, CardDescription, CardHeader, CardTitle, CardFooter } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Button } from "@/components/ui/button";
import { Checkbox } from "@/components/ui/checkbox";
import { 
  Select, 
  SelectContent, 
  SelectItem, 
  SelectTrigger, 
  SelectValue 
} from "@/components/ui/select";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { DatePicker } from "@/components/ui/date-picker";
import { ChevronLeft, Save, AlertCircle } from "lucide-react";
import Link from "next/link";

// Mock data based on the schemas and requirements
const mockUser = {
  id: 1,
  login: "ivan_admin",
  email: "ivanov@company.ru",
  role: "admin",
  is_blocked: false,
  full_name: "Иванов Иван Иванович",
  phone: "+7 (999) 123-45-67",
  telegram_handle: "@ivan_dev",
  position: "Старший системный администратор",
  department: "IT Департамент",
  employee_code: "EMP-001",
  birth_date: new Date("1990-05-15"),
  hire_date: new Date("2023-01-10"),
  dismissal_date: null,
  emergency_contact: "Жена: +7 (999) 765-43-21",
  avatar_url: null,
};

export default function AdminSettingsPage() {
  return (
    <div className="max-w-4xl mx-auto space-y-8 animate-in fade-in slide-in-from-bottom-4 duration-500 pb-12">
      <div className="flex items-center gap-4">
        <Button variant="ghost" size="sm" asChild>
          <Link href="/users" className="gap-2">
            <ChevronLeft className="h-4 w-4" />
            Назад к списку
          </Link>
        </Button>
        <h1 className="text-3xl font-bold tracking-tight">Настройки профиля</h1>
      </div>

      <div className="grid gap-8">
        {/* Section 1: Identity (Read-only) */}
        <Card className="border-border/50 shadow-sm overflow-hidden">
          <CardHeader className="bg-muted/30 pb-4">
            <CardTitle className="text-xl">Идентификация</CardTitle>
            <CardDescription>Основные данные учетной записи (только чтение)</CardDescription>
          </CardHeader>
          <CardContent className="pt-6">
            <div className="flex flex-col md:flex-row gap-8 items-start">
              <div className="flex flex-col items-center gap-3">
                <Avatar className="h-24 w-24 border-2 border-border/20 shadow-md">
                  <AvatarImage src={mockUser.avatar_url || ""} />
                  <AvatarFallback className="text-xl bg-primary text-primary-foreground font-bold">
                    {mockUser.full_name.split(" ").map(n => n[0]).join("")}
                  </AvatarFallback>
                </Avatar>
                <div className="text-center">
                  <span className="text-xs font-semibold text-muted-foreground uppercase tracking-widest">Аватар</span>
                  <p className="text-[10px] text-muted-foreground mt-1">ID: #{mockUser.id}</p>
                </div>
              </div>
              
              <div className="flex-1 grid grid-cols-1 md:grid-cols-2 gap-6 w-full">
                <ReadOnlyField label="Логин" value={mockUser.login} />
                <ReadOnlyField label="Email" value={mockUser.email} />
                <ReadOnlyField label="Employee Code" value={mockUser.employee_code} />
                <div className="flex items-start gap-2 pt-2 col-span-full">
                  <AlertCircle className="h-4 w-4 text-muted-foreground mt-0.5" />
                  <p className="text-xs text-muted-foreground italic">
                    Данные поля не подлежат изменению администратором в этой форме.
                  </p>
                </div>
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Section 2: Personal Information */}
        <Card className="border-border/50 shadow-sm">
          <CardHeader className="bg-muted/30 pb-4">
            <CardTitle className="text-xl">Личная информация</CardTitle>
            <CardDescription>Контактные данные и персональные сведения сотрудника</CardDescription>
          </CardHeader>
          <CardContent className="pt-6 grid grid-cols-1 md:grid-cols-2 gap-6">
            <div className="space-y-2 col-span-full">
              <Label htmlFor="full_name">ФИО полностью</Label>
              <Input id="full_name" defaultValue={mockUser.full_name} className="h-10 transition-all focus:ring-secondary/30" />
            </div>
            <div className="space-y-2">
              <Label htmlFor="phone">Телефон</Label>
              <Input id="phone" defaultValue={mockUser.phone} placeholder="+7 (___) ___-__-__" className="h-10 transition-all focus:ring-secondary/30" />
            </div>
            <div className="space-y-2">
              <Label htmlFor="telegram">Telegram @handle</Label>
              <Input id="telegram" defaultValue={mockUser.telegram_handle} placeholder="@username" className="h-10 transition-all focus:ring-secondary/30" />
            </div>
            <div className="space-y-2">
              <Label htmlFor="birth_date">Дата рождения</Label>
              <div className="w-full">
                <DatePicker date={mockUser.birth_date} setDate={() => {}} />
              </div>
            </div>
            <div className="space-y-2 col-span-full">
              <Label htmlFor="emergency">Контакт для связи в экстренных случаях</Label>
              <Input id="emergency" defaultValue={mockUser.emergency_contact} placeholder="ФИО, телефон, степень родства" className="h-10 transition-all focus:ring-secondary/30" />
            </div>
          </CardContent>
        </Card>

        {/* Section 3: Work Information */}
        <Card className="border-border/50 shadow-sm">
          <CardHeader className="bg-muted/30 pb-4">
            <CardTitle className="text-xl">Рабочая информация</CardTitle>
            <CardDescription>Должность, роль и параметры доступа в системе</CardDescription>
          </CardHeader>
          <CardContent className="pt-6 grid grid-cols-1 md:grid-cols-2 gap-6">
            <div className="space-y-2">
              <Label htmlFor="position">Должность</Label>
              <Input id="position" defaultValue={mockUser.position} className="h-10 transition-all focus:ring-secondary/30" />
            </div>
            <div className="space-y-2">
              <Label htmlFor="department">Отдел /Подразделение</Label>
              <Input id="department" defaultValue={mockUser.department} className="h-10 transition-all focus:ring-secondary/30" />
            </div>
            <div className="space-y-2">
              <Label htmlFor="role">Роль в системе</Label>
              <Select defaultValue={mockUser.role}>
                <SelectTrigger id="role" className="h-10 transition-all focus:ring-secondary/30">
                  <SelectValue placeholder="Выберите роль" />
                </SelectTrigger>
                <SelectContent className="animate-in fade-in zoom-in-95 duration-200">
                  <SelectItem value="admin">Администратор</SelectItem>
                  <SelectItem value="employee">Сотрудник</SelectItem>
                  <SelectItem value="unverified">Не верифицирован</SelectItem>
                </SelectContent>
              </Select>
            </div>
            <div className="flex items-center gap-3 pt-8">
              <Checkbox id="is_blocked" defaultChecked={mockUser.is_blocked} className="h-5 w-5 border-2 rounded-md transition-all data-[state=checked]:bg-destructive data-[state=checked]:border-destructive" />
              <div className="grid gap-1.5 leading-none">
                <Label htmlFor="is_blocked" className="cursor-pointer font-bold text-destructive">Заблокировать доступ</Label>
                <p className="text-xs text-muted-foreground">Пользователь не сможет войти в систему</p>
              </div>
            </div>
            <div className="space-y-2">
              <Label htmlFor="hire_date">Дата приема на работу</Label>
              <div className="w-full">
                <DatePicker date={mockUser.hire_date} setDate={() => {}} />
              </div>
            </div>
            <div className="space-y-2">
              <Label htmlFor="dismissal_date" className="flex items-center gap-2">
                Дата увольнения
                <span className="text-[10px] font-normal px-1.5 py-0.5 bg-muted rounded uppercase text-muted-foreground border border-border/50">опционально</span>
              </Label>
              <div className="w-full">
                <DatePicker date={mockUser.dismissal_date || undefined} setDate={() => {}} />
              </div>
            </div>
          </CardContent>
        </Card>
      </div>

      <div className="flex justify-end gap-4">
        <Button variant="outline" className="h-11 px-8 transition-all hover:bg-muted" asChild>
          <Link href="/users">Отмена</Link>
        </Button>
        <Button className="h-11 px-10 font-bold shadow-lg shadow-primary/20 transition-all active:scale-[0.98] gap-2">
          <Save className="h-5 w-5" />
          Сохранить изменения
        </Button>
      </div>
    </div>
  );
}

function ReadOnlyField({ label, value }: { label: string, value: string | null }) {
  return (
    <div className="space-y-1.5">
      <Label className="text-xs font-semibold text-muted-foreground uppercase tracking-wider">{label}</Label>
      <div className="p-3 bg-muted/40 border border-border/30 rounded-md text-sm font-medium select-none shadow-inner opacity-70">
        {value || "—"}
      </div>
    </div>
  );
}
