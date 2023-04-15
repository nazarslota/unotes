import {AxiosResponse} from 'axios'

import $api from '../http/note'

type CreateNoteRequest = {
    "title": string
    "content": string
}

type CreateNoteResponse = {
    "id": number
    "userId": string
}

type UpdateNoteRequest = {
    "id": number
    "newTitle": string
    "newContent": string
}

type UpdateNoteResponse = {}

type DeleteNoteRequest = {
    "id": number
}

type DeleteNoteResponse = {}

type GetNoteRequest = {
    "id": string
}

type GetNoteResponse = {
    "title": string
    "content": string
    "userId": string
}

type GetNotesRequest = {};

type Note = {
    "id": number
    "title": string
    "content": string
    "userId": string
}

type GetNotesResponse = {
    "notes": Note[]
};

class NoteService {
    static async create(request: CreateNoteRequest): Promise<AxiosResponse<CreateNoteResponse>> {
        return $api.post('/note', request).then(response => response)
    }

    static async update(request: UpdateNoteRequest): Promise<AxiosResponse<UpdateNoteResponse>> {
        return $api.put('/note', request).then(response => response)
    }

    static async delete(request: DeleteNoteRequest): Promise<AxiosResponse<DeleteNoteResponse>> {
        return $api.delete(`/note/${request["id"]}`).then(response => response)
    }

    static async note(request: GetNoteRequest): Promise<AxiosResponse<GetNoteResponse>> {
        return $api.get(`/note/${request["id"]}`).then(response => response)
    }

    static async notes(_: GetNotesRequest): Promise<AxiosResponse<GetNotesResponse>> {
        return $api.get('/notes').then(response => {
            if (response.data["result"]) {
                const data: GetNotesResponse = {notes: [response.data["result"]]}
                return {...response, data: data}
            }

            const jsons = response.data.trim().split('\n')
            const notes = jsons.map((json: string) => JSON.parse(json)["result"])

            const data: GetNotesResponse = {notes: notes}
            return {...response, data: data}
        })
    }
}

export default NoteService;
