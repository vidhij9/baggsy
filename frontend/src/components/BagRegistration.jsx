import React, { useState } from "react";
import { registerBag } from "../api/api";

const BagRegistration = ({ onParentRegistered }) => {
  const [qrCode, setQrCode] = useState("");
  const [childCount, setChildCount] = useState(0);
  const [message, setMessage] = useState("");

  const handleSubmit = async () => {
    if (!qrCode || childCount <= 0) {
      setMessage("QR Code and positive child count are required!");
      return;
    }

    try {
      const payload = { qrCode, bagType: "Parent", childCount };
      const response = await registerBag(payload);

      // Success Message
      setMessage(response.data.message);

      // Trigger Parent Registration Callback
      onParentRegistered(response.data.parentBag);

      // Clear inputs
      setQrCode("");
      setChildCount(0);
    } catch (error) {
      console.error("Error during registration:", error.response?.data || error.message);
      setMessage(error.response?.data?.error || "Something went wrong");
    }
  };

  return (
    <div className="bg-lightGreen p-6 rounded shadow-md">
      <h2 className="text-2xl font-bold mb-4">Register Parent Bag</h2>
      <input
        type="text"
        placeholder="QR Code"
        value={qrCode}
        onChange={(e) => setQrCode(e.target.value)}
        className="w-full p-2 mb-4 border rounded"
      />
      <input
        type="number"
        placeholder="Number of Child Bags"
        value={childCount}
        onChange={(e) => setChildCount(Number(e.target.value))}
        className="w-full p-2 mb-4 border rounded"
      />
      <button
        onClick={handleSubmit}
        className="bg-primary text-white py-2 px-4 rounded hover:bg-dark transition-all"
      >
        Submit
      </button>
      {message && <p className="text-darkGreen mt-4">{message}</p>}
    </div>
  );
};

export default BagRegistration;