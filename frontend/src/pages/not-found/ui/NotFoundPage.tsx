import { ArrowLeft, Home, FileQuestion } from "lucide-react";
import { useNavigate } from "react-router-dom";

import { Button } from "@/shared/ui/button";

export function NotFoundPage() {
  const navigate = useNavigate();

  return (
    <div className="flex min-h-[70vh] flex-col items-center justify-center px-6 py-12 text-center animate-in fade-in duration-500">
      {/* 404 Visual Indicator */}
      <div className="relative mb-6 flex h-32 w-32 items-center justify-center rounded-full bg-primary/5 dark:bg-primary/10">
        <div className="absolute inset-0 animate-pulse rounded-full bg-primary/5" />
        <FileQuestion className="h-16 w-16 text-primary" />
      </div>

      {/* Hero 404 Text */}
      <h1 className="bg-gradient-to-r from-primary to-primary/70 bg-clip-text text-8xl font-black tracking-tight text-transparent select-none">
        404
      </h1>

      {/* Title */}
      <h2 className="mt-4 text-2xl font-bold tracking-tight text-foreground sm:text-3xl">
        Страница не найдена
      </h2>

      {/* Description */}
      <p className="mt-3 max-w-md text-sm text-muted-foreground leading-relaxed">
        Запрашиваемый вами адрес не существует, был изменен или временно недоступен. Проверьте правильность ввода URL или воспользуйтесь кнопками ниже.
      </p>

      {/* Action Buttons */}
      <div className="mt-8 flex flex-col sm:flex-row items-center gap-3">
        <Button
          variant="outline"
          size="lg"
          onClick={() => navigate(-1)}
          className="w-full sm:w-auto h-10 gap-2 border-border/80 hover:bg-muted/50 cursor-pointer transition-all"
        >
          <ArrowLeft className="h-4 w-4" />
          Вернуться назад
        </Button>

        <Button
          variant="default"
          size="lg"
          onClick={() => navigate("/")}
          className="w-full sm:w-auto h-10 gap-2 cursor-pointer shadow-sm transition-all"
        >
          <Home className="h-4 w-4" />
          На главную
        </Button>
      </div>
    </div>
  );
}
