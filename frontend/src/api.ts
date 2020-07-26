import request from "superagent";

const API_ROOT = "http://localhost:53494/"
const UPLOAD_PATH = API_ROOT + 'upload'
const SMAPI_PATH = API_ROOT + 'smapi'

export const uploadFile = async (file: File): Promise<void> => {
    const formData = new FormData()
    formData.append('zip', file)
    await fetch(UPLOAD_PATH, {method: 'POST', body: formData})
}

export async function hasSMAPI(): Promise<boolean> {
    const res = await request.get(SMAPI_PATH)
    return parseInt(res.text) === 1;
}