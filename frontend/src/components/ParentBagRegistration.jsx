import React, { useState } from "react";
import { registerBag } from "../api/api";  // Adjust import as needed

const ParentBagRegistration = ({ onParentRegistered }) => {
  const [qrCode, setQrCode] = useState("");
  const [message, setMessage] = useState("");

  const handleSubmit = async () => {
    if (!qrCode) {
      setMessage("Parent Bag QR Code is required!");
      return;
    }

    try {
      // We specify bagType = "Parent". Backend can parse the childCount from the QR code
      const payload = { qrCode, bagType: "Parent" };
      console.log("Registering Parent Bag Payload:", payload);
      const response = await registerBag(payload);

      console.log("Parent Bag Registration Response:", response.data);
      setMessage(response.data.message);

      // The backend typically returns { message: "...", bag: { QRCode, childCount, ...} }
      if (response.data.bag) {
        onParentRegistered(response.data.bag);
      }
      setQrCode("");
    } catch (error) {
      console.error("Parent Bag Registration Error:", error.response?.data || error.message);
      setMessage(error.response?.data?.error || "Something went wrong");
    }
  };

  return (
    <div className="bg-white p-6 rounded-lg shadow-lg max-w-md mx-auto">
      <h2 className="text-2xl font-bold text-primary mb-4">Register Parent Bag</h2>
      <input
        type="text"
        placeholder="Enter/Scan Parent Bag QR Code"
        value={qrCode}
        onChange={(e) => setQrCode(e.target.value)}
        className="w-full p-3 border rounded mb-4 focus:outline-none focus:ring-2 focus:ring-primary"
      />
      <button
        onClick={handleSubmit}
        className="w-full bg-primary text-white py-3 rounded-lg hover:bg-accent transition-all"
      >
        Register Parent Bag
      </button>
      {message && <p className="text-primary mt-4">{message}</p>}
    </div>
  );
};

export default ParentBagRegistration;
