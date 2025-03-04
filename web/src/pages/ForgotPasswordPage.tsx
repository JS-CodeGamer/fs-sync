import React, { useState } from "react";
import AuthLayout from "../components/AuthLayout";

export default function ForgotPasswordForm() {
  const [email, setEmail] = useState("");
  const [submitted, setSubmitted] = useState(false);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    // Handle password reset logic here
    console.log("Reset password for:", email);
    setSubmitted(true);
  };

  if (submitted) {
    return (
      <AuthLayout
        title="Check your email"
        subtitle="We've sent you instructions to reset your password"
      >
        <div className="text-center">
          <p className="mt-2 text-sm text-gray-600">
            If you don't receive an email within a few minutes, please check
            your spam folder.
          </p>
          <button
            type="button"
            onClick={() => {
              /* Handle navigation */
            }}
            className="mt-6 text-sm font-medium text-blue-600 hover:text-blue-500"
          >
            Return to sign in
          </button>
        </div>
      </AuthLayout>
    );
  }

  return (
    <AuthLayout
      title="Reset your password"
      subtitle="Enter your email and we'll send you instructions"
    >
      <form className="space-y-6" onSubmit={handleSubmit}>
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
          <button
            type="submit"
            className="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
          >
            Send reset instructions
          </button>
        </div>

        <div className="text-sm text-center">
          <button
            type="button"
            onClick={() => {
              /* Handle navigation */
            }}
            className="font-medium text-blue-600 hover:text-blue-500"
          >
            Return to sign in
          </button>
        </div>
      </form>
    </AuthLayout>
  );
}
