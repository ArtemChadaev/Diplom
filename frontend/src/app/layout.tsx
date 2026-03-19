import type { Metadata } from "next";
import "./globals.css";
import { Inter } from "next/font/google";
import { cn } from "@/lib/utils";
import {NuqsAdapter} from "nuqs/adapters/next/app";

const inter = Inter({subsets:['latin'],variable:'--font-sans'});

export const metadata: Metadata = {
  title: "Diplom App",
  description: "Employee & Admin portal",
};

export default function RootLayout({
  children,
}: Readonly<{ children: React.ReactNode }>) {
  return (
    <html lang="ru" className={cn("font-sans", inter.variable)}>
      <NuqsAdapter>
        <body>{children}</body>
      </NuqsAdapter>
    </html>
  );
}
