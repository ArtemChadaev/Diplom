import { Header } from "@/components/header";
import { Footer } from "@/components/footer";

export default function UserLayout({ children }: { children: React.ReactNode }) {
    return (
        <div className="min-h-screen flex flex-col bg-background selection:bg-secondary/20">
            <Header />
            <main className="flex-1 w-full mx-auto p-6 md:p-10">
                {children}
            </main>
            <Footer />
        </div>
    );
}
