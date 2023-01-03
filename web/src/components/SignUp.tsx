import React, {FC, FormEvent, useState} from 'react';
import {Navigate, useLocation} from 'react-router-dom';
import axios from 'axios';

import './SignUp.css';

type SignUpProps = {};

const SignUp: FC<SignUpProps> = () => {
    const location = useLocation();

    const [error, setError] = useState<string>("");
    const [redirect, setRedirect] = useState<boolean>(false);

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
        if (username.length < 4) {
            setError('Username is too short.');
            return;
        } else if (username.length > 32) {
            setError('Username is too long.');
            return;
        }

        if (password.length < 8) {
            setError('Password is too short.');
            return;
        } else if (password.length > 64) {
            setError('Password is too long.');
            return;
        } else if (password !== confirmPassword) {
            setError('The entered passwords do not match.');
            return;
        }

        try {
            const response = await axios.post(
                `${process.env.REACT_APP_AUTH_SERVICE_URL}/api/auth/oauth2/sign-up`,
                {
                    "username": username,
                    "password": password,
                },
                {
                    "headers": {
                        "Content-Type": "application/json",
                        "Accept": "application/json",
                    },
                }
            );

            if (response.status === 204) {
                setError("");
                setRedirect(true);
            }
        } catch (error) {
            if (axios.isAxiosError(error)) {
                if (error.response && error.response.status === 409) {
                    setError('This username is already in use.')
                } else {
                    setError('Unknown error.')
                }
            }
        }
    }

    return (
        <>
            {redirect && <Navigate to="/sign-in" state={{from: location}} replace/>}
            <div className="sign-up-form">
                <h1 className="sign-up-form__title">Sign Up</h1>
                {error !== "" && <label className="sign-up-form__error">{error}</label>}
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
