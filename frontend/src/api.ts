const API_ROOT = "http://localhost:53494/"
const UPLOAD_PATH = API_ROOT + 'upload'

export const uploadFile = async (file: File): Promise<void> => {
    const formData = new FormData()
    formData.append('zip', file)
    await fetch(UPLOAD_PATH, {method: 'POST', body: formData})
}