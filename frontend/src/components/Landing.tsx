import * as React from 'react'
import { useState } from 'react'
import { Button } from 'rsuite'
import Dropzone from 'react-dropzone'
import { uploadFile } from '../api'

export default function () {
  const [uploadQueue, setUploadQueue] = useState<File[]>([])

  const uploadFiles = async (): Promise<void> => {
    await Promise.all(uploadQueue.map(file => {
      return uploadFile(file).then(() => setUploadQueue(uploadQueue.filter((f: File): boolean => {
        return file.name !== f.name
      })))
    }))
    hide()
  }

  const addOne = (file: File): void => {
    if (uploadQueue.map(f => f.name).includes(file.name)) {
      console.warn('duplicate file detected')
    } else {
      uploadQueue.push(file)
      setUploadQueue(uploadQueue.map(f => f))
    }
  }

  const addToQueue = (files: File | File[]): void => {
    if (Array.isArray(files)) {
      files.forEach(addOne)
    } else {
      addOne(files)
    }
    setUploadQueue(uploadQueue)
  }

  const hide = () => {
    setTimeout(() => setUploadQueue([]), 500)
  }

  return (
    <div className="App">
      <Dropzone onDrop={(acceptedFiles) => {
        console.log('file dropped')
        addToQueue(acceptedFiles)
      }} accept={'.zip'}>
        {({ getRootProps, getInputProps }) => (
          <section>
            <div {...getRootProps()}>
              <input {...getInputProps()} />
              <h1>Drag 'n' drop some files here, or click to select files</h1>
            </div>
          </section>
        )}
      </Dropzone>
      {
        uploadQueue.map((file, key) => {
          return <p>{file.name}</p>
        })
      }
      <Button onClick={uploadFiles}>Add Mod Files</Button>
    </div>
  )
}
