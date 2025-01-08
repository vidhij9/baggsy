import React from "react";
import { Link } from "react-router-dom";
import { AiOutlineHome, AiOutlineAppstoreAdd, AiOutlineLink, AiOutlineSearch } from "react-icons/ai";

const BottomNav = () => {
  return (
    <nav className="fixed bottom-0 left-0 right-0 bg-white shadow-lg border-t py-2">
      <div className="flex justify-around items-center">
        <Link to="/" className="text-primary flex flex-col items-center hover:text-accent transition">
          <AiOutlineHome size={24} />
          <span className="text-sm">Home</span>
        </Link>
        <Link to="/register" className="text-primary flex flex-col items-center hover:text-accent transition">
          <AiOutlineAppstoreAdd size={24} />
          <span className="text-sm">Register</span>
        </Link>
        <Link to="/link-bag-to-bill" className="text-primary flex flex-col items-center hover:text-accent transition">
          <AiOutlineLink size={24} />
          <span className="text-sm">Link</span>
        </Link>
        <Link to="/search-bill" className="text-primary flex flex-col items-center hover:text-accent transition">
          <AiOutlineSearch size={24} />
          <span className="text-sm">Search</span>
        </Link>
      </div>
    </nav>
  );
};

export default BottomNav;