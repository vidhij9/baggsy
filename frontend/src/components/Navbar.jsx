import React from "react";
import { Link } from "react-router-dom";

const Navbar = () => {
  return (
    <nav className="bg-primary text-white px-6 py-4 shadow-lg">
      <div className="container mx-auto flex justify-between items-center">
        <h1 className="text-2xl font-bold">Star Agriseeds Pvt Ltd</h1>
        <ul className="flex space-x-6">
          <li>
            <Link to="/" className="hover:text-secondary transition-all duration-300">
              Dashboard
            </Link>
          </li>
          <li>
            <Link to="/register" className="hover:text-secondary transition-all duration-300">
              Register Bag
            </Link>
          </li>
          <li>
            <Link to="/link-bags" className="hover:text-secondary transition-all duration-300">
              Link Bags
            </Link>
          </li>
          <li>
            <Link to="/search-bill" className="hover:text-secondary transition-all duration-300">
              Search Bill
            </Link>
          </li>
        </ul>
      </div>
    </nav>
  );
};

export default Navbar;