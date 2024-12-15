import { useState } from "react";
import { useQuery } from "@tanstack/react-query";
import TodoItem from "./TodoItem";
import { BASE_URL } from "../App";

export type Todo = {
	_id: number;
	body: string;
	completed: boolean;
};

const TodoList = () => {
	const { data: todos, isLoading } = useQuery<Todo[]>({
		queryKey: ["todos"],
		queryFn: async () => {
			try {
				const res = await fetch(BASE_URL + "/todos");
				if (!res.ok) {
					// Throw an error if the response is not OK
					throw new Error("Failed to fetch todos");
				}
				const data = await res.json();
				return data || []; // Ensure we return an array
			} catch (error) {
				console.error("Error fetching todos:", error);
				return []; // Return an empty array on error
			}
		},
		// Add retry option to prevent infinite retries
		retry: 1
	});

	return (
		<>
			<h1 className="text-4xl uppercase font-bold text-center my-2 bg-gradient-to-l from-blue-500 to-cyan-500 bg-clip-text text-transparent">
				Today's Tasks
			</h1>
			
			{isLoading && (
				<div className="flex justify-center my-4">
					<div className="w-16 h-16 border-4 border-t-4 border-blue-500 border-t-transparent rounded-full animate-spin"></div>
				</div>
			)}
			
			{!isLoading && todos?.length === 0 && (
				<div className="flex flex-col items-center gap-3">
					<p className="text-xl text-center text-gray-500">
						All tasks completed! 🤞
					</p>
					<img src='/go.png' alt='Go logo' width={70} height={70} />
				</div>
			)}
			
			<div className="space-y-3">
				{todos?.map((todo) => (
					<TodoItem key={todo._id} todo={todo} />
				))}
			</div>
		</>
	);
};

export default TodoList;

