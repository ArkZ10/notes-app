"use client";
import { useState } from "react";

export default function NoteForm({ onNoteCreated, categories = [] }) {
  const [title, setTitle] = useState("");
  const [body, setBody] = useState("");
  const [categoryId, setCategoryId] = useState("");

  const handleSubmit = (e) => {
    e.preventDefault();
    if (!title.trim()) return alert("Title is required");
    onNoteCreated({ title, body, category_id: categoryId || null });
    setTitle("");
    setBody("");
    setCategoryId("");
  };

//   const handleImageUpload = async (file) => {
//     // Since note doesn’t exist yet, we just embed the image in the body
//     const formData = new FormData();
//     formData.append("image", file);

//     const token = localStorage.getItem("token");

//     const res = await fetch(`http://localhost:8080/notes/${note.id}/images`, {
//     method: "POST",
//     headers: {
//         Authorization: `Bearer ${token}`,
//     },
//     body: formData,
//     });

//     if (!res.ok) {
//       alert("Failed to upload image");
//       return;
//     }

//     const data = await res.json();
//     setBody((prev) => prev + `\n\n![image](${data.url})`);
//   };

  return (
    <form onSubmit={handleSubmit} style={{ marginBottom: "1rem" }}>
      <input
        type="text"
        placeholder="Title"
        value={title}
        onChange={(e) => setTitle(e.target.value)}
        className="w-full border rounded px-2 py-1 text-black placeholder-gray-500 mb-2"
      />
      <textarea
        placeholder="Body"
        value={body}
        onChange={(e) => setBody(e.target.value)}
        rows={3}
        className="w-full border rounded px-2 py-1 text-black placeholder-gray-500 mb-2"
      />
      {/* <input
        type="file"
        accept="image/*"
        className="w-full border rounded px-2 py-1 text-black placeholder-gray-500 mb-2"
        onChange={(e) => {
          if (e.target.files[0]) handleImageUpload(e.target.files[0]);
        }}
      /> */}

      <select
        value={categoryId}
        onChange={(e) => setCategoryId(e.target.value)}
        className="w-full border rounded px-2 py-1 text-black placeholder-gray-500 mb-2"
      >
        <option value="">No category</option>
        {categories.map((cat) => (
          <option key={cat.id} value={cat.id}>
            {cat.name}
          </option>
        ))}
      </select>

      <button
        type="submit"
        className="w-full border rounded px-2 py-1 bg-blue-500 text-white"
      >
        Create Note
      </button>
    </form>
  );
}
