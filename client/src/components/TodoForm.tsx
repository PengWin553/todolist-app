import React, { useState } from 'react';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { IoMdAdd } from "react-icons/io";
import { BASE_URL } from "../App"

interface TodoMutationData {
  body: string;
}

const TodoForm: React.FC = () => {
  const [newTodo, setNewTodo] = useState<string>("");
  const queryClient = useQueryClient();

  const { mutate: createTodo, isPending: isCreating } = useMutation({
    mutationKey: ["createTodo"],
    mutationFn: async (todoData: TodoMutationData) => {
      try {
        const res = await fetch(`${BASE_URL}/todos`, {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify(todoData),
        });

        if (!res.ok) {
          const errorData = await res.json();
          throw new Error(errorData.error || "Something went wrong");
        }

        return res.json();
      } catch (error) {
        throw error instanceof Error ? error : new Error(String(error));
      }
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["todos"] });
      setNewTodo(""); // Reset input after successful creation
    },
    onError: (error: Error) => {
      console.error("Todo creation error:", error);
      alert(error.message);
    },
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (newTodo.trim()) {
      createTodo({ body: newTodo });
    }
  };

  return (
    <form onSubmit={handleSubmit} className="w-full">
      <div className="flex gap-2">
        <input
          type="text"
          value={newTodo}
          onChange={(e) => setNewTodo(e.target.value)}
          className="
            flex-grow 
            px-2 
            py-1 
            border 
            rounded-md 
            focus:outline-none 
            focus:ring-2 
            focus:ring-blue-500 
            text-black 
            dark:text-black 
            dark:bg-gray-200 
            dark:border-gray-300
          "
          placeholder="Enter a new todo"
        />
        <button
          type="submit"
          disabled={isCreating}
          className="
            mx-2 
            p-2 
            bg-blue-500 
            text-white 
            rounded-md 
            hover:bg-blue-600 
            active:scale-95 
            transition-transform 
            flex 
            items-center 
            justify-center 
            disabled:opacity-50
          "
        >
          {isCreating ? (
            <div className="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin" />
          ) : (
            <IoMdAdd size={30} />
          )}
        </button>
      </div>
    </form>
  );
};

export default TodoForm;