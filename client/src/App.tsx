import Navbar from "./components/Navbar";
// import TodoForm from "./components/TodoForm";
// import TodoList from "./components/TodoList";

// export const BASE_URL = import.meta.env.MODE === "development" ? "http://localhost:5000/api" : "/api";

function App() {
    return (
        <div className="flex flex-col h-screen">
            <Navbar />
            <div className="max-w-xl mx-auto w-full px-4 py-6">
                {/* <TodoForm />
                <TodoList /> */}
            </div>
        </div>
    );
}

export default App;