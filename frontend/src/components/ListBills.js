import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { motion } from 'framer-motion';
import { toast } from 'react-toastify';
import { DocumentTextIcon } from '@heroicons/react/24/solid';

function ListBills({ setError }) {
  const [bills, setBills] = useState([]);
  const [expanded, setExpanded] = useState({});
  const [expandedParents, setExpandedParents] = useState({});

  useEffect(() => {
    fetchBills();
  }, []);

  const fetchBills = async () => {
    try {
      const res = await axios.get('http://localhost:8080/api/bills');
      setBills(res.data);
      setError(null);
    } catch (err) {
      setError(err.response?.data?.error || 'Failed to fetch bills');
      toast.error(err.response?.data?.error || 'Failed to fetch bills', { position: 'top-center' });
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
        {bills.map((bill) => (
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
                {bill.bags.map((bag) => (
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
                      onClick={async () => {
                        try {
                          await axios.delete(`http://localhost:8080/api/unlink-bag/${bag.id}`);
                          fetchBills();
                          toast.success('Bag unlinked successfully!', { position: 'top-center' });
                        } catch (err) {
                          setError(err.response?.data?.error || 'Failed to unlink bag');
                          toast.error(err.response?.data?.error || 'Failed to unlink bag', {
                            position: 'top-center',
                          });
                        }
                      }}
                      className="text-red-500 hover:text-red-700 text-sm mt-1"
                    >
                      Unlink
                    </button>
                  </div>
                ))}
              </div>
            )}
          </div>
        ))}
      </div>
    </motion.div>
  );
}

export default ListBills;