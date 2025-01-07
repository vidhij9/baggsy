import React, { useState } from "react";
import { registerBag } from "../api/api";

const ParentBagRegistration = ({ onParentRegistered }) => {
  const [qrCode, setQrCode] = useState("");
  const [message, setMessage] = useState("");

  const handleSubmit = async () => {
    if (!qrCode) {
      setMessage("QR Code is required!");
      return;
    }

    try {
      const payload = { qrCode, bagType: "Parent" };
      const response = await registerBag(payload);

      setMessage(response.data.message);
      onParentRegistered(response.data.bag); // Notify parent component
      setQrCode("");
    } catch (error) {
      setMessage(error.response?.data?.error || "Something went wrong");
    }
  };

  return (
    <div className="bg-white p-6 rounded-lg shadow-lg max-w-md mx-auto">
      <h2 className="text-2xl font-bold text-primary mb-4">Register Parent Bag</h2>
      <input
        type="text"
        placeholder="QR Code"
        value={qrCode}
        onChange={(e) => setQrCode(e.target.value)}
        className="w-full p-3 border rounded mb-4 focus:outline-none focus:ring-2 focus:ring-primary"
      />
      <button
        onClick={handleSubmit}
        className="w-full bg-primary text-white py-3 rounded-lg hover:bg-accent transition-all"
      >
        Submit
      </button>
      {message && <p className="text-primary mt-4">{message}</p>}
    </div>
  );
};

export default ParentBagRegistration;