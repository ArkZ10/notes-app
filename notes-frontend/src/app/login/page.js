"use client";
import { useState } from "react";
import { useRouter } from "next/navigation";
import { login } from "@/api/auth";

export default function LoginPage() {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const router = useRouter();

  async function handleLogin(e) {
    e.preventDefault();
    setError("");

    try {
      await login(username, password);
      router.push("/notes"); // redirect after success
    } catch (err) {
      console.error(err);
      setError("Invalid username or password");
    }
  }

  return (
    <div className="flex flex-col items-center justify-center min-h-screen bg-gray-100">
      <form
        onSubmit={handleLogin}
        className="bg-white p-6 rounded shadow-md w-80"
      >
        <h2 className="text-2xl font-bold mb-4 text-black">Login</h2>

        <input
          type="username"
          placeholder="username"
          value={username}
          onChange={(e) => setUsername(e.target.value)}
          className="w-full border px-3 py-2 rounded mb-2 text-black"
        />

        <input
          type="password"
          placeholder="Password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          className="w-full border px-3 py-2 rounded mb-4 text-black"
        />

        <button
          type="submit"
          className="w-full bg-blue-500 text-white py-2 rounded hover:bg-blue-600"
        >
          Login
        </button>

        {/* ✅ Register button */}
        <button
          type="button"
          onClick={() => router.push("/register")}
          className="w-full mt-3 border border-gray-400 py-2 rounded hover:bg-gray-100 text-black"
        >
          Register
        </button>
      </form>
    </div>
  );
}
