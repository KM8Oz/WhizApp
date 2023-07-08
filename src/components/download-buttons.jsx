import React, { useState } from 'react';
import { team, content } from "../data";
export const DropdownButtons = () => {
  const [isLinuxOpen, setLinuxOpen] = useState(false);
  const [isWindowsOpen, setWindowsOpen] = useState(false);
  const [isMacOSOpen, setMacOSOpen] = useState(false);

  const handleLinuxClick = () => {
    setLinuxOpen(!isLinuxOpen);
    if (!isLinuxOpen) {
        setMacOSOpen(false)
        setWindowsOpen(false)
    }
  };

  const handleWindowsClick = () => {
    setWindowsOpen(!isWindowsOpen);
    if (!isWindowsOpen) {
        setMacOSOpen(false)
        setLinuxOpen(false)
    }
  };

  const handleMacOSClick = () => {
    setMacOSOpen(!isMacOSOpen);
    if (!isMacOSOpen) {
        setWindowsOpen(false)
        setLinuxOpen(false)
    }
  };

  return (
    <div className="flex  w-full justify-center">
      <div className="relative">
        <button
          type="button"
          className="bg-black hover:bg-gray-900 text-white font-mono text-sm font-bold min-w-[150px] py-2 px-4 rounded"
          onClick={handleLinuxClick}
        >
          Linux
        </button>
        {isLinuxOpen && (
          <ul className="absolute mt-2 rounded overflow-hidden animate-pulse">
            <li   className="bg-black hover:bg-gray-200 text-white font-mono text-sm font-bold min-w-[150px] py-1 px-4">
              <a href={`${content.main.link}/linux-386/`}>linux-386</a>
            </li>
            <li className="bg-black hover:bg-gray-200 text-white font-mono text-sm font-bold min-w-[150px] py-1 px-4">
              <a href={`${content.main.link}/linux-amd64/`}>linux-amd64</a>
            </li>
            <li className="bg-black hover:bg-gray-200 text-white font-mono text-sm font-bold min-w-[150px] py-1 px-4">
              <a href={`${content.main.link}/linux-arm/`}>linux-arm</a>
            </li>
            <li className="bg-black hover:bg-gray-200 text-white font-mono text-sm font-bold min-w-[150px] py-1 px-4">
              <a href={`${content.main.link}/linux-arm64/`}>linux-arm64</a>
            </li>
          </ul>
        )}
      </div>

      <div className="relative ml-2">
        <button
          type="button"
          className="bg-black hover:bg-gray-900 text-white font-mono text-sm font-bold min-w-[150px] py-2 px-4 rounded"
          onClick={handleWindowsClick}
        >
          Windows
        </button>
        {isWindowsOpen && (
            <ul className="absolute mt-2 rounded overflow-hidden animate-pulse">
            <li   className="bg-black hover:bg-gray-200 text-white font-mono text-sm font-bold min-w-[150px] py-1 px-4">
              <a href={`${content.main.link}/windows-386/`}>windows-386</a>
            </li>
            <li className="bg-black hover:bg-gray-200 text-white font-mono text-sm min-w-[150px] py-1 px-4">
              <a href={`${content.main.link}/windows-amd64/`}>windows-amd64</a>
            </li>
            <li className="bg-black hover:bg-gray-200 text-white font-mono text-sm font-bold min-w-[150px] py-1 px-4">
              <a href={`${content.main.link}/windows-arm64/`}>windows-arm64</a>
            </li>
          </ul>
        )}
      </div>

      <div className="relative ml-2">
        <button
          type="button"
          className="bg-black hover:bg-gray-900 text-white font-mono text-sm font-bold min-w-[150px] py-2 px-4 rounded"
          onClick={handleMacOSClick}
        >
          macOS
        </button>
        {isMacOSOpen && (
            <ul className="absolute mt-2 rounded overflow-hidden animate-pulse">
            <li   className="bg-black hover:bg-gray-200 text-white font-mono text-sm  font-bold min-w-[150px] py-1 px-4">
              <a href={`${content.main.link}/darwin/`}>darwin</a>
            </li>
          </ul>
        )}
      </div>
    </div>
  );
};

