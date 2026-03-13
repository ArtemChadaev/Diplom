import Link from "next/link";
import { Button } from "@/components/ui/button";

export function Footer() {
  return (
    <footer className="mt-20 border-t border-border/40 py-6 px-10 bg-card">
      <div className="mx-auto flex flex-col md:flex-row justify-between items-center gap-8">
        
        <div className="flex flex-col items-center md:items-start gap-2">
          <span className="font-sans text-lg font-bold text-primary">
            MedicamentHouse
          </span>
          <p className="text-xs text-muted-foreground max-w-xs text-center md:text-left">
            © 2026 ПК «Прибор» — Система управления MedicamentHouse. Версия 1.0.0
          </p>
        </div>

        <div className="flex gap-20">
          <div className="flex gap-10">
            <div className="flex flex-col gap-1">
              <Button variant="link" className="text-sm text-muted-foreground justify-start h-auto py-1 px-4" asChild>
                <Link href="#">Support</Link>
              </Button>
              <Button variant="link" className="text-sm text-muted-foreground justify-start h-auto py-1 px-4" asChild>
                <Link href="#">Documentation</Link>
              </Button>
            </div>

            <div className="flex flex-col gap-1">
              <Button variant="link" className="text-sm text-muted-foreground justify-start h-auto py-1 px-4" asChild>
                <Link href="#">Privacy Policy</Link>
              </Button>
              <Button variant="link" className="text-sm text-muted-foreground justify-start h-auto py-1 px-4" asChild>
                <Link href="#">Terms of Service</Link>
              </Button>
            </div>
          </div>

          {/*TODO*/}
          <div className="bg-muted/50 p-4 rounded-lg flex items-center gap-4">
            <div className="text-right">
              <div className="text-[10px] font-bold text-muted-foreground uppercase">
                Server Status
              </div>
              <div className="text-sm font-bold flex items-center gap-1 justify-end text-amber-500">
                <span className="w-2 h-2 rounded-full bg-amber-500"></span>
                Operational
              </div>
            </div>
          </div>
        </div>
      </div>
    </footer>
  );
}
