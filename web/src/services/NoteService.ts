import {AxiosResponse} from "axios";
import $api from "./http/note";

// Requests

type CreateNoteRequest = {
    title: string;
    content: string;
    priority?: string;
    completionTime?: Date;
};

type UpdateNoteRequest = {
    id: number;
    newTitle: string;
    newContent: string;
};

type DeleteNoteRequest = {
    id: number;
};

type GetNoteRequest = {
    id: number;
};

type GetNotesRequest = {};

// Responses

type CreateNoteResponse = {
    id: number;
    userId: string;
};

type UpdateNoteResponse = {};

type DeleteNoteResponse = {};

type GetNoteResponse = {
    title: string;
    content: string;
    userId: string;
    createdAt: Date;
    priority?: string;
    completionTime?: Date;
}

type GetNotesResponse = {
    notes: {
        id: number;
        title: string;
        content: string;
        userId: string;
        createdAt: Date;
        priority?: string;
        completionTime?: Date;
    }[];
}

// NoteService

class NoteService {
    static async create(request: CreateNoteRequest): Promise<AxiosResponse<CreateNoteResponse>> {
        return $api.post("/note", request).then(response => response);
    }

    static async update(request: UpdateNoteRequest): Promise<AxiosResponse<UpdateNoteResponse>> {
        return $api.put("/note", request).then(response => response);
    }

    static async delete(request: DeleteNoteRequest): Promise<AxiosResponse<DeleteNoteResponse>> {
        return $api.delete(`/note/${request["id"]}`).then(response => response);
    }

    static async note(request: GetNoteRequest): Promise<AxiosResponse<GetNoteResponse>> {
        return $api.get(`/note/${request["id"]}`).then(response => {
            const note = response.data;
            note["createdAt"] = new Date(note["createdAt"]);
            if (note["completionTime"]) {
                note["completionTime"] = new Date(note["completionTime"]);
            }

            const data: GetNoteResponse = note;
            return {...response, data: data};
        });
    }

    static async notes(_: GetNotesRequest): Promise<AxiosResponse<GetNotesResponse>> {
        return $api.get("/notes").then(response => {
            if (response.data["result"]) {
                const note = response.data["result"];
                note["createdAt"] = new Date(note["createdAt"]);
                if (note["completionTime"]) {
                    note["completionTime"] = new Date(note["completionTime"]);
                }

                const data: GetNotesResponse = {notes: [note]};
                return {...response, data: data};
            }

            const jsons = response.data.trim().split('\n');
            const notes = jsons.map((json: string) => {
                const note = JSON.parse(json)["result"];
                note["createdAt"] = new Date(note["createdAt"]);
                if (note["completionTime"]) {
                    note["completionTime"] = new Date(note["completionTime"]);
                }
                return note;
            });

            const data: GetNotesResponse = {notes: notes};
            return {...response, data: data};
        })
    }
}

export default NoteService;
