"use client";
import { useEffect, useState } from "react";
import { useParams, useRouter } from "next/navigation";
import { getNoteById, updateNote, deleteNote } from "@/api/notes";

export default function NoteDetailPage() {
  const { id } = useParams();
  const router = useRouter();
  const [note, setNote] = useState(null);
  const [isEditing, setIsEditing] = useState(false);
  const [title, setTitle] = useState("");
  const [body, setBody] = useState("");

  useEffect(() => {
    async function fetchNote() {
      try {
        const data = await getNoteById(id);
        setNote(data);
        setTitle(data.title);
        setBody(data.body);
      } catch (err) {
        console.error(err);
      }
    }

    fetchNote();
  }, [id]);

  async function handleUpdate() {
  try {
    await updateNote(id, { title, body });
    setNote({ ...note, title, body }); // update local state instantly
    setIsEditing(false);
  } catch (err) {
    console.error("Failed to update note:", err.message);
  }
}

  async function handleDelete() {
    try {
      await deleteNote(id);
      router.push("/notes"); // back to notes list
    } catch (err) {
      console.error(err);
    }
  }

  if (!note) return <p>Loading...</p>;

  return (
    <div>
      {isEditing ? (
        <form onSubmit={handleUpdate}>
          <h1>Edit Note</h1>
          <input
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            placeholder="Title"
          />
          <textarea
            value={body}
            onChange={(e) => setBody(e.target.value)}
            placeholder="Body"
          />
          <button type="submit">Save</button>
          <button type="button" onClick={() => setIsEditing(false)}>
            Cancel
          </button>
        </form>
      ) : (
        <>
          <h1>{note.title}</h1>
          <p>{note.body}</p>
          <button onClick={() => setIsEditing(true)}>Edit</button>
          <button onClick={handleDelete} style={{ color: "red" }}>
            Delete
          </button>
        </>
      )}
    </div>
  );
}