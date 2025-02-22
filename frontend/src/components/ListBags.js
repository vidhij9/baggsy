import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { motion } from 'framer-motion';
import { toast } from 'react-toastify';
import { ArchiveBoxIcon } from '@heroicons/react/24/solid';

function ListBags({ setError }) {
  const [bags, setBags] = useState([]);
  const [filters, setFilters] = useState({ type: '', startDate: '', endDate: '', unlinked: false });
  const [expanded, setExpanded] = useState({});

  useEffect(() => {
    fetchBags();
  }, [filters]);

  const fetchBags = async () => {
    try {
      const params = new URLSearchParams(filters).toString();
      const res = await axios.get(`http://localhost:8080/api/bags?${params}`);
      setBags(res.data);
      setError(null);
    } catch (err) {
      setError(err.response?.data?.error || 'Failed to fetch bags');
      toast.error(err.response?.data?.error || 'Failed to fetch bags', { position: 'top-center' });
    }
  };

  const toggleExpand = (id) => {
    setExpanded((prev) => ({ ...prev, [id]: !prev[id] }));
  };

  return (
    <motion.div
      initial={{ opacity: 0, scale: 0.95 }}
      animate={{ opacity: 1, scale: 1 }}
      className="bg-white p-6 rounded-xl shadow-lg border-t-4 border-primary"
    >
      <h2 className="text-2xl font-semibold text-accent mb-4 flex items-center">
        <ArchiveBoxIcon className="w-6 h-6 text-primary mr-2" />
        List Bags
      </h2>
      <div className="space-y-4">
        <select
          value={filters.type}
          onChange={(e) => setFilters({ ...filters, type: e.target.value })}
          className="w-full p-2 border rounded-lg"
        >
          <option value="">All Types</option>
          <option value="parent">Parent</option>
          <option value="child">Child</option>
        </select>
        <input
          type="date"
          value={filters.startDate}
          onChange={(e) => setFilters({ ...filters, startDate: e.target.value })}
          className="w-full p-2 border rounded-lg"
        />
        <input
          type="date"
          value={filters.endDate}
          onChange={(e) => setFilters({ ...filters, endDate: e.target.value })}
          className="w-full p-2 border rounded-lg"
        />
        <label className="flex items-center">
          <input
            type="checkbox"
            checked={filters.unlinked}
            onChange={(e) => setFilters({ ...filters, unlinked: e.target.checked })}
            className="mr-2 accent-primary"
          />
          Unlinked Parents Only
        </label>
        <div className="max-h-64 overflow-y-auto">
          {bags.map((bag) => (
            <div key={bag.bag.id} className="border-b py-2">
              <div
                className="flex justify-between items-center cursor-pointer"
                onClick={() => toggleExpand(bag.bag.id)}
              >
                <span>{bag.bag.qrCode} ({bag.bag.type})</span>
                <span>{expanded[bag.bag.id] ? '▲' : '▼'}</span>
              </div>
              {expanded[bag.bag.id] && (
                <div className="ml-4 mt-2">
                  {bag.bag.type === 'parent' && bag.children && (
                    <ul>
                      {bag.children.map((child) => (
                        <li key={child.id}>{child.qrCode}</li>
                      ))}
                    </ul>
                  )}
                  {bag.bag.type === 'child' && bag.parentQR && (
                    <p>Parent: {bag.parentQR}</p>
                  )}
                  {bag.billID && <p>Bill ID: {bag.billID}</p>}
                </div>
              )}
            </div>
          ))}
        </div>
      </div>
    </motion.div>
  );
}

export default ListBags;