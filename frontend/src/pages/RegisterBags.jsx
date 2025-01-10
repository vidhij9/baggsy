// src/pages/RegisterBags.jsx
import React, { useState } from "react";
import ParentBagRegistration from "../components/ParentBagRegistration";
import ChildBagRegistration from "../components/ChildBagRegistration";

const RegisterBags = () => {
  const [parentBag, setParentBag] = useState(null);

  // Triggered when a parent bag is successfully registered
  const handleParentRegistered = (bag) => {
    console.log("Parent Bag Registered:", bag);
    setParentBag(bag); // Move to the Child Bag Registration flow
  };

  // Triggered when all child bags are successfully linked
  const handleChildBagsCompleted = () => {

    // If bag.ChildCount === 0, skip child bag registration:
    if (bag.ChildCount === 0) {
      alert("This parent has 0 child capacity. No child registration needed.");
      return; // just don't set parentBag, so we remain in parent mode or redirect somewhere else
    } else {
      console.log("All child bags linked. Returning to parent bag registration...");
    }

    setParentBag(null); // Reset to null, returning to the parent bag registration page
    alert("All child bags linked successfully!");
  };

  return (
    <div className="bg-white p-6 rounded-lg shadow-lg max-w-md mx-auto">
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

export default RegisterBags;
