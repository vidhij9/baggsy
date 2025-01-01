import React, { useState } from "react";
import { registerBag } from "../api/api";

const BagRegistration = () => {
  const [qrCode, setQrCode] = useState("");
  const [bagType, setBagType] = useState("");
  const [message, setMessage] = useState("");

  const handleRegister = async () => {
    if (!qrCode || !bagType) {
      setMessage("QR Code and Bag Type are required!");
      return;
    }

    try {
      const response = await registerBag({ qr_code: qrCode, bag_type: bagType });
      setMessage(response.data.message);
      setQrCode("");
      setBagType("");
    } catch (error) {
      setMessage(error.response?.data?.error || "Something went wrong");
    }
  };

  return (
    <div className="bg-lightGreen p-6 rounded-lg shadow-md">
      <h2 className="text-2xl font-bold text-darkGreen mb-4">Register Bag</h2>
      <input
        type="text"
        placeholder="Enter QR Code"
        value={qrCode}
        onChange={(e) => setQrCode(e.target.value)}
        className="border rounded w-full p-2 mb-4 focus:outline-none focus:ring focus:ring-green"
      />
      <select
        value={bagType}
        onChange={(e) => setBagType(e.target.value)}
        className="border rounded w-full p-2 mb-4 focus:outline-none focus:ring focus:ring-green"
      >
        <option value="">Select Bag Type</option>
        <option value="Parent">Parent</option>
        <option value="Child">Child</option>
      </select>
      <button
        onClick={handleRegister}
        className="bg-green text-white px-4 py-2 rounded shadow hover:bg-darkGreen transition"
      >
        Register
      </button>
      {message && <p className="text-green mt-4">{message}</p>}
    </div>
  );
};

export default BagRegistration;