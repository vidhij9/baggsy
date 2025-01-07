import React, { useState } from "react";
import { linkChildBag } from "../api/api";

const ChildBagLinking = ({ parentBag }) => {
  const [qrCode, setQrCode] = useState("");
  const [message, setMessage] = useState("");
  const [linkedCount, setLinkedCount] = useState(0);

  const handleLinkChildBag = async () => {
    if (!qrCode) {
      setMessage("Child Bag QR Code is required!");
      return;
    }

    try {
      const payload = { parentBag: parentBag.qrCode, childBag: qrCode };
      const response = await linkChildBag(payload);

      setMessage(response.data.message);
      setLinkedCount((prev) => prev + 1);
      setQrCode("");

      if (linkedCount + 1 === parentBag.childCount) {
        alert("All child bags linked successfully!");
      }
    } catch (error) {
      console.error("Error:", error.response?.data || error.message);
      setMessage(error.response?.data?.error || "Something went wrong");
    }
  };

  return (
    <div className="bg-white p-6 rounded-lg shadow-lg max-w-md mx-auto">
      <h2 className="text-2xl font-bold text-primary mb-4">Link Child Bags</h2>
      <p className="text-gray-600 mb-4">Parent Bag: {parentBag.qrCode}</p>
      <p className="text-gray-600 mb-4">
        Remaining: {parentBag.childCount - linkedCount}
      </p>
      <input
        type="text"
        placeholder="Child Bag QR Code"
        value={qrCode}
        onChange={(e) => setQrCode(e.target.value)}
        className="w-full p-3 border rounded mb-4 focus:outline-none focus:ring-2 focus:ring-primary"
      />
      <button
        onClick={handleLinkChildBag}
        className="w-full bg-primary text-white py-3 rounded-lg hover:bg-accent transition-all"
      >
        Link
      </button>
      {message && <p className="text-primary mt-4">{message}</p>}
    </div>
  );
};

export default ChildBagLinking;