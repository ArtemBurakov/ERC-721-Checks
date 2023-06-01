import React, { useState, useRef } from 'react'
import { create } from 'ipfs-http-client'

import ToastError from './ToastError'

const ipfs = create('/ip4/127.0.0.1/tcp/5001')

function FileUpload({ setUploadedFiles }) {
  const fileInputRef = useRef(null)

  const [file, setFile] = useState(null)
  const [error, setError] = useState(null)
  const [showToast, setShowToast] = useState(false)

  const handleFileChange = (e) => {
    if (e.target.files) setFile(e.target.files[0])
  }

  const toggleShowToast = () => setShowToast(!showToast)

  const handleUpload = async () => {
    if (!file) return

    try {
      const fileData = await file.arrayBuffer()
      const pinnedFile = await ipfs.add(fileData)
      const uploadedFile = {
        name: file.name,
        cid: pinnedFile.cid.toString(),
      }

      setUploadedFiles((prevUploadedFiles) => [
        ...prevUploadedFiles,
        uploadedFile,
      ])

      setFile(null)
      if (fileInputRef.current) fileInputRef.current.value = null

      if (showToast) toggleShowToast()
    } catch (error) {
      console.error('Error uploading file to IPFS:', error)
      setError(
        'An error occurred while trying to reach the IPFS node. Please try again later.'
      )
      if (!showToast) toggleShowToast()
    }
  }

  return (
    <div className="d-flex">
      <input
        type="file"
        className="form-control"
        aria-label="Upload"
        onChange={handleFileChange}
        ref={fileInputRef}
      />
      <button
        className="btn btn-primary ms-2"
        type="button"
        onClick={handleUpload}
      >
        Upload
      </button>
      <ToastError
        show={showToast}
        toggleShow={toggleShowToast}
        position="bottom-end"
        background="light"
        message={error}
      />
    </div>
  )
}

export default FileUpload
