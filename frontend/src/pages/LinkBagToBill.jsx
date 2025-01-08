import React, { useState } from "react";
import axios from "axios";

const LinkBagToBill = () => {
  const [parentQRCode, setParentQRCode] = useState("");
  const [billID, setBillID] = useState("");
  const [message, setMessage] = useState("");

  const handleLinkBag = async () => {
    // Basic validation on frontend
    if (!parentQRCode || !billID) {
      setMessage("Please enter both Parent Bag QR Code and Bill ID.");
      return;
    }

    try {
      const payload = {
        ParentBag: parentQRCode,
        BillID: billID
      };

      // Make POST request to our backend
      const response = await axios.post("/link-bag-to-bill", payload);
      // On success, server typically returns { message: "...", ... }
      setMessage(response.data.message);

      // Clear input fields on success
      setParentQRCode("");
      setBillID("");
    } catch (error) {
      console.error("Error linking bag to bill:", error.response?.data || error.message);
      // The backend might send { error: "some error message" }
      setMessage(error.response?.data?.error || "Something went wrong");
    }
  };

  return (
    <div className="bg-white p-6 rounded-lg shadow-lg max-w-md mx-auto">
      <h2 className="text-2xl font-bold text-primary mb-4">Link Parent Bag to Bill</h2>
      
      <input
        type="text"
        placeholder="Parent Bag QR Code"
        value={parentQRCode}
        onChange={(e) => setParentQRCode(e.target.value)}
        className="w-full p-3 border rounded mb-4 focus:outline-none focus:ring-2 focus:ring-primary"
      />

      <input
        type="text"
        placeholder="Bill ID"
        value={billID}
        onChange={(e) => setBillID(e.target.value)}
        className="w-full p-3 border rounded mb-4 focus:outline-none focus:ring-2 focus:ring-primary"
      />

      <button
        onClick={handleLinkBag}
        className="w-full bg-primary text-white py-3 rounded-lg hover:bg-accent transition-all"
      >
        Link Bag to Bill
      </button>

      {message && <p className="text-primary mt-4">{message}</p>}
    </div>
  );
};

export default LinkBagToBill;
