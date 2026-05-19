import { BrowserRouter, Routes, Route } from "react-router-dom";

import { AppLayout } from "./app/layouts/AppLayout";
import { AuthLayout } from "./app/layouts/AuthLayout";
import { ProfileSettingsPage } from "./pages/admin/profile-settings";
import { UsersPage } from "./pages/admin/users";
import { AuthPage } from "./pages/auth";
import { DashboardPage } from "./pages/dashboard";
import { NotFoundPage } from "./pages/not-found";
import { SearchPage } from "./pages/search";

export function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route element={<AuthLayout />}>
          <Route path="/auth" element={<AuthPage />} />
        </Route>

        <Route element={<AppLayout />}>
          <Route path="/" element={<DashboardPage />} />
          <Route path="/search" element={<SearchPage />} />
          <Route path="/admin/users" element={<UsersPage />} />
          <Route path="/admin/profile/:id/settings" element={<ProfileSettingsPage />} />
          <Route path="/me/settings" element={<ProfileSettingsPage />} />
          <Route path="*" element={<NotFoundPage />} />
        </Route>
      </Routes>
    </BrowserRouter>
  );
}

export default App;
