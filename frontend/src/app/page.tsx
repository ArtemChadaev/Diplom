import Link from "next/link";

export default function Home() {
  return (
    <main style={{ display: "flex", flexDirection: "column", alignItems: "center", justifyContent: "center", minHeight: "100vh", gap: 24 }}>
      <h1 style={{ fontSize: 28, fontWeight: 700 }}>Diplom App</h1>
      <div style={{ display: "flex", gap: 16 }}>
        <Link href="/employee" style={{ padding: "12px 28px", background: "#2563eb", color: "#fff", borderRadius: 8, fontWeight: 600 }}>
          Сотрудник
        </Link>
        <Link href="/admin" style={{ padding: "12px 28px", background: "#16a34a", color: "#fff", borderRadius: 8, fontWeight: 600 }}>
          Администратор
        </Link>
      </div>
    </main>
  );
}
