import Link from "next/link";
import { Button } from "@/components/ui/button";

export function Footer() {
  return (
    <footer className="mt-8 border-t border-border/40 py-8 px-10 bg-card">
      <div className="mx-auto flex flex-col md:flex-row justify-between items-start gap-12">
        
        {/* Project Description */}
        <div className="flex flex-col gap-4 max-w-sm">
          <span className="text-xs font-bold uppercase text-primary tracking-widest">
            MedicamentHouse
          </span>
          <p className="text-sm text-muted-foreground leading-relaxed">
            Премиальная система управления медицинским складом. Точность, контроль и безопасность на каждом этапе логистики.
          </p>
        </div>

        {/* Links (Система) */}
        <div className="flex flex-col gap-4">
          <h4 className="text-xs font-bold uppercase text-primary tracking-widest">Система</h4>
          <div className="flex gap-8">
            <div className="flex flex-col gap-2 items-start">
              <Button variant="link" className="text-sm text-muted-foreground justify-start h-auto py-0 px-0 hover:text-secondary transition-colors" asChild>
                <Link href="#">Поддержка</Link>
              </Button>
              <Button variant="link" className="text-sm text-muted-foreground justify-start h-auto py-0 px-0 hover:text-secondary transition-colors" asChild>
                <Link href="#">Документация</Link>
              </Button>
            </div>
            <div className="flex flex-col gap-2 items-start">
              <Button variant="link" className="text-sm text-muted-foreground justify-start h-auto py-0 px-0 hover:text-secondary transition-colors" asChild>
                <Link href="#">Политика конфиденциальности</Link>
              </Button>
              <Button variant="link" className="text-sm text-muted-foreground justify-start h-auto py-0 px-0 hover:text-secondary transition-colors" asChild>
                <Link href="#">Условия использования</Link>
              </Button>
            </div>
          </div>
        </div>

        {/* Service Contacts */}
        <div className="flex flex-col gap-4">
          <h4 className="text-xs font-bold uppercase text-primary tracking-widest">Контакты службы</h4>
          <div className="flex flex-col gap-2">
            <p className="text-sm text-muted-foreground">Экстренная линия: +7 (800) 555-01-99</p>
            <p className="text-sm text-muted-foreground">Email: logistics@medicamenthouse.ru</p>
          </div>
        </div>

      </div>

      {/* Divider and Bottom Section */}
      <div className="mt-6 pt-5 border-t border-border/40 flex flex-col md:flex-row justify-between items-center gap-6">
        <p className="text-xs text-muted-foreground">
          © 2024 MedicamentHouse. All rights reserved.
        </p>
        
        <div className="flex items-center gap-8">
          {/* Server Status (One line) */}
          <div className="bg-muted/40 px-3 py-1.5 rounded-md flex items-center gap-2 border border-border/30">
            <span className="w-1.5 h-1.5 rounded-full bg-amber-500 animate-pulse"></span>
            <div className="flex items-baseline gap-1.5">
              <span className="text-[10px] font-bold text-muted-foreground uppercase leading-none">
                Status:
              </span>
              <span className="text-xs font-bold text-amber-500 leading-none">
                Operational
              </span>
            </div>
          </div>
        </div>
      </div>
    </footer>
  );
}
