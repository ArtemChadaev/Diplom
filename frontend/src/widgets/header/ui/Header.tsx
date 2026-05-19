import { Bell, CircleUserRound } from "lucide-react";
import { Link } from "react-router-dom";


import { SearchBar } from "@/widgets/search-bar";

import { Avatar, AvatarFallback } from "@/shared/ui/avatar";
import { Button } from "@/shared/ui/button";

export function Header() {
  return (
    <header className="sticky top-0 z-50 bg-background/70 backdrop-blur-md border-b border-border/40 px-6 py-4 min-h-[72px] flex items-center justify-between">
      <div className="flex items-center gap-8">
        <Link to="/" className="font-sans text-xl font-bold text-primary tracking-tight">
          pharma-hub
        </Link>
        <nav className="hidden md:flex items-center gap-3 pt-1">
          <Button variant="link" className="text-sm font-medium text-muted-foreground hover:text-primary transition-colors" asChild>
            <Link to="#">
              Warehouse
            </Link>
          </Button>
          <Button variant="link" className="text-sm font-medium text-muted-foreground hover:text-primary transition-colors" asChild>
            <Link to="#">
             Operations
            </Link>
          </Button>
          <Button variant="link" className="text-sm font-medium text-muted-foreground hover:text-primary transition-colors" asChild>
            <Link to="#">
              Reports
            </Link>
          </Button>
        </nav>
      </div>

      <div className="flex items-center gap-4">
        <SearchBar />
        
        <Button variant="ghost" size="icon" className="rounded-full hover:bg-muted/80 transition-colors h-10 w-10">
          <Bell className="h-5 w-5 text-muted-foreground" />
        </Button>

        {/*TODO: Сделать скрытую для обычных пользователей*/}
        <Button className="rounded-lg text-sm font-medium transition-opacity" asChild>
          <Link to="/admin/users">Admin Panel</Link>
        </Button>

        <Button variant="ghost" size="icon" className="rounded-full size-10 transition-transform" asChild>
          <Link to="/admin/profile/1/settings">
            <Avatar>
              <AvatarFallback>
                <CircleUserRound className="size-6" />
              </AvatarFallback>
            </Avatar>
          </Link>
        </Button>
      </div>
    </header>
  );
}
