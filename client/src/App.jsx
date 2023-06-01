import React, { useState } from 'react'

import './App.css'

import FileList from './components/FileList'
import FileUpload from './components/FileUpload'

function App() {
  const [uploadedFiles, setUploadedFiles] = useState([])

  return (
    <div className="container">
      <div className="row justify-content-center pt-5 pb-4">
        <div className="col-12 col-md-10 col-lg-7 col-xl-6">
          <h4>ERC-721-Checks IPFS pin service</h4>
          <FileUpload setUploadedFiles={setUploadedFiles} />
        </div>
      </div>
      <div className="row justify-content-center">
        <div className="col-12 col-md-10 col-lg-7 col-xl-6">
          <FileList uploadedFiles={uploadedFiles} />
        </div>
      </div>
    </div>
  )
}

export default App
