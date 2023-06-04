import axios from "axios";
import NoteService from "../services/NoteService";

export default class utils {
    public static isUserSignedIn = async () => {
        try {
            await NoteService.notes({});
            return true;
        } catch (e) {
            if (!axios.isAxiosError(e) || !e.response) {
                return false;
            }
            return e.response.status === 404;
        }
    }
}