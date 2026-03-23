import Link from "next/link";
import { Bell } from "lucide-react";
import { Button } from "@/components/ui/button";

export default function AuthLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <div className="min-h-screen flex flex-col bg-background">
      {/* Reduced Header */}
      <header className="sticky top-0 z-50 bg-background/70 backdrop-blur-md border-b border-border/40 px-6 py-4 min-h-[72px] flex items-center justify-between">
        <div className="flex items-center gap-8">
          <Link href="/" className="font-sans text-xl font-bold text-primary tracking-tight">
            MedicamentHouse
          </Link>
        </div>

        <div className="flex items-center gap-4">
          <Button variant="ghost" size="icon" className="rounded-full hover:bg-muted/80 transition-colors h-10 w-10">
            <Bell className="h-5 w-5 text-muted-foreground" />
          </Button>
        </div>
      </header>

      {/* Main Content */}
      <main className="flex-1 flex items-center justify-center p-6">
        {children}
      </main>

      {/* Reduced Footer */}
      <footer className="border-t border-border/40 py-6 px-10 bg-card">
        <div className="mx-auto flex flex-col md:flex-row justify-between items-center gap-4">
          <p className="text-xs text-muted-foreground">
            © 2024 MedicamentHouse. All rights reserved.
          </p>
          
          <Button variant="link" className="text-xs text-muted-foreground hover:text-primary transition-colors h-auto p-0" asChild>
            <Link href="#">
              Вернуться на главный сайт
            </Link>
          </Button>
        </div>
      </footer>
    </div>
  );
}
