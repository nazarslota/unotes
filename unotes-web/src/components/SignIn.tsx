import React, {FC, FormEvent, useState} from 'react';
import axios from 'axios';

type SignInUserModel = {
    username: string;
    password: string;
};

type SignInUserResponse = {
    "access_token": string;
    "refresh_token": string;
};

const signInUser = async (model: SignInUserModel): Promise<SignInUserResponse> => {
    const {data} = await axios.post<SignInUserResponse>(
        `${process.env.REACT_APP_AUTH_SERVICE_URL}/api/auth/oauth2/sign-in`,
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
}

type SignInProps = {};

const SignIn: FC<SignInProps> = () => {
    const [username, setUsername] = useState<string>("");
    const usernameOnChange = (e: FormEvent<HTMLInputElement>) => {
        setUsername(e.currentTarget.value);
    }

    const [password, setPassword] = useState<string>("");
    const passwordOnChange = (e: FormEvent<HTMLInputElement>) => {
        setPassword(e.currentTarget.value);
    }

    const signInOnClick = async (_: FormEvent<HTMLButtonElement>) => {
        const data = await signInUser({"username": username, "password": password});
        console.log(data["access_token"]);
        console.log(data["refresh_token"]);
    }

    return (
        <>
            <div className="sign-in-form">
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
