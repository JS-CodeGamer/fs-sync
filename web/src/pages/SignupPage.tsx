import React, { useState } from "react";
import { Eye, EyeOff } from "lucide-react";
import { Link, useNavigate, useSearchParams } from "react-router";

import { useAuth } from "../context/AuthProvider";
import AuthLayout from "../components/AuthLayout";
import { Tokens } from "../models/user";
import backend from "../helpers/backend";
import { toast } from "react-toastify";
import { handleError } from "../helpers/errorHandler";

export default function SignupPage() {
  const [searchParams, _] = useSearchParams();
  const { setTokens } = useAuth();
  const navigate = useNavigate();

  const [showPassword, setShowPassword] = useState(false);
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [username, setName] = useState("");

  async function registerUser({
    email,
    username,
    password,
  }: {
    email: string;
    username: string;
    password: string;
  }) {
    try {
      const res = await backend.post<Tokens>("/register", {
        email,
        username,
        password,
      });
      if (res) {
        setTokens(res?.data);
        toast.success("Registration Success!");
        navigate(decodeURIComponent(searchParams.get("return") ?? "/"));
      } else {
        toast.error("Registration Failed!!!");
      }
    } catch (error) {
      handleError(error);
    }
  }

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    registerUser({ email, username, password });
  };

  return (
    <AuthLayout
      title="Create an account"
      subtitle="Get started with your free account"
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
              type="text"
              autoComplete="username"
              required
              value={username}
              onChange={(e) => setName(e.target.value)}
              className="appearance-none block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm placeholder-gray-400 focus:outline-none focus:ring-blue-500 focus:border-blue-500"
            />
          </div>
        </div>

        <div>
          <label
            htmlFor="email"
            className="block text-sm font-medium text-gray-700"
          >
            Email address
          </label>
          <div className="mt-1">
            <input
              id="email"
              name="email"
              type="email"
              autoComplete="email"
              required
              value={email}
              onChange={(e) => setEmail(e.target.value)}
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
              autoComplete="new-password"
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

        <div>
          <button
            type="submit"
            className="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
          >
            Create account
          </button>
        </div>

        <div className="text-sm text-center">
          <span className="text-gray-600">Already have an account?</span>{" "}
          <Link
            to={{
              pathname: "/login",
              search: "?" + searchParams.toString(),
            }}
            className="font-medium text-blue-600 hover:text-blue-500"
          >
            Sign in
          </Link>
        </div>
      </form>
    </AuthLayout>
  );
}
