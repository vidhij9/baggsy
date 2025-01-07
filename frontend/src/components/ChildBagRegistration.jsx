// import React, { useState } from "react";
// import { linkChildBag } from "../api/api";

// const ChildBagRegistration = ({ parentBag, onChildBagsCompleted }) => {
//   const [childBags, setChildBags] = useState([]); // Tracks linked child bags
//   const [qrCode, setQrCode] = useState(""); // Current child bag QR code being scanned
//   const [message, setMessage] = useState(""); // Feedback message

//   const handleAddChildBag = async () => {
//     if (!qrCode) {
//       setMessage("Child Bag QR Code is required!");
//       return;
//     }

//     try {
//       const payload = { parentBag: parentBag.qrCode, childBag: qrCode };
//       const response = await linkChildBag(payload); // Use LinkChildBag API

//       setChildBags((prev) => [...prev, qrCode]); // Add QR code to childBags state
//       setQrCode(""); // Reset the input
//       setMessage(response.data.message); // Show success message

//       // Check if all child bags have been linked
//       if (childBags.length + 1 === parentBag.childCount) {
//         onChildBagsCompleted(); // Notify the parent component
//       }
//     } catch (error) {
//       console.error("Error:", error.response?.data || error.message);
//       setMessage(error.response?.data?.error || "Something went wrong");
//     }
//   };

//   return (
//     <div className="bg-white p-6 rounded shadow-lg">
//       <h2 className="text-2xl font-bold text-primary mb-4">Register Child Bags</h2>
//       <p className="text-gray-600 mb-2">Parent Bag: {parentBag.qrCode}</p>
//       <p className="text-gray-600 mb-4">
//         Remaining: {parentBag.childCount - childBags.length}
//       </p>
//       <input
//         type="text"
//         placeholder="Child Bag QR Code"
//         value={qrCode}
//         onChange={(e) => setQrCode(e.target.value)}
//         className="w-full p-3 border rounded mb-4 focus:outline-none focus:ring-2 focus:ring-primary"
//       />
//       <button
//         onClick={handleAddChildBag}
//         className="bg-primary text-white py-2 px-6 rounded hover:bg-accent transition"
//       >
//         Add Child Bag
//       </button>
//       {message && <p className="text-primary mt-4">{message}</p>}
//     </div>
//   );
// };

// export default ChildBagRegistration;


import React, { useState } from "react";
import { registerBag } from "../api/api";

const ChildBagRegistration = ({ parentBag }) => {
  const [qrCode, setQrCode] = useState("");
  const [message, setMessage] = useState("");

  const handleSubmit = async () => {
    if (!qrCode) {
      setMessage("QR Code is required!");
      return;
    }

    try {
      const payload = { qrCode, bagType: "Child", parentBag: parentBag.qrCode };
      const response = await registerBag(payload);

      setMessage(response.data.message);
      setQrCode("");
    } catch (error) {
      setMessage(error.response?.data?.error || "Something went wrong");
    }
  };

  return (
    <div className="bg-white p-6 rounded-lg shadow-lg max-w-md mx-auto">
      <h2 className="text-2xl font-bold text-primary mb-4">Register Child Bag</h2>
      <input
        type="text"
        placeholder="QR Code"
        value={qrCode}
        onChange={(e) => setQrCode(e.target.value)}
        className="w-full p-3 border rounded mb-4 focus:outline-none focus:ring-2 focus:ring-primary"
      />
      <button
        onClick={handleSubmit}
        className="w-full bg-primary text-white py-3 rounded-lg hover:bg-accent transition-all"
      >
        Submit
      </button>
      {message && <p className="text-primary mt-4">{message}</p>}
    </div>
  );
};

export default ChildBagRegistration;