"use client"
import * as React from "react"
import { Check, ChevronsUpDown, X } from "lucide-react"

import { cn } from "@/lib/utils"
import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
} from "@/components/ui/command"
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover"
import { Badge } from "@/components/ui/badge"

export type Option = {
  label: string
  value: string
}

interface MultiSelectProps {
  options: Option[]
  selected: string[]
  onChange: (values: string[]) => void
  placeholder?: string
  badgeClassName?: string
}

export function MultiSelect({ 
  options, 
  selected, 
  onChange, 
  placeholder = "Select...",
  badgeClassName
}: MultiSelectProps) {
  const [open, setOpen] = React.useState(false)

  const handleUnselect = (item: string) => {
    onChange(selected.filter((i) => i !== item))
  }

  return (
    <Popover open={open} onOpenChange={setOpen}>
      <PopoverTrigger asChild>
        <div
          role="combobox"
          aria-expanded={open}
          className="flex w-full items-center justify-between rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50 min-h-12 cursor-pointer hover:bg-accent/50 transition-colors"
        >
          <div className="flex gap-1 flex-wrap">
            {selected.length === 0 && <span className="text-muted-foreground py-1">{placeholder}</span>}
            {selected.map((item) => (
              <Badge
                variant="secondary"
                key={item}
                className={cn("mr-1 mb-1 mt-1 font-normal", badgeClassName)}
                onClick={(e) => {
                  e.stopPropagation()
                  handleUnselect(item)
                }}
              >
                {options.find((opt) => opt.value === item)?.label || item}
                <X className="ml-1 h-3 w-3 hover:text-destructive cursor-pointer shrink-0" />
              </Badge>
            ))}
          </div>
          <ChevronsUpDown className="ml-2 h-4 w-4 shrink-0 opacity-50" />
        </div>
      </PopoverTrigger>
      <PopoverContent className="w-[var(--radix-popover-trigger-width)] p-0 bg-background" align="start">
        <Command>
          <CommandInput placeholder="Поиск..." />
          <CommandList>
            <CommandEmpty>Ничего не найдено.</CommandEmpty>
            <CommandGroup>
              {options.map((option) => (
                <CommandItem
                  key={option.value}
                  onSelect={() => {
                    onChange(
                      selected.includes(option.value)
                        ? selected.filter((item) => item !== option.value)
                        : [...selected, option.value]
                    )
                  }}
                >
                  <Check
                    className={cn(
                      "mr-2 h-4 w-4 shrink-0",
                      selected.includes(option.value) ? "opacity-100" : "opacity-0"
                    )}
                  />
                  {option.label}
                </CommandItem>
              ))}
            </CommandGroup>
          </CommandList>
        </Command>
      </PopoverContent>
    </Popover>
  )
}
