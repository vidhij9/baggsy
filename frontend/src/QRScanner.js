import { QrReader } from 'react-qr-reader';
import React, { useState } from 'react';

const QRScanner = ({ onScan }) => {
    const [qrCodeData, setQrCodeData] = useState("");

    const handleScan = (data) => {
        if (data) {
            setQrCodeData(data);
            onScan(data); // Pass scanned data to parent component
        }
    };

    const handleError = (err) => {
        console.error(err);
    };

    return (
        <div>
            <QrReader
                delay={300}
                onResult={(result, error) => {
                    if (!!result) {
                        handleScan(result.text);
                    }
                    if (!!error) {
                        handleError(error);
                    }
                }}
                style={{ width: '100%' }}
            />
            <p>Scanned QR Code: {qrCodeData}</p>
        </div>
    );
};

export default QRScanner;