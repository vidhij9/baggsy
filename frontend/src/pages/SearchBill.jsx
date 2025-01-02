import React, { useState } from "react";
import { searchBill } from "../api/api";

const SearchBill = () => {
  const [qrCode, setQrCode] = useState("");
  const [billId, setBillId] = useState("");
  const [error, setError] = useState("");

  const handleSearch = async () => {
    if (!qrCode) {
      setError("QR Code is required!");
      return;
    }

    try {
      const response = await searchBill(qrCode);
      setBillId(response.data.bill_id);
      setError("");
    } catch (error) {
      setBillId("");
      setError(error.response?.data?.error || "Bill not found!");
    }
  };

  return (
    <div className="min-h-screen bg-background px-6 py-10">
      <h2 className="text-3xl font-bold text-dark mb-6">Search Bill</h2>
      <div className="bg-white p-6 rounded shadow-md">
        <input
          type="text"
          placeholder="Child Bag QR Code"
          value={qrCode}
          onChange={(e) => setQrCode(e.target.value)}
          className="w-full p-2 mb-4 border rounded"
        />
        <button
          onClick={handleSearch}
          className="bg-primary text-white py-2 px-4 rounded hover:bg-dark transition-all"
        >
          Search Bill
        </button>
        {billId && (
          <p className="text-green-600 mt-4">Bill ID: <strong>{billId}</strong></p>
        )}
        {error && <p className="text-red-600 mt-4">{error}</p>}
      </div>
    </div>
  );
};

export default SearchBill;