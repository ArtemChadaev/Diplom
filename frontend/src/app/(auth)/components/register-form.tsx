"use client";

import React from "react";
import { Input } from "@/components/ui/input";
import { Mail, Lock } from "lucide-react";

export function RegisterForm() {
  return (
    <div className="space-y-4 animate-in fade-in slide-in-from-right-4 duration-300">
      <div className="space-y-1.5 group">
        <div className="relative">
          <Mail className="absolute left-3 top-3 h-4 w-4 text-muted-foreground transition-colors group-focus-within:text-primary" />
          <Input 
            placeholder="Электронная почта" 
            type="email" 
            className="pl-10 h-11 bg-background/50 border-border/40 focus:border-primary/50 transition-all text-base"
          />
        </div>
      </div>
      <div className="space-y-1.5 group">
        <div className="relative">
          <Lock className="absolute left-3 top-3 h-4 w-4 text-muted-foreground transition-colors group-focus-within:text-primary" />
          <Input 
            placeholder="Пароль" 
            type="password" 
            className="pl-10 h-11 bg-background/50 border-border/40 focus:border-primary/50 transition-all text-base"
          />
        </div>
      </div>
      <div className="space-y-1.5 group">
        <div className="relative">
          <Lock className="absolute left-3 top-3 h-4 w-4 text-muted-foreground transition-colors group-focus-within:text-primary" />
          <Input 
            placeholder="Подтвердите пароль" 
            type="password" 
            className="pl-10 h-11 bg-background/50 border-border/40 focus:border-primary/50 transition-all text-base"
          />
        </div>
      </div>
    </div>
  );
}
