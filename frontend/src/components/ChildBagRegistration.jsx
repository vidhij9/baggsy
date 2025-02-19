import React, { useState, useEffect } from "react";
import { linkChildBag } from "../api/api";

const ChildBagRegistration = ({ parentBag, onChildBagsCompleted }) => {
  const [qrCode, setQrCode] = useState("");
  const [message, setMessage] = useState("");
  const [linkedCount, setLinkedCount] = useState(0);

  useEffect(() => {
    console.log("ChildBagRegistration loaded. ParentBag:", parentBag);
    console.log("ChildCount: ", parentBag.ChildCount || 0)
    console.log("LinkedCount: ", linkedCount)
  }, [parentBag, linkedCount]);

  const handleLinkChildBag = async () => {
    if (!qrCode) {
      setMessage("Child Bag QR Code is required!");
      return;
    }

    if (linkedCount >= parentBag.ChildCount) {
      setMessage("Child bag limit already reached!");
      onChildBagsCompleted(); // Redirect to parent bag registration
      return;
    }

    try {
      const payload = {
        parentBag: parentBag.QRCode,
        childBag: qrCode,
      };
      console.log("Payload for Link Child Bag:", payload);

      const response = await linkChildBag(payload);

      setMessage(response.data.message);

      setLinkedCount((prev) => {
        const newCount = prev + 1;
        if (newCount === parentBag.ChildCount) {
          onChildBagsCompleted(); // Redirect to parent bag registration
        }
        return newCount;
      });

      setQrCode(""); // Clear input
    } catch (error) {
      console.error("Child Bag Linking Error:", error.response?.data || error.message);
      setMessage(error.response?.data?.error || "Something went wrong");
    }
  };

  return (
    <div className="bg-white p-6 rounded-lg shadow-lg max-w-md mx-auto">
      <h2 className="text-2xl font-bold text-primary mb-4">Register Child Bags</h2>
      <p className="text-gray-600 mb-2">
        Parent Bag: {parentBag.QRCode} &nbsp;|&nbsp; 
        Allowed: {parentBag.ChildCount || 0} 
      </p>
      <p className="text-gray-600 mb-4">
        Remaining: {Math.max(0, parentBag.ChildCount - linkedCount)}
      </p>
      <input
        type="text"
        placeholder="Scan/Enter Child Bag QR Code"
        value={qrCode}
        onChange={(e) => setQrCode(e.target.value)}
        className="w-full p-3 border rounded mb-4 focus:outline-none focus:ring-2 focus:ring-primary"
      />
      <button
        onClick={handleLinkChildBag}
        className="w-full bg-primary text-white py-3 rounded-lg hover:bg-accent transition-all"
      >
        Add Child Bag
      </button>
      {message && <p className="text-primary mt-4">{message}</p>}
    </div>
  );
};

export default ChildBagRegistration;
