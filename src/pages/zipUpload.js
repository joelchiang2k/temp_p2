import React, { useState } from 'react';

function FileInput({ handleSubmit }) {
    const [selectedFile, setSelectedFile] = useState(null);

    const handleFileInputChange = (event) => {
        const file = event.target.files[0];
        setSelectedFile(file);
    };

    const handleFormSubmit = (event) => {
        event.preventDefault();

        if (selectedFile) {
            const reader = new FileReader();
            reader.onload = (event) => {
                const fileData = event.target.result;
                const base64Data = btoa(fileData);
                handleSubmit(base64Data);
            };
            reader.readAsBinaryString(selectedFile);
        }
    };

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