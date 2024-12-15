import React, { useState } from 'react';
import { IoMoon } from "react-icons/io5";
import { LuSun } from "react-icons/lu";

export default function Navbar() {
    const [colorMode, setColorMode] = useState('light');

    const toggleColorMode = () => {
        setColorMode(colorMode === 'light' ? 'dark' : 'light');
    };

    return (
        <div className="max-w-[900px] mx-auto">
            <div className={`
                ${colorMode === 'light' ? 'bg-gray-400' : 'bg-gray-700'} 
                px-4 my-4 rounded-md
            `}>
                <div className="flex h-16 items-center justify-between">
                    {/* LEFT SIDE */}
                    <div 
                        className="hidden sm:flex justify-center items-center gap-3"
                    >
                        <img src='/react.png' alt='logo' width={50} height={50} />
                        <span className="text-[40px]">+</span>
                        <img src='/go.png' alt='logo' width={40} height={40} />
                        <span className="text-[40px]">=</span>
                        <img src='/explode.png' alt='logo' width={50} height={50} />
                    </div>
                    
                    {/* RIGHT SIDE */}
                    <div className="flex items-center gap-3">
                        <span className="text-lg font-medium">
                            Daily Tasks
                        </span>
                        
                        {/* Toggle Color Mode */}
                        <button 
                            onClick={toggleColorMode} 
                            className="focus:outline-none"
                        >
                            {colorMode === "light" ? <IoMoon /> : <LuSun size={20} />}
                        </button>
                    </div>
                </div>
            </div>
        </div>
    );
}