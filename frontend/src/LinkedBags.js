import React, { useState, useEffect } from "react";
import axios from "axios";

function LinkedBags() {
    const [linkedBags, setLinkedBags] = useState([]);
    const [page, setPage] = useState(1);
    const [loading, setLoading] = useState(false);

    const fetchLinkedBags = async () => {
        setLoading(true);
        try {
            const response = await axios.get(`/linked-bags?page=${page}&limit=10`);
            setLinkedBags(response.data.data);
        } catch (error) {
            console.error("Error fetching linked bags:", error);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchLinkedBags();
    }, [page]);

    return (
        <div>
            <h1>Linked Bags</h1>
            {loading ? <p>Loading...</p> : (
                <ul>
                    {linkedBags.map((bag) => (
                        <li key={bag.id}>
                            Parent Bag: {bag.parent_bag}, Child Bag: {bag.child_bag}
                        </li>
                    ))}
                </ul>
            )}
            <button onClick={() => setPage((prev) => Math.max(prev - 1, 1))}>Previous</button>
            <button onClick={() => setPage((prev) => prev + 1)}>Next</button>
        </div>
    );
}

export default LinkedBags;