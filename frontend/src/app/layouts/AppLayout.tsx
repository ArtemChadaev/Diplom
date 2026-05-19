import { Outlet } from "react-router-dom";

import { Footer } from "@/widgets/footer";
import { Header } from "@/widgets/header";

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
