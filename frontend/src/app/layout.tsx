import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "Diplom App",
  description: "Employee & Admin portal",
};

export default function RootLayout({
  children,
}: Readonly<{ children: React.ReactNode }>) {
  return (
    <html lang="ru">
      <body>{children}</body>
    </html>
  );
}
