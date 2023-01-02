import React, {FC, FormEvent, useState} from 'react';
import axios from 'axios';

import './SignUp.css'

type SignUpUserModel = {
    username: string;
    password: string;
};

type SignUpUserResponse = {};

const signUpUser = async (model: SignUpUserModel): Promise<SignUpUserResponse> => {
    const {data} = await axios.post<SignUpUserResponse>(
        `${process.env.REACT_APP_AUTH_SERVICE_URL}/api/auth/oauth2/sign-up`,
        {
            "username": model.username,
            "password": model.password,
        },
        {
            "headers": {
                "Content-Type": "application/json",
                "Accept": "application/json",
            },
        },
    );
    return data;
};

type SignUpProps = {};

const SignUp: FC<SignUpProps> = () => {
    const [username, setUsername] = useState<string>("");
    const usernameOnChange = (e: FormEvent<HTMLInputElement>) => {
        setUsername(e.currentTarget.value);
    }

    const [password, setPassword] = useState<string>("");
    const passwordOnChange = (e: FormEvent<HTMLInputElement>) => {
        setPassword(e.currentTarget.value);
    }

    const [confirmPassword, setConfirmPassword] = useState<string>("");
    const confirmPasswordOnChange = (e: FormEvent<HTMLInputElement>) => {
        setConfirmPassword(e.currentTarget.value);
    }

    const signUpOnClick = async (_: FormEvent<HTMLButtonElement>) => {
        if (password !== confirmPassword) {
            return; // TODO
        }
        await signUpUser({username: username, password: password}); // TODO
    }

    return (
        <>
            <div className="sign-up-form">
                <h1 className="sign-up-form__title">Sign Up</h1>
                <div className="sign-up-form__username">
                    <label className="sign-up-form__username__label">Username</label>
                    <input
                        className="sign-up-form__username__input"
                        id="sign-up-form__password__input"
                        type="text" value={username} placeholder="Username"
                        onChange={usernameOnChange}
                    />
                </div>
                <div className="sign-up-form__password">
                    <label className="sign-up-form__password__label">Password</label>
                    <input
                        className="sign-up-form__password__input" id="sign-up-form__password__input"
                        type="password" value={password} placeholder="Password"
                        onChange={passwordOnChange}
                    />
                </div>
                <div className="sign-up-form__confirm-password">
                    <label className="sign-up-form__confirm-password__label">Confirm Password</label>
                    <input
                        className="sign-up-form__confirm-password__input" id="sign-up-form__confirm-password__input"
                        type="password" value={confirmPassword} placeholder="Confirm Password"
                        onChange={confirmPasswordOnChange}
                    />
                </div>
                <div className="sign-up-form__sign-up">
                    <button className="sign-up-form__sign-up__button" type="submit"
                            onClick={signUpOnClick}>Sign Up
                    </button>
                </div>
            </div>
        </>
    );
}

export default SignUp;
