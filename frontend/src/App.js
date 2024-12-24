import React, { useState } from "react";
import axios from "axios";
import "./App.css";

function App() {
  // Unique state variables for each input field
  const [registerQrCode, setRegisterQrCode] = useState("");
  const [bagType, setBagType] = useState("");
  const [parentBag, setParentBag] = useState("");
  const [childBag, setChildBag] = useState("");
  const [billID, setBillID] = useState("");
  const [searchQrCode, setSearchQrCode] = useState("");
  const [searchResult, setSearchResult] = useState(null);
  const [error, setError] = useState("");

  // API calls
  const handleRegisterBag = async () => {
    try {
      await axios.post("/register-bag", { qr_code: registerQrCode, bag_type: bagType });
      alert("Bag registered successfully");
      setRegisterQrCode("");
      setBagType("");
      setError(""); // Clear any previous error
    } catch (err) {
      setError(err.response?.data?.error || "Failed to register bag");
    }
  };

  const handleLinkBags = async () => {
    try {
      await axios.post("/link-bags", { parent_bag: parentBag, child_bag: childBag });
      alert("Bags linked successfully");
      setParentBag("");
      setChildBag("");
      setError("");
    } catch (err) {
      setError(err.response?.data?.error || "Failed to link bags");
    }
  };

  const handleLinkBagToBill = async () => {
    try {
      await axios.post("/link-bag-to-bill", { parent_bag: parentBag, bill_id: billID });
      alert("Bag linked to bill successfully");
      setParentBag("");
      setBillID("");
      setError("");
    } catch (err) {
      setError(err.response?.data?.error || "Failed to link bag to bill");
    }
  };

  const handleSearchBillByBag = async () => {
    try {
      const response = await axios.get(`/search-bill-by-bag?qr_code=${searchQrCode}`);
      setSearchResult(response.data);
      setError("");
    } catch (err) {
      setError(err.response?.data?.error || "Failed to find bill by bag");
    }
  };

  const handleSearchBagsByBill = async () => {
    try {
      const response = await axios.get(`/search-bags-by-bill?bill_id=${billID}`);
      setSearchResult(response.data);
      setError("");
    } catch (err) {
      setError(err.response?.data?.error || "Failed to find bags by bill");
    }
  };

  return (
    <div className="app-container">
      <header className="app-header">
        <h1>Baggsy</h1>
        <p className="sub-header">A product by Star Agriseeds</p>
      </header>
      {error && <p style={{ color: "red" }}>{error}</p>}

      <div className="panel-container">
        <div className="panel">
          <h2>Register Bag</h2>
          <input
            type="text"
            placeholder="QR Code"
            value={registerQrCode}
            onChange={(e) => setRegisterQrCode(e.target.value)}
          />
          <select value={bagType} onChange={(e) => setBagType(e.target.value)}>
            <option value="">Select Bag Type</option>
            <option value="Parent">Parent</option>
            <option value="Child">Child</option>
          </select>
          <button onClick={handleRegisterBag}>Register Bag</button>
        </div>

        <div className="panel">
          <h2>Link Bags</h2>
          <input
            type="text"
            placeholder="Parent Bag QR Code"
            value={parentBag}
            onChange={(e) => setParentBag(e.target.value)}
          />
          <input
            type="text"
            placeholder="Child Bag QR Code"
            value={childBag}
            onChange={(e) => setChildBag(e.target.value)}
          />
          <button onClick={handleLinkBags}>Link Bags</button>
        </div>

        <div className="panel">
          <h2>Link Bag to Bill</h2>
          <input
            type="text"
            placeholder="Parent Bag QR Code"
            value={parentBag}
            onChange={(e) => setParentBag(e.target.value)}
          />
          <input
            type="text"
            placeholder="Bill ID"
            value={billID}
            onChange={(e) => setBillID(e.target.value)}
          />
          <button onClick={handleLinkBagToBill}>Link Bag to Bill</button>
        </div>

        <div className="panel">
          <h2>Search Bill by Bag</h2>
          <input
            type="text"
            placeholder="QR Code"
            value={searchQrCode}
            onChange={(e) => setSearchQrCode(e.target.value)}
          />
          <button onClick={handleSearchBillByBag}>Search</button>
        </div>

        <div className="panel">
          <h2>Search Bags by Bill</h2>
          <input
            type="text"
            placeholder="Bill ID"
            value={billID}
            onChange={(e) => setBillID(e.target.value)}
          />
          <button onClick={handleSearchBagsByBill}>Search</button>
        </div>
      </div>

      <div className="results">
        <h2>Results</h2>
        {searchResult && <pre>{JSON.stringify(searchResult, null, 2)}</pre>}
      </div>
    </div>
  );
}

export default App;
