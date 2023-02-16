import {AxiosResponse} from 'axios';
import $api from './http/auth';

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
        return $api.post('/sign-up', request).then(response => response.data);
    }

    static async signIn(request: SignInRequest): Promise<AxiosResponse<SignInResponse>> {
        return $api.post('/sign-in', request).then(response => response.data);
    }

    static async signOut(request: SignOutRequest): Promise<AxiosResponse<SignOutResponse>> {
        return $api.post('/sign-out', request).then(response => response.data);
    }

    static async refresh(request: RefreshRequest): Promise<AxiosResponse<RefreshResponse>> {
        return $api.get(`/refresh?t=${request["refresh_token"]}`).then(response => response.data);
    }
}

export default AuthService;
