import React, { useState } from "react";
import BagRegistration from "../components/BagRegistration";
import ChildBagRegistration from "../components/ChildBagRegistration";

const Dashboard = () => {
  const [parentBag, setParentBag] = useState(null);

  const handleParentRegistered = (bag) => {
    console.log("Parent Bag Registered:", bag); // Debug log
    setParentBag(bag); // Update the state
  };

  const handleChildBagsCompleted = () => {
    console.log("Child Bags Completed"); // Debug log
    setParentBag(null); // Reset to Parent Bag Registration
    alert("All child bags registered successfully!");
  };

  return (
    <div className="min-h-screen bg-background px-6 py-10">
      <h1 className="text-3xl font-bold text-dark mb-6">Star Agriseeds</h1>
      <p className="text-gray-600 mb-6">Manage your bags and bills efficiently.</p>

      {/* Conditional Rendering */}
      {!parentBag ? (
        <BagRegistration onParentRegistered={handleParentRegistered} />
      ) : (
        <ChildBagRegistration
          parentBag={parentBag}
          onChildBagsCompleted={handleChildBagsCompleted}
        />
      )}
    </div>
  );
};

export default Dashboard;