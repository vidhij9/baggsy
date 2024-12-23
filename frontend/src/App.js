import React, { useState } from 'react';
import axios from 'axios';
import './App.css'; // Optional for styling

function App() {
  const [qrCode, setQrCode] = useState('');
  const [bagType, setBagType] = useState('');
  const [status, setStatus] = useState('');
  const [bags, setBags] = useState([]);
  const [billId, setBillId] = useState('');
  const [sapBillId, setSapBillId] = useState('');

  // Register a bag
  const handleRegisterBag = async () => {
    try {
      const response = await axios.post('http://localhost:8080/create-bag', {
        qr_code: qrCode,
        bag_type: bagType,
        status: status,
      });
      alert('Bag registered successfully!');
      setBags([...bags, response.data.bag]);
      setQrCode('');
      setBagType('');
      setStatus('');
    } catch (error) {
      alert('Failed to register bag. Please try again.');
      console.error(error);
    }
  };

  // Fetch all bags
  const fetchBags = async () => {
    try {
      const response = await axios.get('http://localhost:8080/bags');
      setBags(response.data.bags);
    } catch (error) {
      alert('Failed to fetch bags. Please try again.');
      console.error(error);
    }
  };

  // Link bags to a SAP Bill
  const handleLinkBagsToSAPBill = async () => {
    try {
      await axios.post('http://localhost:8080/link-bags-to-sap-bill', {
        sap_bill_id: sapBillId,
      });
      alert('Bags linked to SAP Bill successfully!');
      setSapBillId('');
    } catch (error) {
      alert('Failed to link bags. Please try again.');
      console.error(error);
    }
  };

  // Create a new bill
  const handleCreateBill = async () => {
    try {
      await axios.post('http://localhost:8080/create-bill', {
        bill_id: billId,
        description: `Bill for SAP Bill ID: ${sapBillId}`,
      });
      alert('Bill created successfully!');
      setBillId('');
      setSapBillId('');
    } catch (error) {
      alert('Failed to create bill. Please try again.');
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
            style={{
              width: '90%',
              margin: '5px 0',
              padding: '10px',
              fontSize: '14px',
              border: '1px solid #ddd',
              borderRadius: '5px',
              boxSizing: 'border-box',
            }}
          />
          <input
            type="text"
            placeholder="Bag Type (e.g., Parent, Child)"
            value={bagType}
            onChange={(e) => setBagType(e.target.value)}
            style={{
              width: '90%',
              margin: '5px 0',
              padding: '10px',
              fontSize: '14px',
              border: '1px solid #ddd',
              borderRadius: '5px',
              boxSizing: 'border-box',
            }}
          />
          <input
            type="text"
            placeholder="Status (e.g., Active)"
            value={status}
            onChange={(e) => setStatus(e.target.value)}
            style={{
              width: '90%',
              margin: '5px 0',
              padding: '10px',
              fontSize: '14px',
              border: '1px solid #ddd',
              borderRadius: '5px',
              boxSizing: 'border-box',
            }}
          />
          <button onClick={handleRegisterBag}>Register Bag</button>
        </div>

        <div className="panel">
          <h2>Bag List</h2>
          <button onClick={fetchBags}>Refresh List</button>
          <ul>
            {bags.map((bag, index) => (
              <li key={index}>
                QR Code: {bag.qr_code}, Type: {bag.bag_type}, Status: {bag.status}
              </li>
            ))}
          </ul>
        </div>

        <div className="panel">
          <h2>Link Bags to SAP Bill</h2>
          <input
            type="text"
            placeholder="SAP Bill ID"
            value={sapBillId}
            onChange={(e) => setSapBillId(e.target.value)}
          />
          <button onClick={handleLinkBagsToSAPBill}>Link to SAP Bill</button>
        </div>

        <div className="panel">
          <h2>Create Bill</h2>
          <input
            type="text"
            placeholder="Bill ID"
            value={billId}
            onChange={(e) => setBillId(e.target.value)}
          />
          <button onClick={handleCreateBill}>Create Bill</button>
        </div>
      </div>
    </div>
  );
}

export default App;
