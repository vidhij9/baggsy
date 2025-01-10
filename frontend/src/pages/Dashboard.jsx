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

import React, { useState } from "react";
import ParentBagRegistration from "../components/ParentBagRegistration";
import ChildBagRegistration from "../components/ChildBagRegistration";
import LinkBagToBill from "./LinkBagToBill";
import SearchBill from "./SearchBill";

const Dashboard = () => {
  const [view, setView] = useState("parent");

  return (
    <div className="min-h-screen bg-lightGray px-6 py-10">
      <h1 className="text-4xl font-bold text-primary text-center mb-8">Baggsy</h1>

      <div className="flex gap-4 justify-center mb-6">
        <button onClick={() => setView("parent")}>Register Parent</button>
        <button onClick={() => setView("link")}>Link Bag to Bill</button>
        <button onClick={() => setView("search")}>Search Bill</button>
      </div>

      {view === "parent" && <ParentBagRegistration /* props if needed */ />}
      {view === "child" && <ChildBagRegistration /* props if needed */ />}
      {view === "link" && <LinkBagToBill />}
      {view === "search" && <SearchBill />}
    </div>
  );
};

export default Dashboard;
