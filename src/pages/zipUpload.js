import React, { useState } from 'react';

function FileInput({ handleSubmit, onZipUpload }) {
    //const [selectedFile, setSelectedFile] = useState(null);
    const [b64String, setB64String] = useState('')
    const [rows, setRows] = useState([]);
    const [postID, setPostID] = useState();
    const [postName, setPostName] = useState('');
    const [postVersion, setPostVersion] = useState('');

    const addRow = () => {
        const newRow = { ID: postID, Name: postName, Version: postVersion };
        setRows([...rows, newRow]);
        setPostID('');
        setPostName('');
        setPostVersion('');
      };

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
        fetch(`${process.env.NEXT_PUBLIC_BACKEND_API}` + "/package", {
            method: 'POST',
            headers: { "Content-Type": "application/json",
                        "X-Authorization": "bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"

            },
            body: JSON.stringify({ Content: b64String})
        })
        .then(response => response.json())
        .then(responseJSON => onZipUpload(responseJSON))
        .catch(error => console.error(error));
        
        /*.then(responseJSON => {
            console.log(responseJSON);
            console.log(responseJSON.metadata)
            const addRow = {ID: responseJSON.metadata.ID, Name: responseJSON.metadata.Name, Version: responseJSON.metadata.Version};
            setRows([...rows, addRow]);
      }).catch(error => console.error(error));*/
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