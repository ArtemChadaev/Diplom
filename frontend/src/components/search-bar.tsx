"use client"

import { useState } from "react"
import { useRouter, useSearchParams } from "next/navigation"
import { Input } from "@/components/ui/input"
import { Button } from "@/components/ui/button"
import { ButtonGroup } from "@/components/ui/button-group"
import { Search } from "lucide-react"

export function SearchBar() {
  const router = useRouter()
  const searchParams = useSearchParams()
  const [query, setQuery] = useState(searchParams.get("q") || "")

  const handleSearch = (e?: React.FormEvent) => {
    e?.preventDefault()
    if (query.trim()) {
      const params = new URLSearchParams(searchParams.toString())
      params.set("q", query.trim())
      router.push(`/search?${params.toString()}`)
    } else {
      const params = new URLSearchParams(searchParams.toString())
      params.delete("q")
      router.push(`/search?${params.toString()}`)
    }
  }

  return (
    <form onSubmit={handleSearch} className="hidden sm:flex items-center">
      <ButtonGroup>
        <Input
          placeholder="Search..."
          value={query}
          onChange={(e) => setQuery(e.target.value)}
        />
        <Button variant="outline" type="submit">
          <Search className="h-4 w-4" />
        </Button>
      </ButtonGroup>
    </form>
  )
}
