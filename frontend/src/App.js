import React from "react";
import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import Dashboard from "./pages/Dashboard";
import ParentBagRegistration from "./components/ParentBagRegistration";
// import ChildBagRegistration from "./components/ChildBagRegistration";
import LinkBags from "./pages/BagLinking";
import SearchBill from "./pages/SearchBill";
import BottomNav from "./components/BottomNav";

const App = () => {
  return (
    <Router>
      <div className="flex flex-col min-h-screen">
        <div className="flex-grow">
          <Routes>
            <Route path="/" element={<Dashboard />} />
            <Route path="/register" element={<ParentBagRegistration />} />
            <Route path="/link-bags" element={<LinkBags />} />
            <Route path="/search-bill" element={<SearchBill />} />
          </Routes>
        </div>
        <BottomNav />
      </div>
    </Router>
  );
};

export default App;