"use client";
import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { getNotes, createNote, getCategories, createCategory, deleteCategory } from "@/api/notes";
import NoteCard from "@/components/NoteCard";
import NoteForm from "@/components/NoteForm";
import Layout from "@/components/Layout";

export default function NotesPage() {
  const router = useRouter();
  const [notes, setNotes] = useState([]);
  const [loading, setLoading] = useState(true);
  const [creating, setCreating] = useState(false);

  const [searchTerm, setSearchTerm] = useState("");
  const [showFavorites, setShowFavorites] = useState(false);
  const [categories, setCategories] = useState([]);
  const [selectedCategory, setSelectedCategory] = useState("");
  const [newCategory, setNewCategory] = useState("");

  useEffect(() => {
    const token = localStorage.getItem("token");
    if (!token) {
      setNotes([]);
      window.location.href = "/login";
    } else {
      fetchNotes();
      fetchCategories();
    }
  }, []);

  useEffect(() => {
    fetchNotes();
  }, [searchTerm, showFavorites, selectedCategory]);

  async function fetchNotes() {
    const token = localStorage.getItem("token");
    if (!token) return setNotes([]);

    try {
      setLoading(true);
      const params = {};
      if (searchTerm.trim()) params.search = searchTerm.trim();
      if (showFavorites) params.favorite = "true";

      if (selectedCategory === "") {
        // all categories → do nothing
      } else if (selectedCategory === "none") {
        params.category_id = "none";
      } else {
        params.category_id = selectedCategory;
      }

      const data = await getNotes(params);
      setNotes(data.notes || []);
    } catch (err) {
      console.error("Failed to fetch notes:", err.response?.data?.error || err.message);
      setNotes([]);
    } finally {
      setLoading(false);
    }
  }

  async function fetchCategories() {
    try {
      const data = await getCategories();
      const cats = Array.isArray(data) ? data : data.categories || [];
      setCategories(cats);
    } catch (err) {
      console.error("Failed to fetch categories:", err.response?.data?.error || err.message);
      setCategories([]);
    }
  }

  async function handleCreate(note) {
    if (creating) return;
    setCreating(true);

    try {
      const payload = {
        title: note.title.trim(),
        body: note.body.trim(),
        category_id: note.category_id ? Number(note.category_id) : null,
        is_favorite: note.is_favorite ?? false,
        visibility: note.visibility || "private",
      };
      await createNote(payload);
      await fetchNotes();
    } catch (err) {
      console.error("Failed to create note:", err.response?.data?.error || err.message);
      alert(err.response?.data?.error || "Failed to create note. Please try again.");
    } finally {
      setCreating(false);
    }
  }

  const handleAddCategory = async () => {
    if (!newCategory.trim()) return;
    try {
      await createCategory({ name: newCategory.trim() });
      setNewCategory("");
      fetchCategories();
    } catch (err) {
      console.error("Failed to add category:", err);
      alert("Failed to add category");
    }
  };

  const handleDeleteCategory = async (id) => {
    if (!confirm("Are you sure you want to delete this category?")) return;
    try {
      await deleteCategory(id);
      fetchCategories();
      if (selectedCategory === String(id)) setSelectedCategory("");
    } catch (err) {
      console.error("Failed to delete category:", err);
      alert("Failed to delete category");
    }
  };

  return (
    <Layout
      categories={categories}
      notesCount={notes.length}
      favoritesCount={notes.filter(n => n.is_favorite).length}
      onAddCategory={async (name) => {
        try {
          await createCategory({ name });
          fetchCategories();
        } catch (err) {
          console.error("Failed to add category:", err);
        }
      }}
      onDeleteCategory={handleDeleteCategory}
      onShowAllNotes={() => { setShowFavorites(false); setSelectedCategory(""); }}
      onShowFavorites={() => { setShowFavorites(true); setSelectedCategory(""); }}
      onSelectCategory={(id) => { setShowFavorites(false); setSelectedCategory(id); }}
    >
      <NoteForm onNoteCreated={handleCreate} categories={categories} />

      {loading && <p>Loading...</p>}
      {!loading && notes.length === 0 && <p>No notes found</p>}
      <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-4">
        {notes.map((note) => (
          <NoteCard key={note.id} note={note} onUpdated={fetchNotes} />
        ))}
      </div>
    </Layout>
  );
}
