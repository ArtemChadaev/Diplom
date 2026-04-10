import type { Metadata } from "next";
import "./globals.css";
import { Plus_Jakarta_Sans, Inter } from "next/font/google";
import { cn } from "@/lib/utils";
import { NuqsAdapter } from "nuqs/adapters/next/app";

const inter = Inter({subsets:['latin'],variable:'--font-sans'});

const plusJakartaSans = Plus_Jakarta_Sans({
  subsets: ["latin"],
  variable: "--font-plus-jakarta",
  weight: ["300", "400", "500", "600", "700"],
  display: "swap",
});

export const metadata: Metadata = {
  title: "MedicamentHouse — Управление фармацевтическим складом",
  description:
    "ERP-система для управления медицинским складом: учёт партий, FEFO, инвентаризация, операции.",
};

export default function RootLayout({
  children,
}: Readonly<{ children: React.ReactNode }>) {
  return (
    <html lang="ru" className={cn( plusJakartaSans.variable, "font-sans", inter.variable)}>
      <NuqsAdapter>
        <body>{children}</body>
      </NuqsAdapter>
    </html>
  );
}
