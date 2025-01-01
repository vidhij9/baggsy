import React, { useState } from "react";
import { linkBags } from "../api/api";

const BagLinking = () => {
  const [parentBag, setParentBag] = useState("");
  const [childBag, setChildBag] = useState("");
  const [message, setMessage] = useState("");

  const handleLink = async () => {
    if (!parentBag || !childBag) {
      setMessage("Parent Bag and Child Bag QR Codes are required!");
      return;
    }

    try {
      const response = await linkBags({ parent_bag: parentBag, child_bag: childBag });
      setMessage(response.data.message);
      setParentBag("");
      setChildBag("");
    } catch (error) {
      setMessage(error.response?.data?.error || "Something went wrong");
    }
  };

  return (
    <div className="bg-gold p-6 rounded-lg shadow-md">
      <h2 className="text-2xl font-bold text-darkGreen mb-4">Link Bags</h2>
      <input
        type="text"
        placeholder="Enter Parent Bag QR Code"
        value={parentBag}
        onChange={(e) => setParentBag(e.target.value)}
        className="border rounded w-full p-2 mb-4 focus:outline-none focus:ring focus:ring-green"
      />
      <input
        type="text"
        placeholder="Enter Child Bag QR Code"
        value={childBag}
        onChange={(e) => setChildBag(e.target.value)}
        className="border rounded w-full p-2 mb-4 focus:outline-none focus:ring focus:ring-green"
      />
      <button
        onClick={handleLink}
        className="bg-darkGreen text-white px-4 py-2 rounded shadow hover:bg-green transition"
      >
        Link Bags
      </button>
      {message && <p className="text-green mt-4">{message}</p>}
    </div>
  );
};

export default BagLinking;