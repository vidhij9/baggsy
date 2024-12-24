import React, { useState } from 'react';
import axios from 'axios';
import './App.css';

function App() {
  const [qrCode, setQrCode] = useState('');
  const [bagType, setBagType] = useState('Parent');
  const [parentBagQrCode, setParentBagQrCode] = useState('');
  const [childBagQrCode, setChildBagQrCode] = useState('');
  const [billId, setBillId] = useState('');
  const [bagSearchQrCode, setBagSearchQrCode] = useState('');
  const [billSearchId, setBillSearchId] = useState('');
  const [newBillId, setNewBillId] = useState('');

  // Register a bag
  const handleRegisterBag = async () => {
    try {
      const response = await axios.post('http://localhost:8080/create-bag', {
        qr_code: qrCode,
        bag_type: bagType,
      });
      alert('Bag registered successfully!');
      setQrCode('');
      setBagType('Parent');
    } catch (error) {
      alert('Failed to register bag. Please try again.');
      console.error(error);
    }
  };

  // Link parent and child bags
  const handleLinkBags = async () => {
    try {
      await axios.post('http://localhost:8080/link-bags', {
        parent_bag_qr_code: parentBagQrCode,
        child_bag_qr_code: childBagQrCode,
      });
      alert('Bags linked successfully!');
      setParentBagQrCode('');
      setChildBagQrCode('');
    } catch (error) {
      alert('Failed to link bags. Please try again.');
      console.error(error);
    }
  };

  // Link parent bag to a bill ID
  const handleLinkBagToBill = async () => {
    try {
      await axios.post('http://localhost:8080/link-bag-to-bill', {
        parent_bag_qr_code: parentBagQrCode,
        bill_id: billId,
      });
      alert('Bag linked to bill successfully!');
      setParentBagQrCode('');
      setBillId('');
    } catch (error) {
      alert('Failed to link bag to bill. Please try again.');
      console.error(error);
    }
  };

  // Search for a bag to get bill ID
  const handleSearchBillByBag = async () => {
    try {
      const response = await axios.get(`http://localhost:8080/search-bill-by-bag?qr_code=${bagSearchQrCode}`);
      alert(`Bill ID: ${response.data.bill_id}`);
      setBagSearchQrCode('');
    } catch (error) {
      alert('Failed to find bill ID. Please try again.');
      console.error(error);
    }
  };

  // Search for a bill ID to get linked bags
  const handleSearchBagsByBill = async () => {
    try {
      const response = await axios.get(`http://localhost:8080/search-bags-by-bill?bill_id=${billSearchId}`);
      alert(`Linked Bags: ${response.data.bags.join(', ')}`);
      setBillSearchId('');
    } catch (error) {
      alert('Failed to find linked bags. Please try again.');
      console.error(error);
    }
  };

  // Edit bill ID for a parent bag
  const handleEditBillId = async () => {
    try {
      await axios.put('http://localhost:8080/edit-bill-id', {
        parent_bag_qr_code: parentBagQrCode,
        new_bill_id: newBillId,
      });
      alert('Bill ID updated successfully!');
      setParentBagQrCode('');
      setNewBillId('');
    } catch (error) {
      alert('Failed to update bill ID. Please try again.');
      console.error(error);
    }
  };

  return (
    <div style={{ padding: '20px', fontFamily: 'Arial, sans-serif' }}>
      <h1 style={{ textAlign: 'center', marginBottom: '20px', color: '#2c3e50' }}>Baggsy Dashboard</h1>

      <div className="dashboard">
        <div className="panel">
          <h2>Register Bag</h2>
          <input
            type="text"
            placeholder="QR Code"
            value={qrCode}
            onChange={(e) => setQrCode(e.target.value)}
          />
          <div>
            <label>
              <input
                type="radio"
                name="bagType"
                value="Parent"
                checked={bagType === 'Parent'}
                onChange={(e) => setBagType(e.target.value)}
              />
              Parent
            </label>
            <label style={{ marginLeft: '10px' }}>
              <input
                type="radio"
                name="bagType"
                value="Child"
                checked={bagType === 'Child'}
                onChange={(e) => setBagType(e.target.value)}
              />
              Child
            </label>
          </div>
          <button onClick={handleRegisterBag}>Register Bag</button>
        </div>

        <div className="panel">
          <h2>Link Bags</h2>
          <input
            type="text"
            placeholder="Parent Bag QR Code"
            value={parentBagQrCode}
            onChange={(e) => setParentBagQrCode(e.target.value)}
          />
          <input
            type="text"
            placeholder="Child Bag QR Code"
            value={childBagQrCode}
            onChange={(e) => setChildBagQrCode(e.target.value)}
          />
          <button onClick={handleLinkBags}>Link Bags</button>
        </div>

        <div className="panel">
          <h2>Link Parent Bag to Bill ID</h2>
          <input
            type="text"
            placeholder="Parent Bag QR Code"
            value={parentBagQrCode}
            onChange={(e) => setParentBagQrCode(e.target.value)}
          />
          <input
            type="text"
            placeholder="Bill ID"
            value={billId}
            onChange={(e) => setBillId(e.target.value)}
          />
          <button onClick={handleLinkBagToBill}>Link to Bill</button>
        </div>

        <div className="panel">
          <h2>Search Bill ID by Bag</h2>
          <input
            type="text"
            placeholder="Enter Bag QR Code"
            value={bagSearchQrCode}
            onChange={(e) => setBagSearchQrCode(e.target.value)}
          />
          <button onClick={handleSearchBillByBag}>Search</button>
        </div>

        <div className="panel">
          <h2>Search Bags by Bill ID</h2>
          <input
            type="text"
            placeholder="Enter Bill ID"
            value={billSearchId}
            onChange={(e) => setBillSearchId(e.target.value)}
          />
          <button onClick={handleSearchBagsByBill}>Search</button>
        </div>

        <div className="panel">
          <h2>Edit Bill ID for Parent Bag</h2>
          <input
            type="text"
            placeholder="Parent Bag QR Code"
            value={parentBagQrCode}
            onChange={(e) => setParentBagQrCode(e.target.value)}
          />
          <input
            type="text"
            placeholder="New Bill ID"
            value={newBillId}
            onChange={(e) => setNewBillId(e.target.value)}
          />
          <button onClick={handleEditBillId}>Save</button>
        </div>
      </div>
    </div>
  );
}

export default App;
