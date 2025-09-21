import axios from "axios";

const API_BASE = "http://localhost:8080"; // backend

export const login = async (username, password) => {
  const res = await axios.post(`${API_BASE}/login`, { username, password });
  // Save JWT
  localStorage.setItem("token", res.data.token);
  return res.data;
};

export const register = async (username, password) => {
  const res = await axios.post(`${API_BASE}/register`, { username, password });
  return res.data;
};

export const logout = () => {
  localStorage.removeItem("token");
};
