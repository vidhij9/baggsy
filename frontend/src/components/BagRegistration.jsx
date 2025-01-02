import React, { useState } from "react";
import { registerBag } from "../api/api";

const BagRegistration = ({ onParentRegistered }) => {
  const [qrCode, setQrCode] = useState("");
  const [message, setMessage] = useState("");

  const decodeChildCountFromQRCode = (qrCode) => {
    // Assuming the child count is encoded in the last character of the QR code
    return parseInt(qrCode.slice(-1), 10);
  };

  const handleSubmit = async () => {
    try {
      const childCount = decodeChildCountFromQRCode(qrCode);
      const payload = { qrCode, bagType: "Parent", childCount };
      console.log("Submitting Parent Bag Data:", payload);
  
      const response = await registerBag(payload);
  
      console.log("Parent Bag Response:", response.data);
  
      // Trigger the callback with the registered bag details
      onParentRegistered(response.data.bag); // Pass the `bag` object
      setQrCode("");
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