import { FaCheckCircle } from "react-icons/fa";
import { MdDelete } from "react-icons/md";
import { Todo } from "./TodoList";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { BASE_URL } from "../App";

const TodoItem = ({ todo }: { todo: Todo }) => {
	const queryClient = useQueryClient();

	const { mutate: updateTodo, isPending: isUpdating } = useMutation({
		mutationKey: ["updateTodo"],
		mutationFn: async () => {
			if (todo.completed) return alert("Todo is already completed");
			try {
				const res = await fetch(`${BASE_URL}/todos/${todo._id}`, {
					method: "PATCH",
				});
				const data = await res.json();
				if (!res.ok) {
					throw new Error(data.error || "Something went wrong");
				}
				return data;
			} catch (error) {
				console.log(error);
				throw error;
			}
		},
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: ["todos"] });
		},
	});

	const { mutate: deleteTodo, isPending: isDeleting } = useMutation({
		mutationKey: ["deleteTodo"],
		mutationFn: async () => {
			try {
				const res = await fetch(`${BASE_URL}/todos/${todo._id}`, {
					method: "DELETE",
				});
				const data = await res.json();
				if (!res.ok) {
					throw new Error(data.error || "Something went wrong");
				}
				return data;
			} catch (error) {
				console.log(error);
				throw error;
			}
		},
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: ["todos"] });
		},
	});

	return (
		<div className="flex gap-2 items-center">
			<div 
				className="flex-1 flex items-center border border-gray-600 p-2 rounded-lg justify-between"
			>
				<span 
					className={`
						${todo.completed ? 'text-green-200 line-through' : 'text-black-100'}
					`}
				>
					{todo.body}
				</span>
				{todo.completed && (
					<span className="ml-1 px-2 py-1 bg-green-500 text-white text-xs rounded">
						Done
					</span>
				)}
				{!todo.completed && (
					<span className="ml-1 px-2 py-1 bg-yellow-500 text-white text-xs rounded">
						In Progress
					</span>
				)}
			</div>
			<div className="flex gap-2 items-center">
				<div 
					className="text-green-500 cursor-pointer" 
					onClick={() => updateTodo()}
				>
					{!isUpdating && <FaCheckCircle size={20} />}
					{isUpdating && <svg className="animate-spin h-5 w-5" viewBox="0 0 24 24">
						<circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
						<path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
					</svg>}
				</div>
				<div 
					className="text-red-500 cursor-pointer" 
					onClick={() => deleteTodo()}
				>
					{!isDeleting && <MdDelete size={25} />}
					{isDeleting && <svg className="animate-spin h-5 w-5" viewBox="0 0 24 24">
						<circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
						<path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
					</svg>}
				</div>
			</div>
		</div>
	);
};
export default TodoItem;