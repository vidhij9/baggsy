// src/App.js
import React from "react";
import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import Dashboard from "./pages/Dashboard";
import RegisterBags from "./pages/RegisterBags";
import LinkBagToBill from "./pages/LinkBagToBill";
import SearchBill from "./pages/SearchBill";
import BottomNav from "./components/BottomNav"; // If you want the bottom nav globally

const App = () => {
  return (
    <Router>
      {/* If you want the BottomNav on every page, render it here or inside each page */}
      <BottomNav />

      <Routes>
        {/* Dashboard with the three features */}
        <Route path="/" element={<Dashboard />} />

        {/* 1) Register Bags => handles parent + child in a single flow */}
        <Route path="/register" element={<RegisterBags />} />

        {/* 2) Link Bag to Bill */}
        <Route path="/link-bag-to-bill" element={<LinkBagToBill />} />

        {/* 3) Search Bill */}
        <Route path="/search-bill" element={<SearchBill />} />
      </Routes>
    </Router>
  );
};

export default App;
