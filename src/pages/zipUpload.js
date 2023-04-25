import React, { useState } from 'react';

function FileInput({ handleSubmit }) {
    //const [selectedFile, setSelectedFile] = useState(null);
    const [b64String, setB64String] = useState('')

    const handleFileInputChange = (event) => {
        const file = event.target.files[0];
      //  setSelectedFile(file);
        const reader = new FileReader();
        reader.onload = () => {
            const newB64String = btoa(reader.result);
            setB64String(newB64String);
            handleSubmit(newB64String);
        };
        reader.readAsBinaryString(file);
    };

    const handleFormSubmit = (event) => {
        event.preventDefault();
        sendPostRequest(b64String);
        /*if (selectedFile) {
            const reader = new FileReader();
            reader.onload = (event) => {
                const fileData = event.target.result;
                const base64Data = btoa(fileData);
                handleSubmit(base64Data);
            };
            reader.readAsBinaryString(selectedFile);
        }*/
    };

    const sendPostRequest = (b64String) => {
        fetch(`${process.env.NEXT_PUBLIC_BACKEND_API}`, {
            method: 'POST',
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ Content: b64String})
        })
    };//response?

    return (
        <form onSubmit={handleFormSubmit}>
            <label htmlFor="file-input">
                Select a zip file:
                <input
                    type="file"
                    id="file-input"
                    accept="application/zip"
                    onChange={handleFileInputChange}
                    style={{ display: 'none' }}
                />
            </label>
            <button type="submit">Submit</button>
        </form>
    );
}

export default FileInput;