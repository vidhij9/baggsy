import React, { useState } from "react";
import { searchBill } from "../api/api";

const SearchBill = () => {
  const [qrCode, setQrCode] = useState("");
  const [billId, setBillId] = useState("");

  const handleSearch = async () => {
    if (!qrCode) {
      setBillId("QR Code is required!");
      return;
    }

    try {
      const response = await searchBill(qrCode);
      setBillId(response.data.bill_id);
    } catch (error) {
      setBillId(error.response?.data?.error || "Something went wrong");
    }
  };

  return (
    <div className="bg-white p-6 rounded-lg shadow-md">
      <h2 className="text-2xl font-bold text-darkGreen mb-4">Search Bill</h2>
      <input
        type="text"
        placeholder="Enter QR Code"
        value={qrCode}
        onChange={(e) => setQrCode(e.target.value)}
        className="border rounded w-full p-2 mb-4 focus:outline-none focus:ring focus:ring-green"
      />
      <button
        onClick={handleSearch}
        className="bg-gold text-darkGreen px-4 py-2 rounded shadow hover:bg-green hover:text-white transition"
      >
        Search
      </button>
      {billId && <p className="text-green mt-4">Bill ID: {billId}</p>}
    </div>
  );
};

export default SearchBill;