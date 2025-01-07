import React, { useState } from "react";
import ParentBagRegistration from "../components/ParentBagRegistration";
import ChildBagRegistration from "../components/ChildBagRegistration";

const Dashboard = () => {
  const [parentBag, setParentBag] = useState(null); // Stores the registered parent bag details

  // Triggered when a parent bag is successfully registered
  const handleParentRegistered = (bag) => {
    console.log("Parent Bag Registered:", bag); // Debug log
    setParentBag(bag); // Move to the Child Bag Registration flow
  };

  // Triggered when all child bags are successfully linked
  const handleChildBagsCompleted = () => {
    console.log("All Child Bags Linked!"); // Debug log
    setParentBag(null); // Reset to allow new parent bag registration
    alert("All child bags linked successfully!");
  };

  return (
    <div className="min-h-screen bg-lightGray px-6 py-10">
      <h1 className="text-4xl font-bold text-primary text-center mb-8">
        Star Agriseeds
      </h1>
      <p className="text-lg text-gray-600 text-center mb-8">
        Manage your bags and bills efficiently.
      </p>

      {/* Conditional Rendering: Switch between Parent Bag Registration and Child Bag Linking */}
      {!parentBag ? (
        <ParentBagRegistration onParentRegistered={handleParentRegistered} />
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