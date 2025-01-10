import React, { useState } from "react";
import { getBillByQRCode } from "../api/api";

const SearchBill = () => {
  const [qrCode, setQrCode] = useState("");
  const [message, setMessage] = useState("");
  const [billID, setBillID] = useState("");

  const handleSearch = async () => {
    if (!qrCode) {
      setMessage("Child Bag QR Code is required!");
      return;
    }

    try {
      // Adjust if you have a proxy or a different endpoint
      const response = await getBillByQRCode(qrCode);

      // Could be { message: "No bill linked..." } or { billId: "BILL-123" }
      setMessage(response.data.message || "");
      if (response.data.billId) {
        setBillID(response.data.billId);
      } else {
        setBillID("");
      }

    } catch (error) {
      console.error("Search Error:", error.response?.data || error.message);
      setMessage(error.response?.data?.error || "Something went wrong");
      setBillID("");
    }
  };

  return (
    <div className="bg-white p-6 rounded-lg shadow-lg max-w-md mx-auto">
      <h2 className="text-2xl font-bold text-primary mb-4">Find Bill by Bag QR Code</h2>

      <input
        type="text"
        placeholder="Enter/Scan Child Bag QR Code"
        value={qrCode}
        onChange={(e) => setQrCode(e.target.value)}
        className="w-full p-3 border rounded mb-4 focus:outline-none focus:ring-2 focus:ring-primary"
      />

      <button
        onClick={handleSearch}
        className="w-full bg-primary text-white py-3 rounded-lg hover:bg-accent transition-all"
      >
        Find Bill
      </button>

      {message && <p className="text-primary mt-4">{message}</p>}
      {billID && <p className="text-green-500 mt-2">Bill ID: {billID}</p>}
    </div>
  );
};

export default SearchBill;