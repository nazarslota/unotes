import {AxiosResponse} from 'axios';
import $api from '../http/auth';

type SignUpRequest = {
    "username": string;
    "password": string;
};

type SignUpResponse = {};

type SignInRequest = {
    "username": string;
    "password": string;
};

type SignInResponse = {
    "access_token": string;
    "refresh_token": string;
};

type SignOutRequest = {
    "access_token": string;
};

type SignOutResponse = {};

type RefreshRequest = {
    "refresh_token": string;
};

type RefreshResponse = {};

class AuthService {
    static async signUp(request: SignUpRequest): Promise<AxiosResponse<SignUpResponse>> {
        return $api.post('/oauth2/sign-up', request).then(response => response);
    }

    static async signIn(request: SignInRequest): Promise<AxiosResponse<SignInResponse>> {
        return $api.post('/oauth2/sign-in', request).then(response => response);
    }

    static async signOut(request: SignOutRequest): Promise<AxiosResponse<SignOutResponse>> {
        return $api.post('/oauth2/sign-out', request).then(response => response);
    }

    static async refresh(request: RefreshRequest): Promise<AxiosResponse<RefreshResponse>> {
        return $api.get(`/oauth2/refresh?t=${request["refresh_token"]}`).then(response => response);
    }
}

export default AuthService;
