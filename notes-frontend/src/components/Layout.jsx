"use client";
import { useState } from "react";

export default function Layout({
  children,
  categories,
  notesCount,
  favoritesCount,
  onAddCategory,
  onShowAllNotes,
  onShowFavorites,
  onSelectCategory,
  onDeleteCategory,
}) {
  const [newCategory, setNewCategory] = useState("");

  const handleAdd = () => {
    if (newCategory.trim()) {
      onAddCategory(newCategory);
      setNewCategory("");
    }
  };

  return (
    <div className="flex h-screen bg-gray-100">
      {/* Sidebar */}
      <aside className="w-64 bg-gray-700 text-white flex flex-col p-4">
        <div className="flex items-center justify-between mb-8">
          <button className="text-2xl">☰</button>
        </div>

        <nav className="flex-1">
          <div
            className="flex items-center justify-between px-2 py-2 rounded-md hover:bg-gray-600 cursor-pointer"
            onClick={onShowAllNotes}
          >
            <span className="flex items-center gap-2">📝 Notes</span>
            <span className="text-sm bg-gray-600 px-2 rounded-full">{notesCount}</span>
          </div>

          <div
            className="flex items-center justify-between px-2 py-2 rounded-md hover:bg-gray-600 cursor-pointer"
            onClick={onShowFavorites}
          >
            <span className="flex items-center gap-2">⭐ Favorites</span>
            <span className="text-sm bg-gray-600 px-2 rounded-full">{favoritesCount}</span>
          </div>

          <div className="mt-6">
            <h4 className="uppercase text-xs text-gray-400 mb-2">Categories</h4>
            {categories?.map((cat) => (
              <div
                key={cat.id}
                className="flex items-center justify-between px-2 py-1 hover:bg-gray-600 rounded-md cursor-pointer"
              >
                <span onClick={() => onSelectCategory(cat.id)}>{cat.name}</span>
                <button
                  onClick={() => onDeleteCategory(cat.id)}
                  className="text-red-400 hover:text-red-200 text-xs ml-2"
                >
                  ✕
                </button>
              </div>
            ))}

            {/* Add new category */}
            <div className="mt-3 flex items-center gap-2">
              <input
                value={newCategory}
                onChange={(e) => setNewCategory(e.target.value)}
                placeholder="New category"
                className="flex-1 bg-gray-600 text-white text-sm rounded px-2 py-1 focus:outline-none"
              />
              <button
                onClick={handleAdd}
                className="text-green-400 hover:text-green-200 text-sm"
              >
                ＋
              </button>
            </div>
          </div>
        </nav>
      </aside>

      {/* Main Content */}
      <main className="flex-1 p-6 overflow-y-auto">{children}</main>
    </div>
  );
}
