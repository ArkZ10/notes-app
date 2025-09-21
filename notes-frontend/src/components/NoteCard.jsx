"use client";
import { useState } from "react";
import { updateNote, deleteNote } from "@/api/notes";
import ReactMarkdown from "react-markdown";

export default function NoteCard({ note, onUpdated }) {
  const [editing, setEditing] = useState(false);
  const [title, setTitle] = useState(note.title);
  const [body, setBody] = useState(note.body);
  const [isFavorite, setIsFavorite] = useState(note.is_favorite);
  const [loading, setLoading] = useState(false);
  const [expanded, setExpanded] = useState(false);
  const maxChars = 250;

  const handleSave = async () => {
    setLoading(true);
    try {
      await updateNote(note.id, { title, body, is_favorite: isFavorite });
      setEditing(false);
      onUpdated?.();
    } catch (err) {
      console.error("Failed to update note:", err.response?.data?.error || err.message);
      alert(err.response?.data?.error || "Failed to update note.");
    } finally {
      setLoading(false);
    }
  };

  const handleImageUpload = async (file) => {
    const formData = new FormData();
    formData.append("image", file);

    const token = localStorage.getItem("token"); // or however you store it

    const res = await fetch(`http://localhost:8080/notes/${note.id}/images`, {
    method: "POST",
    headers: {
        Authorization: `Bearer ${token}`,
    },
    body: formData,
    });

    if (!res.ok) {
      alert("Failed to upload image");
      return;
    }

    const data = await res.json();
    const imageMarkdown = `\n\n![image](${data.url})`;

    setBody((prev) => prev + imageMarkdown);
  };

  const handleCancel = () => {
    setTitle(note.title);
    setBody(note.body);
    setEditing(false);
  };

  const toggleFavorite = async () => {
    setLoading(true);
    try {
      await updateNote(note.id, { is_favorite: !isFavorite });
      setIsFavorite(!isFavorite);
      onUpdated?.();
    } catch (err) {
      console.error("Failed to toggle favorite:", err.response?.data?.error || err.message);
    } finally {
      setLoading(false);
    }
  };

  const handleDelete = async () => {
    if (!confirm("Are you sure you want to delete this note?")) return;
    setLoading(true);
    try {
      await deleteNote(note.id);
      onUpdated?.();
    } catch (err) {
      console.error("Failed to delete note:", err.response?.data?.error || err.message);
      alert(err.response?.data?.error || "Failed to delete note.");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="bg-green-100 rounded-md shadow-md p-4 min-h-[150px] flex flex-col justify-between">
      {/* Title + Favorite */}
      <div className="flex justify-between items-center">
        {editing ? (
          <input
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            className="font-semibold text-lg flex-1 bg-transparent border-b border-gray-400 focus:outline-none text-black"
          />
        ) : (
          <h2 className="font-semibold text-lg text-black">{title}</h2>
        )}

        <button
          onClick={toggleFavorite}
          disabled={loading}
          className="ml-2 text-xl text-yellow-500 hover:scale-110 transition"
        >
          {isFavorite ? "★" : "☆"}
        </button>
      </div>

      {/* Body */}
      <div className="mt-2 text-sm text-gray-800">
        {editing ? (
          <>
            <textarea
              value={body}
              onChange={(e) => setBody(e.target.value)}
              className="w-full h-24 bg-transparent border border-gray-300 rounded-md p-2 focus:outline-none text-black"
            />
            <input
              type="file"
              accept="image/*"
              className="mt-2 text-sm"
              onChange={(e) => {
                if (e.target.files[0]) handleImageUpload(e.target.files[0]);
              }}
            />
          </>
        ) : (
          <>
            <ReactMarkdown
              components={{
                img: ({ node, ...props }) => (
                  <img
                    {...props}
                    alt={props.alt || "note image"}
                    className="rounded-md mt-2 max-h-64 object-contain"
                  />
                ),
              }}
            >
              {expanded || body.length <= maxChars
                ? body
                : body.substring(0, maxChars) + "..."}
            </ReactMarkdown>
            {body.length > maxChars && (
              <button
                className="text-blue-500 text-xs mt-1"
                onClick={() => setExpanded(!expanded)}
              >
                {expanded ? "Show less" : "Read more"}
              </button>
            )}
          </>
        )}
      </div>

      {/* Actions */}
      <div className="mt-2 flex justify-end gap-2">
        {editing ? (
          <>
            <button
              className="text-green-600"
              onClick={handleSave}
              disabled={loading}
            >
              {loading ? "Saving..." : "Save"}
            </button>
            <button className="text-gray-500" onClick={handleCancel}>
              Cancel
            </button>
          </>
        ) : (
          <>
            <button
              className="text-blue-500"
              onClick={() => setEditing(true)}
            >
              Edit
            </button>
            <button
              className="text-red-500"
              onClick={handleDelete}
              disabled={loading}
            >
              Delete
            </button>
          </>
        )}
      </div>
    </div>
  );
}
