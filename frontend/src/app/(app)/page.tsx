import { CalendarDays, TrendingUp, AlertTriangle, AlertCircle, PlusSquare, MinusSquare, AlertOctagon, CheckCircle2, ChevronRight, Package, ArrowUpRight, ArrowDownRight } from "lucide-react";
import { KpiCard } from "./components/kpi-card";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";

//TODO: НЕ ЗАБЫТЬ СДЕЛАТЬ СКЕЛЕТ И НЕ ПОКАЗЫВАТЬ ДАННЫЕ НЕ НОВЫЕ
export default function EmployeePage() {
    return (
        <div className="flex flex-col gap-10">
            {/* Bento Grid Layout */}
            <div className="grid grid-cols-1 md:grid-cols-12 gap-5">
                
                {/* KPI Cards (Top Row) */}
                <KpiCard
                    title="Total Items"
                    value="12,482"
                    description="+4% от прошлого месяца"
                    icon={<Package className="h-5 w-5" />}
                    state="normal"
                    trendIcon={<TrendingUp className="h-4 w-4" />}
                />

                <KpiCard
                    title="Shortage"
                    value="24"
                    description="Позиции требуют пополнения"
                    icon={<AlertTriangle className="h-5 w-5" />}
                    state="error"
                />

                <KpiCard
                    title="Expiring"
                    value="112"
                    description="Срок годности < 30 дней"
                    icon={<AlertCircle className="h-5 w-5" />}
                    state="warning"
                />

                {/* Quick Actions (Center Left) */}
                <Card className="md:col-span-5 bg-muted/40 flex flex-col gap-4 border-none shadow-none">
                    <CardHeader className="pb-2">
                        <CardTitle className="text-xs font-bold uppercase tracking-widest text-muted-foreground">Быстрые операции</CardTitle>
                    </CardHeader>
                    <CardContent className="flex flex-col gap-4">
                        <button className="bg-gradient-to-br from-primary to-primary/80 group w-full p-6 rounded-xl flex items-center justify-between text-left hover:shadow-lg transition-all cursor-pointer">
                            <div className="flex items-center gap-4">
                                <div className="bg-white/10 p-3 rounded-lg">
                                    <PlusSquare className="text-white h-6 w-6" />
                                </div>
                                <div>
                                    <span className="block text-white font-semibold text-lg leading-tight">+ Оформить поступление</span>
                                    <span className="text-white/60 text-xs">Приемка новой партии товаров</span>
                                </div>
                            </div>
                            <ChevronRight className="text-white/40 group-hover:translate-x-1 transition-transform h-5 w-5" />
                        </button>
                        
                        <button className="bg-card group w-full p-6 rounded-xl flex items-center justify-between text-left border border-transparent hover:border-border hover:shadow-sm transition-all cursor-pointer">
                            <div className="flex items-center gap-4">
                                <div className="bg-muted p-3 rounded-lg text-primary">
                                    <MinusSquare className="h-6 w-6" />
                                </div>
                                <div>
                                    <span className="block text-primary font-semibold text-lg leading-tight">– Оформить списание</span>
                                    <span className="text-muted-foreground text-xs">Инвентаризация или брак</span>
                                </div>
                            </div>
                            <ChevronRight className="text-muted-foreground/40 group-hover:translate-x-1 transition-transform h-5 w-5" />
                        </button>
                    </CardContent>
                </Card>

                {/* Reorder Tasks (Center Right) */}
                <Card className="md:col-span-7 relative overflow-hidden">
                    <div className="absolute top-0 right-0 w-32 h-32 bg-secondary/5 rounded-full -mr-16 -mt-16"></div>
                    <CardHeader className="flex flex-row justify-between items-center mb-2 z-10 relative">
                        <CardTitle className="text-xs font-bold uppercase tracking-widest text-muted-foreground">Уведомления и задачи</CardTitle>
                        <span className="bg-secondary/10 text-secondary px-2 py-0.5 rounded text-[10px] font-bold">1 ACTIVE</span>
                    </CardHeader>
                    
                    <CardContent className="space-y-4 relative z-10">
                        {/* Task Item */}
                        <div className="flex flex-col md:flex-row items-start md:items-center justify-between gap-4 p-5 bg-muted/30 rounded-xl border">
                            <div className="flex items-start gap-4">
                                <div className="w-10 h-10 rounded-full bg-destructive/10 flex items-center justify-center text-destructive shrink-0">
                                    <AlertOctagon className="h-5 w-5" />
                                </div>
                                <div>
                                    <h4 className="font-bold text-primary">Reorder needed: Aspirin</h4>
                                    <p className="text-sm text-muted-foreground">Stock level reached zero. Critical for Heart Ward.</p>
                                    <div className="mt-2 flex items-center gap-4">
                                        <span className="text-xs font-semibold px-2 py-0.5 bg-muted rounded text-muted-foreground">SKU: ASP-442</span>
                                        <span className="text-xs font-bold text-destructive">Stock: 0</span>
                                    </div>
                                </div>
                            </div>
                            <button className="whitespace-nowrap bg-secondary text-primary-foreground px-5 py-2.5 rounded-lg text-sm font-semibold hover:bg-secondary/90 transition-colors shadow-lg shadow-secondary/20">
                                Create Order Request
                            </button>
                        </div>

                        {/* Completed Task Item */}
                        <div className="flex items-center gap-4 p-5 opacity-60 grayscale">
                            <div className="w-10 h-10 rounded-full bg-muted flex items-center justify-center text-muted-foreground shrink-0">
                                <CheckCircle2 className="h-5 w-5" />
                            </div>
                            <div className="flex-1 border-b border-border/50 pb-2">
                                <h4 className="font-bold text-primary">Insulin restock completed</h4>
                                <p className="text-xs text-muted-foreground">Successfully received 500 units by User @ivanov_aa</p>
                            </div>
                        </div>
                    </CardContent>
                </Card>

                {/* Recent Activity Table (Bottom Wide) */}
                <Card className="md:col-span-12 overflow-hidden shadow-sm">
                    <CardHeader className="flex flex-row items-center justify-between">
                        <div>
                            <CardTitle className="text-xs font-bold uppercase tracking-widest text-muted-foreground mb-1">Последняя активность</CardTitle>
                            <p className="text-sm text-muted-foreground/60">История операций за последние 24 часа</p>
                        </div>
                        <button className="text-secondary text-sm font-bold flex items-center gap-1 hover:underline">
                            View full operation log
                            <ChevronRight className="h-4 w-4" />
                        </button>
                    </CardHeader>
                    
                    <CardContent className="p-0">
                        <Table>
                            <TableHeader>
                                <TableRow className="hover:bg-transparent">
                                    <TableHead className="px-6 py-4 text-xs font-bold uppercase tracking-wider h-auto">Time</TableHead>
                                    <TableHead className="px-6 py-4 text-xs font-bold uppercase tracking-wider h-auto">User</TableHead>
                                    <TableHead className="px-6 py-4 text-xs font-bold uppercase tracking-wider h-auto">Operation</TableHead>
                                    <TableHead className="px-6 py-4 text-xs font-bold uppercase tracking-wider h-auto">Drug Name</TableHead>
                                    <TableHead className="px-6 py-4 text-xs font-bold uppercase tracking-wider h-auto text-right">Quantity</TableHead>
                                </TableRow>
                            </TableHeader>
                            <TableBody>
                                {/* Table Row 1 */}
                                <TableRow>
                                    <TableCell className="px-6 py-4 text-sm font-medium text-muted-foreground">14:22</TableCell>
                                    <TableCell className="px-6 py-4">
                                        <div className="flex items-center gap-2">
                                            <div className="w-6 h-6 rounded-full bg-primary text-[10px] text-primary-foreground flex items-center justify-center">AS</div>
                                            <span className="text-sm font-medium text-primary">А. Сидоров</span>
                                        </div>
                                    </TableCell>
                                    <TableCell className="px-6 py-4">
                                        <span className="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-[11px] font-bold bg-secondary/10 text-secondary uppercase">
                                            <ArrowDownRight className="w-3 h-3" />
                                            Incoming
                                        </span>
                                    </TableCell>
                                    <TableCell className="px-6 py-4 text-sm font-semibold text-primary">Paracetamol 500mg</TableCell>
                                    <TableCell className="px-6 py-4 text-sm font-bold text-primary text-right">+2,000</TableCell>
                                </TableRow>
                                {/* Table Row 2 */}
                                <TableRow>
                                    <TableCell className="px-6 py-4 text-sm font-medium text-muted-foreground">13:05</TableCell>
                                    <TableCell className="px-6 py-4">
                                        <div className="flex items-center gap-2">
                                            <div className="w-6 h-6 rounded-full bg-primary text-[10px] text-primary-foreground flex items-center justify-center">ИВ</div>
                                            <span className="text-sm font-medium text-primary">И. Волков</span>
                                        </div>
                                    </TableCell>
                                    <TableCell className="px-6 py-4">
                                        <span className="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-[11px] font-bold bg-destructive/10 text-destructive uppercase">
                                            <ArrowUpRight className="w-3 h-3" />
                                            Outgoing
                                        </span>
                                    </TableCell>
                                    <TableCell className="px-6 py-4 text-sm font-semibold text-primary">Ibuprofen 200mg</TableCell>
                                    <TableCell className="px-6 py-4 text-sm font-bold text-primary text-right">-450</TableCell>
                                </TableRow>
                            </TableBody>
                        </Table>
                    </CardContent>
                </Card>

            </div>
        </div>
    );
}
