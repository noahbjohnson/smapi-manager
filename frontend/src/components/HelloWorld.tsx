import React, { useState } from 'react'
import Modal from 'react-modal'

export default function () {
  const [showModal, setShowModal] = useState<boolean>(false)
  const [result, setResult] = useState<string>('')
  const handleOpenModal = () => {
    setShowModal(true)
    window.backend.Initialize().then(setResult)
  }
  const handleCloseModal = () => setShowModal(false)
  return (
    <div className="App">
      <button onClick={handleOpenModal} type="button">
                Hello
      </button>
      <Modal
        isOpen={showModal}
        contentLabel="Minimal Modal Example"
      >
        <p>{result}</p>
        <button onClick={handleCloseModal}>Close Modal</button>
      </Modal>
    </div>
  )
}
