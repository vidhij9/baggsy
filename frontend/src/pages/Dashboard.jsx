// import React, { useState } from "react";
// import ParentBagRegistration from "../components/ParentBagRegistration";
// import ChildBagRegistration from "../components/ChildBagRegistration";

// const Dashboard = () => {
//   const [parentBag, setParentBag] = useState(null); // Stores the registered parent bag details

//   // Triggered when a parent bag is successfully registered
//   const handleParentRegistered = (bag) => {
//     console.log("Parent Bag Registered:", bag); // Debug log
//     setParentBag(bag); // Move to the Child Bag Registration flow
//   };

//   // Triggered when all child bags are successfully linked
//   const handleChildBagsCompleted = () => {
//     console.log("All child bags linked. Returning to parent bag registration...");
//     setParentBag(null); // Reset to null, returning to the parent bag registration page
//     alert("All child bags linked successfully!");
//   };

//   return (
//     <div className="min-h-screen bg-lightGray px-6 py-10">
//       <h1 className="text-4xl font-bold text-primary text-center mb-8">
//         Baggsy
//       </h1>
//       <p className="text-lg text-gray-600 text-center mb-8">
//         Manage and Track your bags efficiently
//       </p>

//       {/* Conditional Rendering: Switch between Parent Bag Registration and Child Bag Linking */}
//       {!parentBag ? (
//         <ParentBagRegistration onParentRegistered={handleParentRegistered} />
//       ) : (
//         <ChildBagRegistration 
//          parentBag={parentBag} 
//          onChildBagsCompleted={handleChildBagsCompleted} 
//          />
//       )}
//     </div>
//   );
// };

// export default Dashboard;

// import React, { useState } from "react";
// import ParentBagRegistration from "../components/ParentBagRegistration";
// import ChildBagRegistration from "../components/ChildBagRegistration";
// import LinkBagToBill from "./LinkBagToBill";
// import SearchBill from "./SearchBill";

// const Dashboard = () => {
//   const [view, setView] = useState("parent");

//   return (
//     <div className="min-h-screen bg-lightGray px-6 py-10">
//       <h1 className="text-4xl font-bold text-primary text-center mb-8">Baggsy</h1>

//       <div className="flex gap-4 justify-center mb-6">
//         <button onClick={() => setView("parent")}>Register Parent</button>
//         <button onClick={() => setView("link")}>Link Bag to Bill</button>
//         <button onClick={() => setView("search")}>Search Bill</button>
//       </div>

//       {view === "parent" && <ParentBagRegistration /* props if needed */ />}
//       {view === "child" && <ChildBagRegistration /* props if needed */ />}
//       {view === "link" && <LinkBagToBill />}
//       {view === "search" && <SearchBill />}
//     </div>
//   );
// };

// export default Dashboard;


// src/pages/Dashboard.jsx
import React from "react";
import { Link } from "react-router-dom";
import "./dashboard.css"; // for the styling you provided

const Dashboard = () => {
  return (
    <div className="min-h-screen bg-lightGray px-6 py-10 dashboard-container">
      <h1 className="text-4xl font-bold text-primary text-center mb-8">Baggsy</h1>

      {/* Three columns (or one column on mobile) for each feature */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        {/* 1) Register Bags */}
        <div className="bg-white p-6 rounded-lg shadow text-center">
          <h2 className="text-xl font-bold mb-4">Register Bags</h2>
          <p className="mb-4">Register your parent and child bags.</p>
          <Link to="/register" className="bg-primary text-white px-4 py-2 rounded">
            Go
          </Link>
        </div>

        {/* 2) Link Bag to Bill */}
        <div className="bg-white p-6 rounded-lg shadow text-center">
          <h2 className="text-xl font-bold mb-4">Link Bag to Bill</h2>
          <p className="mb-4">Assign a bill ID to your parent bag.</p>
          <Link to="/link-bag-to-bill" className="bg-primary text-white px-4 py-2 rounded">
            Go
          </Link>
        </div>

        {/* 3) Search Bill */}
        <div className="bg-white p-6 rounded-lg shadow text-center">
          <h2 className="text-xl font-bold mb-4">Search Bill</h2>
          <p className="mb-4">Find the bill ID by scanning a bag QR code.</p>
          <Link to="/search-bill" className="bg-primary text-white px-4 py-2 rounded">
            Go
          </Link>
        </div>
      </div>
    </div>
  );
};

export default Dashboard;
