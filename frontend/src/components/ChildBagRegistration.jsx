import React, { useState } from "react";
import { registerBag } from "../api/api";

const ChildBagRegistration = ({ parentBag, onChildBagsCompleted }) => {
  const [childBags, setChildBags] = useState([]);
  const [qrCode, setQrCode] = useState("");
  const [message, setMessage] = useState("");

  console.log("Parent Bag for Child Registration:", parentBag); // Debug log

  const handleAddChildBag = async () => {
    if (!qrCode) {
      setMessage("Child Bag QR Code is required!");
      return;
    }

    try {
      const payload = { qrCode, bagType: "Child", parentBag: parentBag.qrCode };
      const response = await registerBag(payload);

      setChildBags((prev) => [...prev, response.data.childBag]);
      setQrCode("");

      if (childBags.length + 1 === parentBag.childCount) {
        onChildBagsCompleted();
      }
    } catch (error) {
      console.error("Error during child bag registration:", error.response?.data || error.message);
      setMessage(error.response?.data?.error || "Something went wrong");
    }
  };

  return (
    <div className="bg-lightGreen p-6 rounded shadow-md">
      <h2 className="text-2xl font-bold mb-4">Register Child Bags</h2>
      <p>Parent Bag: {parentBag.qrCode}</p>
      <p>Remaining: {parentBag.childCount - childBags.length}</p>
      <input
        type="text"
        placeholder="Child Bag QR Code"
        value={qrCode}
        onChange={(e) => setQrCode(e.target.value)}
        className="w-full p-2 mb-4 border rounded"
      />
      <button
        onClick={handleAddChildBag}
        className="bg-primary text-white py-2 px-4 rounded hover:bg-dark transition-all"
      >
        Add Child Bag
      </button>
      {message && <p className="text-darkGreen mt-4">{message}</p>}
    </div>
  );
};


export default ChildBagRegistration;