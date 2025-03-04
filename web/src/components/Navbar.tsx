import { Search, User, LogOut } from "lucide-react";
import { useAuth } from "../context/AuthProvider";
import { Link, useNavigate } from "react-router";

export default function Navbar() {
  const { isLoggedIn, clearAuthContext } = useAuth();
  const navigate = useNavigate();

  function onLogout() {
    if (!isLoggedIn()) return;
    clearAuthContext();
    navigate("/");
  }

  return (
    <nav className="bg-white shadow-md px-6 py-3 flex items-center justify-between">
      <div className="flex items-center space-x-2">
        <span className="text-2xl font-semibold text-blue-600">
          Drive Clone
        </span>
      </div>

      {isLoggedIn() && (
        <div className="flex-1 max-w-2xl mx-8">
          <div className="relative">
            <input
              type="text"
              placeholder="Search in Drive"
              className="w-full px-4 py-2 pl-10 bg-gray-100 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
            <Search className="absolute left-3 top-2.5 text-gray-400 h-5 w-5" />
          </div>
        </div>
      )}

      <div className="flex items-center space-x-4">
        {!isLoggedIn() ? (
          <div className="flex space-x-2">
            <Link
              to="/login"
              className="px-4 py-2 text-gray-700 hover:bg-gray-100 rounded-md"
            >
              Sign In
            </Link>
            <Link
              to="/register"
              className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700"
            >
              Sign Up
            </Link>
          </div>
        ) : (
          <div className="flex items-center space-x-2">
            <button className="p-2 hover:bg-gray-100 rounded-full">
              <User className="h-6 w-6 text-gray-700" />
            </button>
            <button
              onClick={onLogout}
              className="p-2 hover:bg-gray-100 rounded-full"
              title="Sign out"
            >
              <LogOut className="h-6 w-6 text-gray-700" />
            </button>
          </div>
        )}
      </div>
    </nav>
  );
}
