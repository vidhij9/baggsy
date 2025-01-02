import React from "react";
import { Link } from "react-router-dom";
import { AiOutlineHome, AiOutlineAppstoreAdd, AiOutlineSearch, AiOutlineLink } from "react-icons/ai";

const BottomNav = () => {
  return (
    <nav className="fixed bottom-0 left-0 right-0 bg-white shadow-lg border-t">
      <div className="flex justify-around items-center py-2">
        <Link to="/" className="flex flex-col items-center text-primary">
          <AiOutlineHome size={24} />
          <span className="text-sm">Home</span>
        </Link>
        <Link to="/register" className="flex flex-col items-center text-primary">
          <AiOutlineAppstoreAdd size={24} />
          <span className="text-sm">Register</span>
        </Link>
        <Link to="/link-bags" className="flex flex-col items-center text-primary">
          <AiOutlineLink size={24} />
          <span className="text-sm">Link</span>
        </Link>
        <Link to="/search-bill" className="flex flex-col items-center text-primary">
          <AiOutlineSearch size={24} />
          <span className="text-sm">Search</span>
        </Link>
      </div>
    </nav>
  );
};

export default BottomNav;