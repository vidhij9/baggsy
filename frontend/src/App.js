import React, { useState, useEffect } from 'react';
import axios from 'axios';

const App = () => {
  const [bags, setBags] = useState([]);

  useEffect(() => {
    axios.get('/bags').then((res) => setBags(res.data.bags));
  }, []);

  return (
    <div>
      <h1>Baggsy</h1>
      <ul>
        {bags.map((bag) => (
          <li key={bag.id}>{bag.qr_code}</li>
        ))}
      </ul>
    </div>
  );
};

export default App;
