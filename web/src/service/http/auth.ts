import axios from 'axios';

const AUTH_SERVICE_API_URL = `${process.env.REACT_APP_AUTH_SERVICE_URL}/api/oauth2`;

const $api = axios.create({
    baseURL: AUTH_SERVICE_API_URL,
});

export default $api;
