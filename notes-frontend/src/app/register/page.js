"use client";

import { useState } from "react";
import axios from "axios";
import { useRouter } from "next/navigation";

const API_BASE = "http://localhost:8080"; // backend URL

export default function RegisterPage() {
  const [username, setUsername] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [success, setSuccess] = useState("");
  const router = useRouter();

  async function handleRegister(e) {
    e.preventDefault();
    setError("");
    setSuccess("");

    try {
      await axios.post(`${API_BASE}/register`, {
        username,
        email,
        password,
      });

      setSuccess("✅ Registration successful! Redirecting to login...");
      setTimeout(() => {
        router.push("/login");
      }, 1500);
    } catch (err) {
      console.error(err);
      setError("❌ Failed to register. Try another username or email.");
    }
  }

  return (
    <div className="flex flex-col items-center justify-center min-h-screen bg-gray-100">
      <form
        onSubmit={handleRegister}
        className="bg-white p-6 rounded shadow-md w-80"
      >
        <h2 className="text-2xl font-bold mb-4 text-black">Register</h2>

        <label className="block text-sm font-medium mb-1 text-black">
          Username
        </label>
        <input
          type="text"
          value={username}
          onChange={(e) => setUsername(e.target.value)}
          required
          className="w-full border px-3 py-2 rounded mb-3 text-black placeholder-gray-500"
          placeholder="Enter username"
        />

        <label className="block text-sm font-medium mb-1 text-black">
          Email
        </label>
        <input
          type="email"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          required
          className="w-full border px-3 py-2 rounded mb-3 text-black placeholder-gray-500"
          placeholder="Enter email"
        />

        <label className="block text-sm font-medium mb-1 text-black">
          Password
        </label>
        <input
          type="password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          required
          className="w-full border px-3 py-2 rounded mb-3 text-black placeholder-gray-500"
          placeholder="Enter password"
        />

        {error && <p className="text-red-500 text-sm mb-2">{error}</p>}
        {success && <p className="text-green-500 text-sm mb-2">{success}</p>}

        <button
          type="submit"
          className="w-full bg-green-500 text-white py-2 rounded hover:bg-green-600"
        >
          Register
        </button>

        <button
          type="button"
          onClick={() => router.push("/login")}
          className="w-full mt-3 border border-gray-400 py-2 rounded hover:bg-gray-100 text-black"
        >
          Back to Login
        </button>
      </form>
    </div>
  );
}
