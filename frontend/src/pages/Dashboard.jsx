import React, { useState } from "react";
import ParentBagRegistration from "../components/ParentBagRegistration";
import ChildBagRegistration from "../components/ChildBagRegistration";

const Dashboard = () => {
  const [parentBag, setParentBag] = useState(null);

  const handleParentRegistered = (bag) => {
    setParentBag(bag); // Transition to Child Bag Registration
  };

  const handleChildBagsCompleted = () => {
    setParentBag(null); // Reset back to Parent Bag Registration
    alert("All child bags registered successfully!");
  };

  return (
    <div className="min-h-screen bg-lightGray px-6 py-10">
      <h1 className="text-4xl font-bold text-primary text-center mb-8">
        Star Agriseeds
      </h1>
      <p className="text-lg text-gray-600 text-center mb-8">
        Manage your bags and bills efficiently.
      </p>

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