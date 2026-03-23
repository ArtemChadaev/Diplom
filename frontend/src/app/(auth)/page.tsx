"use client";

import React from "react";
import Image from "next/image";
import { 
  Card, 
  CardContent, 
  CardDescription, 
  CardFooter, 
  CardHeader, 
  CardTitle 
} from "@/components/ui/card";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";
import { Mail, Lock } from "lucide-react";

export default function AuthPage() {
  return (
    <div className="w-full max-w-[420px] mx-auto animate-in fade-in duration-500 scale-in-95">
      <Tabs defaultValue="login" className="w-full">
        <Card className="border-border/40 shadow-xl bg-card/50 backdrop-blur-sm overflow-hidden">
          <CardHeader className="space-y-1 pb-6 text-center">
            <CardTitle className="text-2xl font-bold tracking-tight">Добро пожаловать</CardTitle>
            <CardDescription>
              Введите данные для доступа к системе управления
            </CardDescription>
            <TabsList className="grid w-full grid-cols-2 mt-4">
              <TabsTrigger value="login">Вход</TabsTrigger>
              <TabsTrigger value="register">Регистрация</TabsTrigger>
            </TabsList>
          </CardHeader>
          
          <CardContent className="space-y-4">
            <TabsContent value="login" className="space-y-4 mt-0">
              <div className="space-y-4">
                <div className="space-y-1.5 translate-y-0 group">
                  <div className="relative">
                    <Mail className="absolute left-3 top-3 h-4 w-4 text-muted-foreground transition-colors group-focus-within:text-primary" />
                    <Input 
                      placeholder="Электронная почта" 
                      type="email" 
                      className="pl-10 h-11 bg-background/50 border-border/40 focus:border-primary/50 transition-all"
                    />
                  </div>
                </div>
                <div className="space-y-1.5 group">
                  <div className="relative">
                    <Lock className="absolute left-3 top-3 h-4 w-4 text-muted-foreground transition-colors group-focus-within:text-primary" />
                    <Input 
                      placeholder="Пароль" 
                      type="password" 
                      className="pl-10 h-11 bg-background/50 border-border/40 focus:border-primary/50 transition-all"
                    />
                  </div>
                </div>
                <Button className="w-full h-11 font-medium transition-all hover:opacity-90 active:scale-[0.98]">
                  Авторизоваться
                </Button>
              </div>
            </TabsContent>

            <TabsContent value="register" className="space-y-4 mt-0">
              <div className="space-y-4">
                <div className="space-y-1.5 group">
                  <div className="relative">
                    <Mail className="absolute left-3 top-3 h-4 w-4 text-muted-foreground transition-colors group-focus-within:text-primary" />
                    <Input 
                      placeholder="Электронная почта" 
                      type="email" 
                      className="pl-10 h-11 bg-background/50 border-border/40 focus:border-primary/50 transition-all"
                    />
                  </div>
                </div>
                <div className="space-y-1.5 group">
                  <div className="relative">
                    <Lock className="absolute left-3 top-3 h-4 w-4 text-muted-foreground transition-colors group-focus-within:text-primary" />
                    <Input 
                      placeholder="Пароль" 
                      type="password" 
                      className="pl-10 h-11 bg-background/50 border-border/40 focus:border-primary/50 transition-all"
                    />
                  </div>
                </div>
                <div className="space-y-1.5 group">
                  <div className="relative">
                    <Lock className="absolute left-3 top-3 h-4 w-4 text-muted-foreground transition-colors group-focus-within:text-primary" />
                    <Input 
                      placeholder="Подтвердите пароль" 
                      type="password" 
                      className="pl-10 h-11 bg-background/50 border-border/40 focus:border-primary/50 transition-all"
                    />
                  </div>
                </div>
                <Button className="w-full h-11 font-medium transition-all hover:opacity-90 active:scale-[0.98]">
                  Зарегистрироваться
                </Button>
              </div>
            </TabsContent>

            <div className="relative py-2">
              <div className="absolute inset-0 flex items-center">
                <Separator className="w-full" />
              </div>
              <div className="relative flex justify-center text-xs uppercase">
                <span className="bg-card px-2 text-muted-foreground">Или войти через</span>
              </div>
            </div>

            <div className="grid grid-cols-2 gap-4 pb-2">
              <Button variant="outline" className="h-11 border-border/40 hover:bg-muted/50 transition-colors group">
                <svg className="w-5 h-5 mr-0 transition-transform group-hover:scale-110" viewBox="0 0 24 24" fill="currentColor">
                  <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm4.64 6.8c-.15 1.58-.8 5.42-1.13 7.19-.14.75-.42 1-.68 1.03-.58.05-1.02-.38-1.58-.75-.88-.58-1.38-.94-2.23-1.5-.99-.65-.35-1.01.22-1.59.15-.15 2.71-2.48 2.76-2.69a.2.2 0 00-.05-.18c-.06-.05-.14-.03-.21-.02-.09.02-1.49.95-4.22 2.79-.4.27-.76.41-1.08.4-.36-.01-1.04-.2-1.55-.37-.63-.2-1.13-.31-1.08-.66.02-.18.27-.36.74-.55 2.91-1.27 4.85-2.11 5.83-2.51 2.78-1.16 3.35-1.36 3.73-1.36.08 0 .27.02.39.12.1.08.13.19.14.27-.01.06.01.24 0 .38z" />
                </svg>
              </Button>
              <Button variant="outline" className="h-11 border-border/40 hover:bg-muted/50 transition-colors group">
                <svg className="w-5 h-5 mr-0 transition-transform group-hover:scale-110" viewBox="0 0 24 24">
                  <path d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z" fill="#4285F4" />
                  <path d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z" fill="#34A853" />
                  <path d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l3.66-2.84z" fill="#FBBC05" />
                  <path d="M12 5.38c1.62 0 3.06.56 4.21 1.66l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z" fill="#EA4335" />
                </svg>
              </Button>
            </div>
          </CardContent>
        </Card>
      </Tabs>
    </div>
  );
}
