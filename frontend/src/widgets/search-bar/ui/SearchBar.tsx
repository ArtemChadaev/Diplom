import { useState, Suspense } from "react"
import { useNavigate, useSearchParams } from "react-router-dom"
import { Input } from "@/shared/ui/input"
import { Button } from "@/shared/ui/button"
import { ButtonGroup } from "@/shared/ui/button-group"
import { Search } from "lucide-react"

function SearchBarInner() {
  const navigate = useNavigate()
  const [searchParams] = useSearchParams()
  const [query, setQuery] = useState(searchParams.get("q") || "")

  const handleSearch = (e?: React.FormEvent) => {
    e?.preventDefault()
    if (query.trim()) {
      const params = new URLSearchParams(searchParams.toString())
      params.set("q", query.trim())
      navigate(`/search?${params.toString()}`)
    } else {
      const params = new URLSearchParams(searchParams.toString())
      params.delete("q")
      navigate(`/search?${params.toString()}`)
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

export function SearchBar() {
  return (
    <Suspense fallback={<div className="hidden sm:flex w-[200px] h-10" />}>
      <SearchBarInner />
    </Suspense>
  )
}
