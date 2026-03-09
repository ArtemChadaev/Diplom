import Link from "next/link";

export default function AdminLayout({ children }: { children: React.ReactNode }) {
    return (
        <div style={{ display: "flex", minHeight: "100vh" }}>
            {/* Sidebar */}
            <aside style={{ width: 220, background: "#14532d", color: "#fff", display: "flex", flexDirection: "column", padding: "24px 16px", gap: 8 }}>
                <div style={{ fontWeight: 700, fontSize: 18, marginBottom: 24 }}>Администратор</div>
                <Link href="/admin" style={{ color: "#bbf7d0", padding: "8px 12px", borderRadius: 6 }}>Главная</Link>
                <Link href="/admin/users" style={{ color: "#bbf7d0", padding: "8px 12px", borderRadius: 6 }}>Пользователи</Link>
                <Link href="/admin/settings" style={{ color: "#bbf7d0", padding: "8px 12px", borderRadius: 6 }}>Настройки</Link>
                <div style={{ marginTop: "auto" }}>
                    <Link href="/" style={{ color: "#86efac", fontSize: 13 }}>← На главную</Link>
                </div>
            </aside>
            {/* Content */}
            <main style={{ flex: 1, padding: 32 }}>{children}</main>
        </div>
    );
}
