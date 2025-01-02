import React, { useState } from "react";
import { linkBags } from "../api/api";

const LinkBags = () => {
  const [parentBag, setParentBag] = useState("");
  const [childBag, setChildBag] = useState("");
  const [message, setMessage] = useState("");

  const handleLinkBags = async () => {
    if (!parentBag || !childBag) {
      setMessage("Parent Bag and Child Bag QR Codes are required!");
      return;
    }

    try {
      const payload = { parent_bag: parentBag, child_bag: childBag };
      const response = await linkBags(payload);
      setMessage(response.data.message);
      setParentBag("");
      setChildBag("");
    } catch (error) {
      setMessage(error.response?.data?.error || "Something went wrong");
    }
  };

  return (
    <div className="min-h-screen bg-background px-6 py-10">
      <h2 className="text-3xl font-bold text-dark mb-6">Link Bags</h2>
      <div className="bg-white p-6 rounded shadow-md">
        <input
          type="text"
          placeholder="Parent Bag QR Code"
          value={parentBag}
          onChange={(e) => setParentBag(e.target.value)}
          className="w-full p-2 mb-4 border rounded"
        />
        <input
          type="text"
          placeholder="Child Bag QR Code"
          value={childBag}
          onChange={(e) => setChildBag(e.target.value)}
          className="w-full p-2 mb-4 border rounded"
        />
        <button
          onClick={handleLinkBags}
          className="bg-primary text-white py-2 px-4 rounded hover:bg-dark transition-all"
        >
          Link Bags
        </button>
        {message && <p className="text-green-600 mt-4">{message}</p>}
      </div>
    </div>
  );
};

export default LinkBags;