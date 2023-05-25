import React, { useState } from "react";
import { create } from "ipfs-http-client";
import "./App.css";

const ipfs = create("/ip4/127.0.0.1/tcp/5001");

function App() {
  const [file, setFile] = useState(null);

  const handleFileChange = (e) => {
    if (e.target.files) {
      setFile(e.target.files[0]);
    }
  };

  const handleUpload = async () => {
    if (!file) return;

    try {
      const fileData = await file.arrayBuffer();
      const pinnedFile = await ipfs.add(fileData);
      console.log(pinnedFile);
    } catch (error) {
      console.error("Error uploading file to IPFS:", error);
    }
  };

  return (
    <div className="App">
      <header className="App-header">
        <div>
          <input type="file" onChange={handleFileChange} />
          <div>{file && `${file.name} - ${file.type}`}</div>
          <button onClick={handleUpload}>Upload</button>
        </div>
      </header>
    </div>
  );
}

export default App;
