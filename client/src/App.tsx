import React, { createContext, useState, useContext, ReactNode } from 'react';
import Navbar from "./components/Navbar";
import TodoForm from "./components/TodoForm";
import TodoList from "./components/TodoList";

export const BASE_URL = import.meta.env.MODE === "development" ? "http://localhost:5000/api" : "/api";

// Create a context for color mode
interface ColorModeContextType {
  colorMode: 'light' | 'dark';
  toggleColorMode: () => void;
}

const ColorModeContext = createContext<ColorModeContextType>({
  colorMode: 'light',
  toggleColorMode: () => {}
});

// Custom hook to use color mode
export const useColorMode = () => useContext(ColorModeContext);

// Provider component to wrap the entire app
export const ColorModeProvider: React.FC<{ children: ReactNode }> = ({ children }) => {
  const [colorMode, setColorMode] = useState<'light' | 'dark'>('light');

  const toggleColorMode = () => {
    setColorMode(colorMode === 'light' ? 'dark' : 'light');
  };

  return (
    <ColorModeContext.Provider value={{ colorMode, toggleColorMode }}>
      <div className={colorMode === 'dark' ? 'dark' : ''}>
        <div className="bg-white dark:bg-gray-900 text-black dark:text-white min-h-screen">
          {children}
        </div>
      </div>
    </ColorModeContext.Provider>
  );
};

function App() {
  return (
    <ColorModeProvider>
      <div className="flex flex-col h-screen">
        <Navbar />
        <div className="max-w-xl mx-auto w-full px-4 py-6 bg-white dark:bg-gray-900 text-black dark:text-white">
          <TodoForm />
          <TodoList /> 
        </div>
      </div>
    </ColorModeProvider>
  );
}

export default App;