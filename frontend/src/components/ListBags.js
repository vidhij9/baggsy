import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { motion } from 'framer-motion';
import { toast } from 'react-toastify';
import { ArchiveBoxIcon } from '@heroicons/react/24/solid';

function ListBags({ setError, token }) {
    const [bags, setBags] = useState([]);
    const [filters, setFilters] = useState({ type: '', startDate: '', endDate: '', unlinked: false, page: 1, limit: 10 });
    const [expanded, setExpanded] = useState({});
    const [loading, setLoading] = useState(true);
    const [totalPages, setTotalPages] = useState(1);

    useEffect(() => {
        if (token) {
            fetchBags();
        }
    }, [filters, token]);

    const fetchBags = async () => {
        setLoading(true);
        try {
            const params = new URLSearchParams({
                type: filters.type,
                startDate: filters.startDate,
                endDate: filters.endDate,
                unlinked: filters.unlinked.toString(),
                page: filters.page.toString(),
                limit: filters.limit.toString(),
            }).toString();
            const res = await axios.get(`http://baggsy-env.eba-z5m26a8j.ap-south-1.elasticbeanstalk.com/api/bags?${params}`, {
                headers: { Authorization: `Bearer ${token}` },
            });
            if (res.data.message === "No bags found.") {
                setBags([]);
                setTotalPages(1);
            } else {
                setBags(Array.isArray(res.data) ? res.data : []);
                setTotalPages(Math.ceil((res.headers['x-total-count'] || 1) / filters.limit));
            }
            setError(null);
        } catch (err) {
            const errorMsg = err.response?.data?.error || 'Failed to fetch bags';
            console.error("Fetch bags error:", err);
            setBags([]);
            setError(errorMsg);
            toast.error(errorMsg, { position: 'top-center' });
        } finally {
            setLoading(false);
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
                    onChange={(e) => setFilters({ ...filters, type: e.target.value, page: 1 })}
                    className="w-full p-2 border rounded-lg"
                >
                    <option value="">All Types</option>
                    <option value="parent">Parent</option>
                    <option value="child">Child</option>
                </select>
                <input
                    type="date"
                    value={filters.startDate}
                    onChange={(e) => setFilters({ ...filters, startDate: e.target.value, page: 1 })}
                    className="w-full p-2 border rounded-lg"
                />
                <input
                    type="date"
                    value={filters.endDate}
                    onChange={(e) => setFilters({ ...filters, endDate: e.target.value, page: 1 })}
                    className="w-full p-2 border rounded-lg"
                />
                <label className="flex items-center">
                    <input
                        type="checkbox"
                        checked={filters.unlinked}
                        onChange={(e) => setFilters({ ...filters, unlinked: e.target.checked, page: 1 })}
                        className="mr-2 accent-primary"
                    />
                    Unlinked Parents Only
                </label>
                <div className="max-h-64 overflow-y-auto">
                    {loading ? (
                        <p className="text-accent text-center">Loading bags...</p>
                    ) : bags.length === 0 ? (
                        <p className="text-accent text-center">No bags found.</p>
                    ) : (
                        bags.map((bag) => (
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
                                        {bag.bag.type === 'parent' && bag.children && bag.children.length > 0 ? (
                                            <ul>
                                                {bag.children.map((child) => (
                                                    <li key={child.id}>{child.qrCode}</li>
                                                ))}
                                            </ul>
                                        ) : bag.bag.type === 'parent' ? (
                                            <p>No child bags linked.</p>
                                        ) : null}
                                        {bag.bag.type === 'child' && bag.parentQR && (
                                            <p>Parent: {bag.parentQR}</p>
                                        )}
                                        {bag.billID && <p>Bill ID: {bag.billID}</p>}
                                    </div>
                                )}
                            </div>
                        ))
                    )}
                </div>
                <div className="flex justify-between">
                    <button
                        onClick={() => setFilters({ ...filters, page: Math.max(1, filters.page - 1) })}
                        disabled={filters.page === 1 || loading}
                        className="py-2 px-4 rounded-lg bg-gray-200 disabled:opacity-50"
                    >
                        Previous
                    </button>
                    <span>Page {filters.page} of {totalPages}</span>
                    <button
                        onClick={() => setFilters({ ...filters, page: Math.min(totalPages, filters.page + 1) })}
                        disabled={filters.page === totalPages || loading}
                        className="py-2 px-4 rounded-lg bg-gray-200 disabled:opacity-50"
                    >
                        Next
                    </button>
                </div>
            </div>
        </motion.div>
    );
}

export default ListBags;