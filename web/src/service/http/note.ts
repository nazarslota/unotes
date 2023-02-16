import axios from 'axios';

const NOTE_SERVICE_API_URL = `${process.env.REACT_APP_NOTE_SERVICE_URL}/api`;

const $api = axios.create({
    baseURL: NOTE_SERVICE_API_URL,
});

export default $api;
