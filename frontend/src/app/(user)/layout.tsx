import Link from "next/link";

export default function EmployeeLayout({ children }: { children: React.ReactNode }) {
    return (
        <div style={{ display: "flex", minHeight: "100vh" }}>
            {/* Sidebar */}
            <aside style={{ width: 220, background: "#1e3a8a", color: "#fff", display: "flex", flexDirection: "column", padding: "24px 16px", gap: 8 }}>
                <div style={{ fontWeight: 700, fontSize: 18, marginBottom: 24 }}>Сотрудник</div>
                <Link href="/employee" style={{ color: "#bfdbfe", padding: "8px 12px", borderRadius: 6 }}>Главная</Link>
                <Link href="/employee/history" style={{ color: "#bfdbfe", padding: "8px 12px", borderRadius: 6 }}>История</Link>
                <Link href="/employee/profile" style={{ color: "#bfdbfe", padding: "8px 12px", borderRadius: 6 }}>Профиль</Link>
                <div style={{ marginTop: "auto" }}>
                    <Link href="/" style={{ color: "#93c5fd", fontSize: 13 }}>← На главную</Link>
                </div>
            </aside>
            {/* Content */}
            <main style={{ flex: 1, padding: 32 }}>{children}</main>
        </div>
    );
}
