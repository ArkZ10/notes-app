import axios from "axios";

const API_BASE = "http://localhost:8080"; // backend URL

// Helper to get token
function getToken() {
  return localStorage.getItem("token");
}

// Add interceptor to redirect to login on 401
const api = axios.create({
  baseURL: API_BASE,
});

api.interceptors.response.use(
  (res) => res,
  (err) => {
    if (err.response && err.response.status === 401) {
      localStorage.removeItem("token"); // clear old token
      window.location.href = "/login"; // force redirect
    }
    return Promise.reject(err);
  }
);

export const getNotes = async (params = {}) => {
  const token = getToken();

  const queryParams = new URLSearchParams();
  if (params.search) queryParams.append("search", params.search);
  if (params.favorite) queryParams.append("favorite", params.favorite);

  // Add category_id filter
  if (params.category_id) {
    queryParams.append("category_id", params.category_id);
  }

  const res = await api.get(`/notes?${queryParams.toString()}`, {
    headers: { Authorization: `Bearer ${token}` },
  });
  return res.data;
};

export const getNoteById = async (id) => {
  const token = getToken();
  const res = await api.get(`/notes/${id}`, {
    headers: { Authorization: `Bearer ${token}` },
  });
  return res.data;
};

export const createNote = async (noteData) => {
  const token = getToken();

  // Validate title
  if (!noteData.title || !noteData.title.trim()) {
    throw new Error("Title is required");
  }

  const payload = {
    title: noteData.title.trim(),
    body: noteData.body?.trim() || "",
    category_id: noteData.category_id ?? null,
    is_favorite: noteData.is_favorite ?? false,
    visibility: noteData.visibility || "private",
  };

  const res = await api.post("/notes", payload, {
    headers: {
      Authorization: `Bearer ${token}`,
      "Content-Type": "application/json",
    },
  });

  return res.data;
};

export const updateNote = async (id, noteData) => {
  const token = getToken();

  const payload = {};
  if (noteData.title !== undefined) payload.title = noteData.title;
  if (noteData.body !== undefined) payload.body = noteData.body;
  if ("category_id" in noteData) payload.category_id = noteData.category_id ?? null;
  if ("is_favorite" in noteData) payload.is_favorite = noteData.is_favorite ?? false;
  if (noteData.visibility !== undefined) payload.visibility = noteData.visibility;

  const res = await api.patch(`/notes/${id}`, payload, {
    headers: { Authorization:  `Bearer ${token}`, "Content-Type": "application/json" },
  });
  return res.data;
};

export const deleteNote = async (id) => {
  const token = getToken();
  const res = await api.delete(`/notes/${id}`, {
    headers: { Authorization: `Bearer ${token}` },
  });
  return res.data;
};

export const getCategories = async () => {
  const token = getToken();
  const res = await api.get("/categories", {
    headers: { Authorization: `Bearer ${token}` },
  });
  return res.data;
};

// Create a new category
export const createCategory = async ({ name }) => {
  const token = getToken();
  if (!name || !name.trim()) throw new Error("Category name is required");

  const res = await api.post(
    "/categories",
    { name: name.trim() },
    {
      headers: {
        Authorization: `Bearer ${token}`,
        "Content-Type": "application/json",
      },
    }
  );
  return res.data;
};

export const deleteCategory = async (id) => {
  const token = getToken();
  const res = await api.delete(`/categories/${id}`, {
    headers: { Authorization: `Bearer ${token}` },
  });
  return res.data;
};