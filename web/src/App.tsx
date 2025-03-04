import { Routes, Route, Navigate } from "react-router";

import LoadSession from "./components/LoadSession";
import { AuthProvider } from "./context/AuthProvider";
import LoginPage from "./pages/LoginPage";
import SignupPage from "./pages/SignupPage";
import ForgotPasswordPage from "./pages/ForgotPasswordPage";
import { RootAssetFolderView, RouteBasedFolderView } from "./pages/FolderView";
import { ToastContainer } from "react-toastify";

export default function App() {
  return (
    <AuthProvider>
      <Routes>
        <Route index element={<Navigate to="/home" />} />

        {/* public routes */}
        <Route path="login" element={<LoginPage />} />
        <Route path="register" element={<SignupPage />} />
        <Route path="forgot-password" element={<ForgotPasswordPage />} />

        {/* protected routes */}
        <Route element={<LoadSession />}>
          <Route path="/home" element={<RootAssetFolderView />} />
          <Route path="/folders/:assetID" element={<RouteBasedFolderView />} />
        </Route>
      </Routes>
      <ToastContainer />
    </AuthProvider>
  );
}
