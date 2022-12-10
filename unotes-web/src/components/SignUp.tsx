import React, {FC, FormEvent, useState} from 'react';

type SignUpProps = {};

const SignUp: FC<SignUpProps> = () => {
    const [username, setUsername] = useState<string>("");
    const [password, setPassword] = useState<string>("");
    const [confirmPassword, setConfirmPassword] = useState<string>("");

    const usernameOnChange = (e: FormEvent<HTMLInputElement>) => {
        setUsername(e.currentTarget.value);
    }

    const passwordOnChange = (e: FormEvent<HTMLInputElement>) => {
        setPassword(e.currentTarget.value);
    }

    const confirmPasswordOnChange = (e: FormEvent<HTMLInputElement>) => {
        setConfirmPassword(e.currentTarget.value);
    }

    const signUpOnClick = (e: FormEvent<HTMLButtonElement>) => {
        console.log(`Username: ${username}; Password: ${password}; Confirm Password: ${confirmPassword}`);
    }

    return (
        <div className="sign-up-form">
            <div className="sign-up-form__username">
                <label className="sign-up-form__username__label">Username</label>
                <input
                    className="sign-up-form__username__input" id="sign-up-form__password__input"
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
            <div className="sign-up-form__signup">
                <button className="sign-up-form__signup__button" type="submit" onClick={signUpOnClick}>SignUp</button>
            </div>
        </div>
    );
}

export default SignUp;
