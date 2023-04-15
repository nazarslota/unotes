import axios from 'axios';

const AUTH_SERVICE_API_URL = `${process.env.REACT_APP_AUTH_URL}/api`;

const $api = axios.create({
    baseURL: AUTH_SERVICE_API_URL,
});

export default $api;
