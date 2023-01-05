import React, {FC, FormEvent, useState} from 'react';
import {Navigate, useLocation} from 'react-router-dom';
import {useCookies} from 'react-cookie';
import axios from 'axios';

import './SignIn.css';

type SignInUserResponse = {
    "access_token": string;
    "refresh_token": string;
};

type SignInProps = {};

const SignIn: FC<SignInProps> = () => {
    const [, setCookie] = useCookies(['access_token', 'refresh_token'])

    const [error, setError] = useState<string>("");

    const location = useLocation();
    const [redirectToHome, setRedirectToHome] = useState<boolean>(false);

    const [username, setUsername] = useState<string>("");
    const usernameOnChange = (e: FormEvent<HTMLInputElement>) => {
        setUsername(e.currentTarget.value);
    }

    const [password, setPassword] = useState<string>("");
    const passwordOnChange = (e: FormEvent<HTMLInputElement>) => {
        setPassword(e.currentTarget.value);
    }

    const signInOnClick = async (_: FormEvent<HTMLButtonElement>) => {
        try {
            const response = await axios.post<SignInUserResponse>(
                `${process.env.REACT_APP_AUTH_SERVICE_URL}/api/auth/oauth2/sign-in`,
                {
                    "username": username,
                    "password": password,
                },
                {
                    "headers": {
                        "Content-Type": "application/json",
                        "Accept": "application/json",
                    },
                },
            );

            if (response.status === 200) {
                setError("");
                setRedirectToHome(true);

                setCookie("access_token", response.data["access_token"], {secure: true, sameSite: 'none'});
                setCookie("refresh_token", response.data["refresh_token"], {secure: true, sameSite: 'none'});
            }
        } catch (error) {
            if (axios.isAxiosError(error) && error.response) {
                if (error.response.status === 400 || error.response.status === 404) {
                    setError('Incorrect username or password.');
                    return;
                }
                setError('Unknown error.');
            }
        }
    }

    return (
        <>
            {redirectToHome && <Navigate to="/" state={{from: location}} replace/>}
            <div className="sign-in-form">
                <h1 className="sign-in-form__title">Sign In</h1>
                {error !== "" && <label className="sign-up-form__error">{error}</label>}
                <div className="sign-in-form__username">
                    <label className="sign-in-form__username__label">Username</label>
                    <input
                        className="sign-in-form__username__input"
                        id="sign-in-form__password__input"
                        type="text" value={username} placeholder="Username"
                        onChange={usernameOnChange}
                    />
                </div>
                <div className="sign-in-form__password">
                    <label className="sign-in-form__password__label">Password</label>
                    <input
                        className="sign-in-form__password__input"
                        id="sign-in-form__password__input"
                        type="password" value={password} placeholder="Password"
                        onChange={passwordOnChange}
                    />
                </div>
                <div className="sign-in-form__sign-in">
                    <button className="sign-in-form__sign-in__button" type="submit"
                            onClick={signInOnClick}>Sign In
                    </button>
                </div>
            </div>
        </>
    );
}

export default SignIn;
