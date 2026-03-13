import { ReactNode } from "react";
import { cn } from "@/lib/utils";

export type KpiCardState = "normal" | "warning" | "error" | "success";

interface KpiCardProps {
    title: string;
    value: string;
    description: string;
    icon: ReactNode;
    state?: KpiCardState;
    trendIcon?: ReactNode;
}

export function KpiCard({
    title,
    value,
    description,
    icon,
    state = "normal",
    trendIcon
}: KpiCardProps) {

    const styles = {
        normal: {
            wrapper: "bg-card border border-l-4 border-l-muted-foreground/30",
            iconBg: "bg-muted",
            iconColor: "text-primary",
            title: "text-muted-foreground/60",
            value: "text-primary",
            desc: "text-muted-foreground",
        },
        warning: {
            wrapper: "bg-card border border-l-4 border-l-yellow-500 shadow-sm",
            iconBg: "bg-yellow-100",
            iconColor: "text-yellow-700",
            title: "text-muted-foreground/60",
            value: "text-yellow-700",
            desc: "text-muted-foreground",
        },
        error: {
            wrapper: "bg-card border border-l-4 border-l-destructive shadow-sm",
            iconBg: "bg-destructive/10",
            iconColor: "text-destructive",
            title: "text-destructive/60",
            value: "text-destructive",
            desc: "text-destructive/80",
        },
        success: {
            wrapper: "bg-card border border-l-4 border-l-green-500 shadow-sm",
            iconBg: "bg-green-500/10",
            iconColor: "text-green-600",
            title: "text-green-600/60",
            value: "text-green-600",
            desc: "text-green-600/80",
        }
    };

    const currentStyle = styles[state];

    return (
        <div className={cn("md:col-span-4 p-6 rounded-xl flex flex-col justify-between group transition-all duration-300", currentStyle.wrapper)}>
            <div className="flex justify-between items-start mb-4">
                <div className={cn("p-2 rounded-lg", currentStyle.iconBg, currentStyle.iconColor)}>
                    {icon}
                </div>
                <span className={cn("text-xs font-bold uppercase tracking-widest", currentStyle.title)}>
                    {title}
                </span>
            </div>
            <div>
                <div className={cn("font-sans text-5xl font-bold mb-1 tracking-tight", currentStyle.value)}>
                    {value}
                </div>
                <div className={cn("text-sm font-medium flex items-center gap-1", currentStyle.desc)}>
                    {trendIcon}
                    {description}
                </div>
            </div>
        </div>
    );
}
