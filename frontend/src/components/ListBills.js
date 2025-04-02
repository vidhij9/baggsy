import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { motion } from 'framer-motion';
import { toast } from 'react-toastify';
import { DocumentTextIcon } from '@heroicons/react/24/solid';

function ListBills({ setError, token, refresh }) {
  const [bills, setBills] = useState([]);
  const [expanded, setExpanded] = useState({});
  const [expandedParents, setExpandedParents] = useState({});
  const [loading, setLoading] = useState(true);
  const [page, setPage] = useState(1);
  const [limit] = useState(10);
  const [totalPages, setTotalPages] = useState(1);

  useEffect(() => {
    if (token) {
      fetchBills();
    }
  }, [token, page, refresh]); // Add refresh to dependencies

  const fetchBills = async () => {
    setLoading(true);
    try {
      const res = await axios.get(`https://baggsy.app/api/bills?page=${page}&limit=${limit}`, {
        headers: { Authorization: `Bearer ${token}` },
      });
      console.log("Bills response:", res.data);
      setBills(Array.isArray(res.data) ? res.data : []);
      setTotalPages(Math.ceil((res.headers['x-total-count'] || 1) / limit));
      setError(null);
    } catch (err) {
      const errorMsg = err.response?.data?.error || 'Failed to fetch bills';
      console.error("Fetch bills error:", err);
      setBills([]);
      setError(errorMsg);
      toast.error(errorMsg, { position: 'top-center' });
    } finally {
      setLoading(false);
    }
  };

  const toggleExpandBill = (billID) => {
    setExpanded((prev) => ({ ...prev, [billID]: !prev[billID] }));
  };

  const toggleExpandParent = (billID, parentID) => {
    setExpandedParents((prev) => ({
      ...prev,
      [`${billID}-${parentID}`]: !prev[`${billID}-${parentID}`],
    }));
  };

  const handleUnlink = async (bagId) => {
    if (!window.confirm('Are you sure you want to unlink this bag?')) return;
    try {
      await axios.delete(`https://baggsy.app/api/unlink-bag/${bagId}`, {
        headers: { Authorization: `Bearer ${token}` },
      });
      fetchBills();
      toast.success('Bag unlinked successfully!', { position: 'top-center' });
    } catch (err) {
      setError(err.response?.data?.error || 'Failed to unlink bag');
      toast.error(err.response?.data?.error || 'Failed to unlink bag', { position: 'top-center' });
    }
  };

  return (
    <motion.div
      initial={{ opacity: 0, scale: 0.95 }}
      animate={{ opacity: 1, scale: 1 }}
      className="bg-white p-6 rounded-xl shadow-lg border-t-4 border-primary"
    >
      <h2 className="text-2xl font-semibold text-accent mb-4 flex items-center">
        <DocumentTextIcon className="w-6 h-6 text-primary mr-2" />
        List Bills
      </h2>
      <div className="max-h-64 overflow-y-auto">
        {loading ? (
          <p className="text-accent text-center">Loading bills...</p>
        ) : bills.length === 0 ? (
          <p className="text-accent text-center">No bills found.</p>
        ) : (
          bills.map((bill) => (
            <div key={bill.billID} className="border-b py-2">
              <div
                className="flex justify-between items-center cursor-pointer"
                onClick={() => toggleExpandBill(bill.billID)}
              >
                <span>{bill.billID}</span>
                <span>{expanded[bill.billID] ? '▲' : '▼'}</span>
              </div>
              {expanded[bill.billID] && (
                <div className="ml-4 mt-2">
                  {bill.bags && bill.bags.length > 0 ? (
                    bill.bags.map((bag) => (
                      <div key={bag.id} className="border-b py-1">
                        <div
                          className="flex justify-between items-center cursor-pointer"
                          onClick={() => toggleExpandParent(bill.billID, bag.id)}
                        >
                          <span>{bag.qrCode}</span>
                          <span>{expandedParents[`${bill.billID}-${bag.id}`] ? '▲' : '▼'}</span>
                        </div>
                        {expandedParents[`${bill.billID}-${bag.id}`] && bag.children && (
                          <ul className="ml-4">
                            {bag.children.map((child) => (
                              <li key={child.id}>{child.qrCode}</li>
                            ))}
                          </ul>
                        )}
                        <button
                          onClick={() => handleUnlink(bag.id)}
                          className="text-red-500 hover:text-red-700 text-sm mt-1"
                        >
                          Unlink
                        </button>
                      </div>
                    ))
                  ) : (
                    <p className="text-accent">No bags linked to this bill.</p>
                  )}
                </div>
              )}
            </div>
          ))
        )}
      </div>
      <div className="flex justify-between mt-4">
        <button
          onClick={() => setPage((prev) => Math.max(1, prev - 1))}
          disabled={page === 1 || loading}
          className="py-2 px-4 rounded-lg bg-gray-200 disabled:opacity-50"
        >
          Previous
        </button>
        <span>Page {page} of {totalPages}</span>
        <button
          onClick={() => setPage((prev) => Math.min(totalPages, prev + 1))}
          disabled={page === totalPages || loading}
          className="py-2 px-4 rounded-lg bg-gray-200 disabled:opacity-50"
        >
          Next
        </button>
      </div>
    </motion.div>
  );
}

export default ListBills;