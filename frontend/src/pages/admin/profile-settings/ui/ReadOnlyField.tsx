import { Label } from "@/shared/ui/label";

export function ReadOnlyField({ label, value }: { label: string, value: string | null }) {
  return (
    <div className="space-y-1.5">
      <Label className="text-xs font-semibold text-muted-foreground uppercase tracking-wider">{label}</Label>
      <div className="p-3 bg-muted/40 border border-border/30 rounded-md text-sm font-medium select-none shadow-inner opacity-70">
        {value ?? "—"}
      </div>
    </div>
  );
}
