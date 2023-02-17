import {AxiosResponse} from 'axios';
import $api from '../http/note';

type CreateNoteRequest = {
    "title": string;
    "content": string;
};

type CreateNoteResponse = {
    "id": string;
    "userId": string;
};

type UpdateNoteRequest = {
    "id": string;
    "newTitle": string;
    "newContent": string;
};

type UpdateNoteResponse = {};

type DeleteNoteRequest = {
    "id": string;
};

type DeleteNoteResponse = {};

type GetNoteRequest = {
    "id": string;
};

type GetNoteResponse = {
    "title": string;
    "content": string;
    "userId": string;
};

type GetNotesRequest = {};

type GetNotesResponse = {
    "notes": { "id": string; "title": string; "content": string; "userId": string; }[];
};

class NoteService {
    static async create(request: CreateNoteRequest): Promise<AxiosResponse<CreateNoteResponse>> {
        return $api.post('/note', request).then(request => request);
    }

    static async update(request: UpdateNoteRequest): Promise<AxiosResponse<UpdateNoteResponse>> {
        return $api.put('/note', request).then(request => request);
    }

    static async delete(request: DeleteNoteRequest): Promise<AxiosResponse<DeleteNoteResponse>> {
        return $api.delete(`/note/${request["id"]}`).then(response => response);
    }

    static async note(request: GetNoteRequest): Promise<AxiosResponse<GetNoteResponse>> {
        return $api.get(`/note/${request["id"]}`).then(response => response);
    }

    static async notes(_: GetNotesRequest): Promise<AxiosResponse<GetNotesResponse>> {
        return $api.get(`/notes`).then(response => response);
    }
}

export default NoteService;
