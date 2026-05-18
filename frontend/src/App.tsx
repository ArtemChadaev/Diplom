import React from "react";
import { BrowserRouter, Routes, Route } from "react-router-dom";
import { AppLayout } from "./app/layouts/AppLayout";
import { AuthLayout } from "./app/layouts/AuthLayout";
import { AuthPage } from "./pages/auth";
import { DashboardPage } from "./pages/dashboard";
import { SearchPage } from "./pages/search";
import { UsersPage } from "./pages/admin/users";
import { ProfileSettingsPage } from "./pages/admin/profile-settings";

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
        </Route>
      </Routes>
    </BrowserRouter>
  );
}

export default App;
