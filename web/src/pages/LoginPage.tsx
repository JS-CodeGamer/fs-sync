import React, { useState } from "react";
import { Eye, EyeOff } from "lucide-react";
import { Link, useNavigate, useSearchParams } from "react-router";

import { useAuth } from "../context/AuthProvider";
import AuthLayout from "../components/AuthLayout";
import backend from "../helpers/backend";
import { Tokens } from "../models/user";
import { toast } from "react-toastify";
import { handleError } from "../helpers/errorHandler";

export default function LoginPage() {
  const [searchParams, _] = useSearchParams();
  const { setTokens } = useAuth();
  const navigate = useNavigate();

  const [showPassword, setShowPassword] = useState(false);
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");

  async function loginUser({
    username,
    password,
  }: {
    username: string;
    password: string;
  }) {
    try {
      const res = await backend.post<Tokens>("/login", {
        username,
        password,
      });
      if (res) {
        setTokens(res?.data);
        toast.success("Login Success!");
        navigate(decodeURIComponent(searchParams.get("return") ?? "/"));
      } else {
        toast.error("Login Failed!!!");
      }
    } catch (error) {
      handleError(error);
    }
  }

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    loginUser({ username, password });
  };

  return (
    <AuthLayout
      title="Welcome back"
      subtitle="Sign in to your account to continue"
    >
      <form className="space-y-6" onSubmit={handleSubmit}>
        <div>
          <label
            htmlFor="username"
            className="block text-sm font-medium text-gray-700"
          >
            Username
          </label>
          <div className="mt-1">
            <input
              id="username"
              name="username"
              autoComplete="username"
              required
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              className="appearance-none block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm placeholder-gray-400 focus:outline-none focus:ring-blue-500 focus:border-blue-500"
            />
          </div>
        </div>

        <div>
          <label
            htmlFor="password"
            className="block text-sm font-medium text-gray-700"
          >
            Password
          </label>
          <div className="mt-1 relative">
            <input
              id="password"
              name="password"
              type={showPassword ? "text" : "password"}
              autoComplete="current-password"
              required
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              className="appearance-none block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm placeholder-gray-400 focus:outline-none focus:ring-blue-500 focus:border-blue-500"
            />
            <button
              type="button"
              className="absolute inset-y-0 right-0 pr-3 flex items-center"
              onClick={() => setShowPassword(!showPassword)}
            >
              {showPassword ? (
                <EyeOff className="h-5 w-5 text-gray-400" />
              ) : (
                <Eye className="h-5 w-5 text-gray-400" />
              )}
            </button>
          </div>
        </div>

        <div className="flex items-center justify-between">
          <div className="flex items-center">
            <input
              id="remember-me"
              name="remember-me"
              type="checkbox"
              className="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
            />
            <label
              htmlFor="remember-me"
              className="ml-2 block text-sm text-gray-900"
            >
              Remember me
            </label>
          </div>

          <div className="text-sm">
            <Link
              to={{
                pathname: "/forgot-password",
                search: "?" + searchParams.toString(),
              }}
              className="font-medium text-blue-600 hover:text-blue-500"
            >
              Forgot your password?
            </Link>
          </div>
        </div>

        <div>
          <button
            type="submit"
            className="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
          >
            Sign in
          </button>
        </div>

        <div className="text-sm text-center">
          <span className="text-gray-600">Don't have an account?</span>{" "}
          <Link
            to={{
              pathname: "/register",
              search: "?" + searchParams.toString(),
            }}
            className="font-medium text-blue-600 hover:text-blue-500"
          >
            Sign up
          </Link>
        </div>
      </form>
    </AuthLayout>
  );
}
