import React, { useState } from 'react';
import axios from 'axios';

const App = () => {
  const [sapBillID, setSapBillID] = useState('');
  const [description, setDescription] = useState('');
  const [parentBags, setParentBags] = useState('');

  const createBill = async () => {
    try {
      const res = await axios.post('/create-bill', {
        sap_bill_id: sapBillID,
        description,
        parent_bags: parentBags.split(','),
      });
      alert(res.data.message);
      setSapBillID('');
      setDescription('');
      setParentBags('');
    } catch (err) {
      console.error(err);
    }
  };

  return (
    <div>
      <h1>Create a Bill</h1>
      <input
        value={sapBillID}
        onChange={(e) => setSapBillID(e.target.value)}
        placeholder="SAP Bill ID"
      />
      <input
        value={description}
        onChange={(e) => setDescription(e.target.value)}
        placeholder="Description"
      />
      <textarea
        value={parentBags}
        onChange={(e) => setParentBags(e.target.value)}
        placeholder="Comma-separated Parent Bag IDs"
      />
      <button onClick={createBill}>Create Bill and Link Bags</button>
    </div>
  );
};

export default App;
