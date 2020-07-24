import React, {useState} from 'react';
import Modal from 'react-modal';

export default function () {
    let [showModal, setShowModal] = useState<boolean>(false)
    let [result, setResult] = useState<string>("")
    const handleOpenModal = () => {
        setShowModal(true);
        window.backend.Basic().then(setResult);
    };
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