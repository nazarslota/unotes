import {AxiosResponse} from "axios";
import $api from "./http/auth";

// Requests

type SignUpRequest = {
    username: string;
    password: string;
};

type SignInRequest = {
    username: string;
    password: string;
};

type SignOutRequest = {
    access_token: string;
};

type RefreshRequest = {
    refresh_token: string;
};

// Responses

type SignUpResponse = {};

type SignInResponse = {
    access_token: string;
    refresh_token: string;
};

type SignOutResponse = {};

type RefreshResponse = {
    access_token: string;
    refresh_token: string;
};

// AuthService

class AuthService {
    static signUp = async (request: SignUpRequest): Promise<AxiosResponse<SignUpResponse>> => {
        return $api.post("/oauth2/sign-up", request).then(response => response);
    };

    static signIn = async (request: SignInRequest): Promise<AxiosResponse<SignInResponse>> => {
        return $api.post("/oauth2/sign-in", request).then(response => response);
    };

    static signOut = async (request: SignOutRequest): Promise<AxiosResponse<SignOutResponse>> => {
        return $api.post("/oauth2/sign-out", request).then(response => response);
    };

    static refresh = async (request: RefreshRequest): Promise<AxiosResponse<RefreshResponse>> => {
        return $api.get(`/oauth2/refresh?t=${request["refresh_token"]}`).then(response => response);
    };
}

export default AuthService;
