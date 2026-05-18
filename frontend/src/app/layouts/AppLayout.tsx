import React from "react";
import { Outlet } from "react-router-dom";
import { Header } from "@/widgets/header";
import { Footer } from "@/widgets/footer";

export function AppLayout() {
  return (
    <div className="min-h-screen flex flex-col bg-background selection:bg-secondary/20">
      <Header />
      <main className="flex-1 w-full mx-auto p-6 md:p-10">
        <Outlet />
      </main>
      <Footer />
    </div>
  );
}
