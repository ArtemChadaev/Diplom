import { Bell } from "lucide-react"

import { SearchBar } from "@/widgets/search-bar"

import { Button } from "@/shared/ui/button"
import { SidebarTrigger } from "@/shared/ui/sidebar"

export function Header() {
  return (
    <header className="sticky top-0 z-50 bg-background/70 backdrop-blur-md border-b border-border/40 px-6 py-4 min-h-[72px] flex items-center justify-between">
      <div className="flex items-center gap-4">
        <SidebarTrigger className="rounded-none shadow-none" />
      </div>

      <div className="flex items-center gap-4">
        <SearchBar />
        
        <Button variant="ghost" size="icon" className="rounded-full hover:bg-muted/80 transition-colors h-10 w-10">
          <Bell className="h-5 w-5 text-muted-foreground" />
        </Button>
      </div>
    </header>
  )
}
