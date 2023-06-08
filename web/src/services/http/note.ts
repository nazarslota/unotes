import axios from "axios";

const NOTE_SERVICE_API_URL = `${process.env.REACT_APP_NOTE_URL}/api`;

const $api = axios.create({
    baseURL: NOTE_SERVICE_API_URL,
    headers: {"Authorization": "Bearer " + localStorage.getItem("access_token")},
});

export default $api;
