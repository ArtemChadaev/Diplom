import { useEffect, useState } from "react"
import { useParams, useNavigate } from "react-router-dom"
import { Save, Camera, ImageOff } from "lucide-react"

import {
  InputOTP,
  InputOTPGroup,
  InputOTPSeparator,
  InputOTPSlot,
} from "@/shared/ui/input-otp"

import { useUserStore } from "@/entities/user"
import { api } from "@/shared/api"

import { Avatar, AvatarFallback, AvatarImage } from "@/shared/ui/avatar"
import { Button } from "@/shared/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/shared/ui/card"
import { Checkbox } from "@/shared/ui/checkbox"
import { DatePicker } from "@/shared/ui/date-picker"
import { Input } from "@/shared/ui/input"
import { Label } from "@/shared/ui/label"
import { Switch } from "@/shared/ui/switch"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/shared/ui/select"
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/shared/ui/alert-dialog"

import { ReadOnlyField } from "./ReadOnlyField"

export function ProfileSettingsPage() {
  const { id: rawId } = useParams<{ id: string }>()
  const id = rawId === "undefined" || rawId === "null" ? undefined : rawId
  const navigate = useNavigate()
  const currentUser = useUserStore((state) => state.user)
  
  const targetId = id || (currentUser ? String(currentUser.id) : "")
  const isAdmin = currentUser?.role === "admin"
  const isSelf = !id || String(currentUser?.id) === String(targetId)

  const [isLoading, setIsLoading] = useState(true)
  const [isSaving, setIsSaving] = useState(false)
  const [error, setError] = useState<string | null>(null)

  // Form states
  const [email, setEmail] = useState("")
  const [employeeCode, setEmployeeCode] = useState("")
  
  const [fullName, setFullName] = useState("")
  const [phone, setPhone] = useState("")
  const [corporateEmail, setCorporateEmail] = useState("")
  const [medicalBookScanUrl, setMedicalBookScanUrl] = useState("")
  const [gdpTrainingHistory, setGdpTrainingHistory] = useState("")
  const [birthDate, setBirthDate] = useState<Date | undefined>(undefined)
  const [avatarUrl, setAvatarUrl] = useState<string | null>(null)

  const [position, setPosition] = useState("")
  const [department, setDepartment] = useState("")
  const [role, setRole] = useState("")
  const [isBlocked, setIsBlocked] = useState(false)
  const [specialZoneAccess, setSpecialZoneAccess] = useState(false)
  const [hireDate, setHireDate] = useState<Date | undefined>(undefined)
  const [dismissalDate, setDismissalDate] = useState<Date | undefined>(undefined)

  // Avatar Dialog state
  const [isAvatarDialogOpen, setIsAvatarDialogOpen] = useState(false)
  const [newAvatarUrlInput, setNewAvatarUrlInput] = useState("")

  const [employeeCodeError, setEmployeeCodeError] = useState<string | null>(null)

  const handleEmployeeCodeChange = (val: string) => {
    setEmployeeCodeError(null)

    // Check first 2 characters (must be English letters only)
    const firstTwo = val.slice(0, Math.min(val.length, 2))
    const hasNonEnglish = /[^a-zA-Z]/.test(firstTwo)
    if (hasNonEnglish) {
      setEmployeeCodeError("Используйте только латинские буквы")
    }

    // Check next 3 characters (must be digits only)
    const lastThree = val.slice(2)
    const hasNonDigit = /[^0-9]/.test(lastThree)
    if (hasNonDigit) {
      setEmployeeCodeError("Вторая часть должна состоять из цифр")
    }

    // Clean and normalize value:
    let letters = firstTwo.toUpperCase().replace(/[^A-Z]/g, "")
    let digits = lastThree.replace(/[^0-9]/g, "")

    if (letters.length > 0) {
      if (val.length > 2) {
        setEmployeeCode(letters + "-" + digits)
      } else {
        setEmployeeCode(letters)
      }
    } else {
      setEmployeeCode("")
    }
  }

  const handlePhoneChange = (val: string) => {
    const digitsOnly = val.replace(/[^0-9]/g, "").slice(0, 10)
    setPhone("+7" + digitsOnly)
  }

  useEffect(() => {
    if (!currentUser) return

    // Access control: if not admin and not self, eject immediately
    if (!isAdmin && !isSelf) {
      navigate(-1)
      return
    }

    const loadProfile = async () => {
      setIsLoading(true)
      setError(null)
      try {
        // Fetch base user profile (accessible by anyone for self, or admin for anyone)
        const url = isAdmin ? `/admin/users/${id}` : `/users/me`
        const userRes = await api.get<any>(url)

        let empRes: any = null
        if (isAdmin && targetId) {
          try {
            empRes = await api.get<any>(`/admin/employees/${targetId}`)
          } catch (e) {
            console.error("No employee profile found in DB, using fallback:", e)
          }
        }

        // Merge backend profile data
        const merged = { ...userRes, ...empRes }

        // Load emergency contact and telegram handle from localStorage backup
        const backupData = localStorage.getItem(`profile_backup_${targetId}`)
        if (backupData) {
          const parsed = JSON.parse(backupData)
          if (parsed.employee_code) merged.employee_code = parsed.employee_code
          if (parsed.phone && !merged.phone) merged.phone = parsed.phone
          if (parsed.birth_date && !merged.birth_date) merged.birth_date = parsed.birth_date
          if (parsed.hire_date && !merged.hire_date) merged.hire_date = parsed.hire_date
          if (parsed.dismissal_date && !merged.dismissal_date) merged.dismissal_date = parsed.dismissal_date
          
          if (parsed.corporate_email && !merged.corporate_email) merged.corporate_email = parsed.corporate_email
          if (parsed.medical_book_scan_url && !merged.medical_book_scan_url) merged.medical_book_scan_url = parsed.medical_book_scan_url
          if (parsed.gdp_training_history && !merged.gdp_training_history) merged.gdp_training_history = parsed.gdp_training_history
          if (parsed.special_zone_access !== undefined) merged.special_zone_access = parsed.special_zone_access
        }

        setEmail(merged.email ?? "")
        setEmployeeCode(merged.employee_code ?? "")
        
        setFullName(merged.full_name ?? "")
        setPhone(merged.phone ?? "")
        setCorporateEmail(merged.corporate_email ?? merged.corporateEmail ?? "")
        setMedicalBookScanUrl(merged.medical_book_scan_url ?? merged.medicalBookScanUrl ?? "")
        setGdpTrainingHistory(typeof (merged.gdp_training_history ?? merged.gdpTrainingHistory) === 'string'
          ? (merged.gdp_training_history ?? merged.gdpTrainingHistory)
          : JSON.stringify(merged.gdp_training_history ?? merged.gdpTrainingHistory ?? ""))
        setBirthDate(merged.birth_date ? new Date(merged.birth_date) : undefined)
        setAvatarUrl(merged.avatar_url ?? null)

        setPosition(merged.position ?? "")
        setDepartment(merged.department ?? "")
        setRole(merged.role ?? "")
        setIsBlocked(merged.is_blocked ?? false)
        setSpecialZoneAccess(merged.special_zone_access ?? merged.specialZoneAccess ?? false)
        setHireDate(merged.hire_date ? new Date(merged.hire_date) : undefined)
        setDismissalDate(merged.dismissal_date ? new Date(merged.dismissal_date) : undefined)
      } catch (err) {
        console.error("Failed to load user profile:", err)
        setError("Не удалось загрузить данные профиля с сервера")
      } finally {
        setIsLoading(false)
      }
    }

    void loadProfile()
  }, [targetId, currentUser, isAdmin, isSelf, navigate])

  const handleSave = async () => {
    setIsSaving(true)
    try {
      // Save local backup for local storage fields
      const backup = {
        employee_code: employeeCode,
        phone,
        birth_date: birthDate ? birthDate.toISOString() : null,
        corporate_email: corporateEmail,
        medical_book_scan_url: medicalBookScanUrl,
        gdp_training_history: gdpTrainingHistory,
        special_zone_access: specialZoneAccess,
        hire_date: hireDate ? hireDate.toISOString() : null,
        dismissal_date: dismissalDate ? dismissalDate.toISOString() : null,
      }
      localStorage.setItem(`profile_backup_${targetId}`, JSON.stringify(backup))

      if (isAdmin) {
        // 1. Update employee profile
        const empInput = {
          employee_code: employeeCode,
          full_name: fullName,
          phone: phone,
          position: position,
          department: department,
          birth_date: birthDate ? birthDate.toISOString() : null,
          avatar_url: avatarUrl,
          hire_date: hireDate ? hireDate.toISOString() : null,
          dismissal_date: dismissalDate ? dismissalDate.toISOString() : null,
          corporate_email: corporateEmail,
          medical_book_scan_url: medicalBookScanUrl,
          gdp_training_history: gdpTrainingHistory,
          special_zone_access: specialZoneAccess,
        }
        await api.patch(`/admin/employees/${targetId}`, empInput)

        // 2. Update role
        await api.patch(`/admin/users/${targetId}/role`, {
          role: role,
        })
        
        // 3. Update blocked status
        await api.patch(`/admin/users/${targetId}/blocked`, {
          blocked: isBlocked,
        })
      } else {
        // Regular user updates their own profile
        if (isSelf) {
          await api.patch("/users/me", {
            full_name: fullName,
            phone: phone,
            birth_date: birthDate ? birthDate.toISOString() : null,
            avatar_url: avatarUrl,
            corporate_email: corporateEmail,
            medical_book_scan_url: medicalBookScanUrl,
            gdp_training_history: gdpTrainingHistory,
          })

          useUserStore.getState().updateUser({
            full_name: fullName,
            avatar_url: avatarUrl,
          })
        }
      }

      alert("Изменения успешно сохранены!")
    } catch (err) {
      console.error("Failed to save changes:", err)
      alert("Не удалось сохранить изменения")
    } finally {
      setIsSaving(false)
    }
  }

  // Handle Avatar click
  const handleAvatarClick = () => {
    setNewAvatarUrlInput(avatarUrl ?? "")
    setIsAvatarDialogOpen(true)
  }

  const handleSaveAvatar = () => {
    setAvatarUrl(newAvatarUrlInput !== "" ? newAvatarUrlInput : null)
    setIsAvatarDialogOpen(false)
  }

  const handleDeleteAvatar = () => {
    setAvatarUrl(null)
    setIsAvatarDialogOpen(false)
  }

  if (isLoading) {
    return (
      <div className="max-w-4xl mx-auto space-y-8 animate-pulse pb-12 pt-6">
        <div className="h-8 w-64 bg-muted rounded" />
        <div className="grid gap-8">
          <div className="border border-border/50 bg-card rounded-none p-6 space-y-6">
            <div className="h-6 w-32 bg-muted rounded" />
            <div className="flex gap-8 items-start">
              <div className="h-24 w-24 bg-muted rounded-full" />
              <div className="flex-1 grid grid-cols-2 gap-6">
                <div className="h-10 bg-muted rounded" />
                <div className="h-10 bg-muted rounded" />
                <div className="h-10 bg-muted rounded" />
              </div>
            </div>
          </div>
          <div className="border border-border/50 bg-card rounded-none p-6 space-y-6">
            <div className="h-6 w-48 bg-muted rounded" />
            <div className="grid grid-cols-2 gap-6">
              <div className="h-10 bg-muted rounded col-span-2" />
              <div className="h-10 bg-muted rounded" />
              <div className="h-10 bg-muted rounded" />
              <div className="h-10 bg-muted rounded" />
              <div className="h-10 bg-muted rounded" />
            </div>
          </div>
        </div>
      </div>
    )
  }

  if (error) {
    return (
      <div className="p-6 text-center space-y-4">
        <div className="text-destructive font-semibold text-lg">{error}</div>
        <Button onClick={() => window.location.reload()}>Попробовать снова</Button>
      </div>
    )
  }

  const initials = fullName
    ? fullName
        .split(" ")
        .filter(Boolean)
        .map((n) => n[0])
        .join("")
    : "?"

  const canEdit = isAdmin || isSelf

  return (
    <div className="max-w-4xl mx-auto space-y-8 animate-in fade-in duration-300 pb-12 pt-6">
      <div className="flex items-center justify-between">
        <h1 className="text-3xl font-bold tracking-tight">Настройки профиля</h1>
      </div>

      <div className="grid gap-8">
        {/* CARD 1: ИДЕНТИФИКАЦИЯ */}
        <Card className="border-border/50 shadow-sm overflow-hidden rounded-none pt-0">
          <CardHeader className="bg-muted/30 pt-4 pb-4">
            <CardTitle className="text-xl">Идентификация</CardTitle>
            <CardDescription>Основные данные учетной записи (только чтение)</CardDescription>
          </CardHeader>
          <CardContent className="pt-6">
            <div className="flex flex-col md:flex-row gap-8 items-start">
              {/* Аватар с оверлеем изменения */}
              <div className="flex flex-col items-center gap-3">
                <div className="relative group cursor-pointer" onClick={handleAvatarClick} title="Нажмите, чтобы изменить аватар">
                  <Avatar className="h-24 w-24 border-2 border-border/20 shadow-md transition-transform duration-300 group-hover:scale-105" size="lg">
                    {avatarUrl && <AvatarImage src={avatarUrl} alt={fullName} />}
                    <AvatarFallback className="text-xl bg-primary text-primary-foreground font-bold flex items-center justify-center">
                      {initials}
                    </AvatarFallback>
                  </Avatar>
                  <div className="absolute inset-0 bg-black/40 rounded-full flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity duration-300">
                    <Camera className="h-6 w-6 text-white" />
                  </div>
                </div>
                <div className="text-center">
                  <span className="text-[10px] font-semibold text-muted-foreground uppercase tracking-widest">Аватар</span>
                  <p className="text-[10px] text-muted-foreground mt-1">ID: #{targetId}</p>
                </div>
              </div>
              
              <div className="flex-1 grid grid-cols-1 md:grid-cols-2 gap-6 w-full">
                <ReadOnlyField label="Логин" value={email} />
                <div className="space-y-2" title={!canEdit ? "Нет прав для изменения" : undefined}>
                  <Label className={`text-xs font-semibold text-muted-foreground uppercase tracking-wider ${!canEdit ? "opacity-60" : ""}`}>
                    Код сотрудника
                  </Label>
                  <div className="flex flex-col gap-1.5">
                    <InputOTP
                      maxLength={5}
                      value={employeeCode.replace("-", "")}
                      onChange={handleEmployeeCodeChange}
                      disabled={!canEdit}
                    >
                      <InputOTPGroup>
                        <InputOTPSlot index={0} className="h-10 w-10 text-sm" />
                        <InputOTPSlot index={1} className="h-10 w-10 text-sm" />
                      </InputOTPGroup>
                      <InputOTPSeparator />
                      <InputOTPGroup>
                        <InputOTPSlot index={2} className="h-10 w-10 text-sm" />
                        <InputOTPSlot index={3} className="h-10 w-10 text-sm" />
                        <InputOTPSlot index={4} className="h-10 w-10 text-sm" />
                      </InputOTPGroup>
                    </InputOTP>
                    {employeeCodeError && (
                      <span className="text-[11px] text-destructive font-medium animate-in fade-in duration-200">
                        {employeeCodeError}
                      </span>
                    )}
                  </div>
                </div>
              </div>
            </div>
          </CardContent>
        </Card>

        {/* CARD 2: ЛИЧНАЯ ИНФОРМАЦИЯ */}
        <Card className="border-border/50 shadow-sm overflow-hidden rounded-none pt-0">
          <CardHeader className="bg-muted/30 pt-4 pb-4">
            <CardTitle className="text-xl">Личная информация</CardTitle>
            <CardDescription>Контактные данные и персональные сведения сотрудника</CardDescription>
          </CardHeader>
          <CardContent className="pt-6 grid grid-cols-1 md:grid-cols-2 gap-6">
            <div className="space-y-2 col-span-full" title={!canEdit ? "Нет прав для изменения" : undefined}>
              <Label htmlFor="full_name" className={!canEdit ? "opacity-60" : ""}>ФИО полностью</Label>
              <Input
                id="full_name"
                value={fullName}
                onChange={(e) => setFullName(e.target.value)}
                disabled={!canEdit}
                className="h-10 transition-all focus:ring-secondary/30 disabled:opacity-50 disabled:bg-muted disabled:cursor-not-allowed"
              />
            </div>
            <div className="space-y-2 col-span-full" title={!canEdit ? "Нет прав для изменения" : undefined}>
              <Label htmlFor="phone" className={!canEdit ? "opacity-60" : ""}>Телефон</Label>
              <div className="flex items-center gap-2">
                <div className="h-10 px-3 bg-muted border border-border text-muted-foreground flex items-center justify-center font-semibold rounded-md select-none text-sm">
                  +7
                </div>
                <InputOTP
                  maxLength={10}
                  value={phone.replace(/[^0-9]/g, "").slice(1, 11)}
                  onChange={handlePhoneChange}
                  disabled={!canEdit}
                >
                  <InputOTPGroup>
                    <InputOTPSlot index={0} className="h-10 w-9 text-sm" />
                    <InputOTPSlot index={1} className="h-10 w-9 text-sm" />
                    <InputOTPSlot index={2} className="h-10 w-9 text-sm" />
                  </InputOTPGroup>
                  <InputOTPSeparator />
                  <InputOTPGroup>
                    <InputOTPSlot index={3} className="h-10 w-9 text-sm" />
                    <InputOTPSlot index={4} className="h-10 w-9 text-sm" />
                    <InputOTPSlot index={5} className="h-10 w-9 text-sm" />
                  </InputOTPGroup>
                  <InputOTPSeparator />
                  <InputOTPGroup>
                    <InputOTPSlot index={6} className="h-10 w-9 text-sm" />
                    <InputOTPSlot index={7} className="h-10 w-9 text-sm" />
                    <InputOTPSlot index={8} className="h-10 w-9 text-sm" />
                    <InputOTPSlot index={9} className="h-10 w-9 text-sm" />
                  </InputOTPGroup>
                </InputOTP>
              </div>
            </div>
            <div className="space-y-2" title={!canEdit ? "Нет прав для изменения" : undefined}>
              <Label htmlFor="corporate_email" className={!canEdit ? "opacity-60" : ""}>Корпоративный Email</Label>
              <Input
                id="corporate_email"
                value={corporateEmail}
                onChange={(e) => setCorporateEmail(e.target.value)}
                disabled={!canEdit}
                placeholder="corporate@company.com"
                className="h-10 transition-all focus:ring-secondary/30 disabled:opacity-50 disabled:bg-muted disabled:cursor-not-allowed"
              />
            </div>
            <div className="space-y-2" title={!canEdit ? "Нет прав для изменения" : undefined}>
              <Label className={!canEdit ? "opacity-60" : ""}>Дата рождения</Label>
              <DatePicker date={birthDate} setDate={setBirthDate} disabled={!canEdit} />
            </div>
            <div className="space-y-2" title={!canEdit ? "Нет прав для изменения" : undefined}>
              <Label htmlFor="medical_book" className={!canEdit ? "opacity-60" : ""}>Ссылка на скан медкнижки</Label>
              <Input
                id="medical_book"
                value={medicalBookScanUrl}
                onChange={(e) => setMedicalBookScanUrl(e.target.value)}
                disabled={!canEdit}
                placeholder="https://..."
                className="h-10 transition-all focus:ring-secondary/30 disabled:opacity-50 disabled:bg-muted disabled:cursor-not-allowed"
              />
            </div>
            <div className="space-y-2" title={!canEdit ? "Нет прав для изменения" : undefined}>
              <Label htmlFor="gdp_history" className={!canEdit ? "opacity-60" : ""}>История обучения GDP</Label>
              <Input
                id="gdp_history"
                value={gdpTrainingHistory}
                onChange={(e) => setGdpTrainingHistory(e.target.value)}
                disabled={!canEdit}
                placeholder="История обучения..."
                className="h-10 transition-all focus:ring-secondary/30 disabled:opacity-50 disabled:bg-muted disabled:cursor-not-allowed"
              />
            </div>
          </CardContent>
        </Card>

        {/* CARD 3: РАБОЧАЯ ИНФОРМАЦИЯ */}
        <Card className="border-border/50 shadow-sm overflow-hidden rounded-none pt-0">
          <CardHeader className="bg-muted/30 pt-4 pb-4">
            <CardTitle className="text-xl">Рабочая информация</CardTitle>
            <CardDescription>Должность, роль и параметры доступа в системе</CardDescription>
          </CardHeader>
          <CardContent className="pt-6 grid grid-cols-1 md:grid-cols-2 gap-6">
            {/* Ряд 1: Должность и Отдел */}
            <div className="space-y-2" title={!isAdmin ? "Нет прав для изменения" : undefined}>
              <Label htmlFor="position" className={!isAdmin ? "opacity-60" : ""}>Должность</Label>
              <Input
                id="position"
                value={position}
                onChange={(e) => setPosition(e.target.value)}
                disabled={!isAdmin}
                className="h-10 transition-all focus:ring-secondary/30 disabled:opacity-50 disabled:bg-muted disabled:cursor-not-allowed"
              />
            </div>
            <div className="space-y-2" title={!isAdmin ? "Нет прав для изменения" : undefined}>
              <Label htmlFor="department" className={!isAdmin ? "opacity-60" : ""}>Отдел / Подразделение</Label>
              <Input
                id="department"
                value={department}
                onChange={(e) => setDepartment(e.target.value)}
                disabled={!isAdmin}
                className="h-10 transition-all focus:ring-secondary/30 disabled:opacity-50 disabled:bg-muted disabled:cursor-not-allowed"
              />
            </div>

            {/* Ряд 2: Даты приема и увольнения */}
            <div className="space-y-2" title={!isAdmin ? "Нет прав для изменения" : undefined}>
              <Label className={!isAdmin ? "opacity-60" : ""}>Дата приема на работу</Label>
              <DatePicker date={hireDate} setDate={setHireDate} disabled={!isAdmin} />
            </div>
            <div className="space-y-2" title={!isAdmin ? "Нет прав для изменения" : undefined}>
              <Label className={!isAdmin ? "opacity-60" : ""}>Дата увольнения</Label>
              <DatePicker date={dismissalDate} setDate={setDismissalDate} disabled={!isAdmin} />
            </div>

            {/* Ряд 3: Роль в системе и Доступ в специальные зоны */}
            <div className="space-y-2" title={!isAdmin ? "Нет прав для изменения" : undefined}>
              <Label htmlFor="role" className={!isAdmin ? "opacity-60" : ""}>Роль в системе</Label>
              <Select value={role} onValueChange={setRole} disabled={!isAdmin}>
                <SelectTrigger id="role" className="w-[240px] h-10 px-3 text-sm font-normal transition-all focus:ring-secondary/30 disabled:opacity-50 disabled:bg-muted disabled:cursor-not-allowed">
                  <SelectValue placeholder="Выберите роль" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="admin">Администратор</SelectItem>
                  <SelectItem value="qp">Уполн. лицо</SelectItem>
                  <SelectItem value="warehouse_manager">Зав. складом</SelectItem>
                  <SelectItem value="storekeeper">Кладовщик</SelectItem>
                  <SelectItem value="pharmacist">Фармацевт</SelectItem>
                </SelectContent>
              </Select>
            </div>

            <div className="space-y-2" title={!isAdmin ? "Нет прав для изменения" : undefined}>
              <Label htmlFor="special_zone_access" className={!isAdmin ? "opacity-60" : ""}>
                Special Zone Access
              </Label>
              <div className="h-10 flex items-center">
                <Switch
                  id="special_zone_access"
                  checked={specialZoneAccess}
                  onCheckedChange={setSpecialZoneAccess}
                  disabled={!isAdmin}
                />
              </div>
            </div>
          </CardContent>
        </Card>
      </div>

      {/* КНОПКА СОХРАНИТЬ И БЛОКИРОВКА */}
      <div className="flex items-center justify-between pt-4 border-t border-border/30">
        <div>
          {isAdmin && (
            <div className="flex items-center gap-2">
              <Checkbox
                id="is_blocked"
                checked={isBlocked}
                onCheckedChange={(checked) => setIsBlocked(!!checked)}
                className="h-5 w-5 border-2 rounded-md transition-all data-[state=checked]:bg-destructive data-[state=checked]:border-destructive"
              />
              <Label htmlFor="is_blocked" className="cursor-pointer text-sm font-normal text-foreground">
                Заблокировать доступ
              </Label>
            </div>
          )}
        </div>
        <Button
          onClick={handleSave}
          disabled={isSaving}
          className="h-11 px-10 font-bold shadow-lg shadow-primary/20 transition-all active:scale-[0.98] gap-2"
        >
          <Save className="h-5 w-5" />
          {isSaving ? "Сохранение..." : "Сохранить изменения"}
        </Button>
      </div>

      {/* ALERT DIALOG ДЛЯ СМЕНЫ АВАТАРА */}
      <AlertDialog open={isAvatarDialogOpen} onOpenChange={setIsAvatarDialogOpen}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Изменение аватара</AlertDialogTitle>
            <AlertDialogDescription className="space-y-4">
              <div>Ниже представлен текущий аватар пользователя:</div>
              {avatarUrl ? (
                <div className="flex justify-center p-4">
                  <img
                    src={avatarUrl}
                    alt="Текущий аватар"
                    className="h-32 w-32 rounded-full border border-border shadow-sm object-cover"
                  />
                </div>
              ) : (
                <div className="border-2 border-dashed border-border p-6 text-center space-y-2 bg-muted/20">
                  <ImageOff className="h-10 w-10 text-muted-foreground mx-auto" />
                  <div className="text-sm font-semibold text-muted-foreground">Аватар не установлен</div>
                </div>
              )}
              
              <div className="space-y-2 text-foreground">
                <Label htmlFor="newAvatarUrl">Ссылка на новое изображение</Label>
                <Input
                  id="newAvatarUrl"
                  value={newAvatarUrlInput}
                  onChange={(e) => setNewAvatarUrlInput(e.target.value)}
                  placeholder="https://example.com/avatar.jpg"
                  className="h-10"
                />
              </div>
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter className="flex justify-between items-center w-full gap-2">
            {avatarUrl && (
              <Button variant="destructive" onClick={handleDeleteAvatar} className="mr-auto">
                Удалить аватар
              </Button>
            )}
            <div className="flex gap-2 justify-end ml-auto">
              <AlertDialogCancel onClick={() => setIsAvatarDialogOpen(false)}>Отмена</AlertDialogCancel>
              <AlertDialogAction onClick={handleSaveAvatar}>Сохранить</AlertDialogAction>
            </div>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </div>
  )
}
