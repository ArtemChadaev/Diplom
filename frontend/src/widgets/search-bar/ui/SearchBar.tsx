import { Search } from "lucide-react"
import { Suspense, useState } from "react"
import { useNavigate, useSearchParams } from "react-router-dom"

import { Button } from "@/shared/ui/button"
import { ButtonGroup } from "@/shared/ui/button-group"
import { Input } from "@/shared/ui/input"

function SearchBarInner() {
  const navigate = useNavigate()
  const [searchParams] = useSearchParams()
  const [query, setQuery] = useState(searchParams.get("q") ?? "")

  const handleSearch = (e?: React.SyntheticEvent) => {
    e?.preventDefault()
    const params = new URLSearchParams(searchParams.toString())
    if (query.trim()) {
      params.set("q", query.trim())
    } else {
      params.delete("q")
    }
    void navigate(`/search?${params.toString()}`)
  }

  return (
    <form onSubmit={handleSearch} className="hidden sm:flex items-center">
      <ButtonGroup>
        <Input
          placeholder="Search..."
          value={query}
          onChange={(e) => { setQuery(e.target.value); }}
        />
        <Button variant="outline" type="submit">
          <Search className="h-4 w-4" />
        </Button>
      </ButtonGroup>
    </form>
  )
}

export function SearchBar() {
  return (
    <Suspense fallback={<div className="hidden sm:flex w-[200px] h-10" />}>
      <SearchBarInner />
    </Suspense>
  )
}
