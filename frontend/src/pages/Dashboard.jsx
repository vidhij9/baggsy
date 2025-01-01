import React from "react";
import BagRegistration from "../components/BagRegistration";
import BagLinking from "../components/BagLinking";
import SearchBill from "../components/SearchBill";

const Dashboard = () => {
  return (
    <div className="min-h-screen bg-lightGreen p-10">
      <h1 className="text-4xl font-bold text-darkGreen mb-10">
        Star Agriseeds Private Limited
      </h1>
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <BagRegistration />
        <BagLinking />
        <SearchBill />
      </div>
    </div>
  );
};

export default Dashboard;