import {FC, useState} from 'react';
import {useCookies} from 'react-cookie';
import {ErrorMessage, Field, Form, Formik, FormikErrors, FormikHelpers} from 'formik';

import './SignIn.css';
import axios from 'axios';
import * as yup from 'yup';

type SignInProps = {};

const SignIn: FC<SignInProps> = () => {
    type Values = {
        username: string;
        password: string;
        signInErr: string;
    };

    const [, setCookie] = useCookies<string, string>(['access_token', 'refresh_token']);
    const [signedIn, setSignedIn] = useState(false);

    const initial: Values = {username: '', password: '', signInErr: ''};
    const validate = (values: Values): FormikErrors<Values> => {
        const errors: { [key: string]: string } = {};
        try {
            yup.string()
                .min(4, "Min length 4")
                .max(32, "Max length 32")
                .required("Required")
                .validateSync(values.username)
        } catch (err) {
            if (err instanceof yup.ValidationError) {
                errors.username = err.message
            }
        }

        try {
            yup.string()
                .min(8, "Min length 8")
                .max(64, "Max length 64")
                .required("Required")
                .validateSync(values.password)
        } catch (err) {
            if (err instanceof yup.ValidationError) {
                errors.password = err.message
            }
        }
        return errors;
    }

    const submit = async (values: Values, actions: FormikHelpers<Values>): Promise<void> => {
        const url = `${process.env.REACT_APP_AUTH_SERVICE_URL}/api/oauth2/sign-in`;
        const data = {
            username: values.username,
            password: values.password,
        };

        try {
            const response = await axios.post(url, data, {
                headers: {'Content-Type': 'application/json', 'Accept': 'application/json'}
            });

            setCookie('access_token', response.data['access_token']);
            setCookie('refresh_token', response.data['refresh_token']);

            setSignedIn(true);
        } catch (err) {
            let message: string;
            if (axios.isAxiosError(err) && err.response) {
                switch (err.response.status) {
                    case 400: {
                        message = 'The username or password is incorrect';
                        break;
                    }
                    case 404: {
                        message = 'The username or password is incorrect';
                        break;
                    }
                    case 500: {
                        message = 'Server error, please try again later';
                        break;
                    }
                    default: {
                        message = 'Unable to sign in. Unknown error';
                        break;
                    }
                }
            } else {
                message = 'Unable to sign in. Unknown error';
            }
            actions.setErrors({signInErr: message});
        }
    }

    return <>
        <div className="sign-in">
            {signedIn ? (
                <div className="sign-in__signed-in">
                    <label className="sign-in__signed-in__label">
                        Login successful! You can now go to the homepage.</label>
                </div>
            ) : (
                <Formik initialValues={initial} validate={validate} onSubmit={submit}>
                    {({isSubmitting}) => (
                        <Form className="sign-in__form">
                            <h1 className="sign-in__form__title">Sign In</h1>
                            <div className="sign-in__form__username">
                                <label className="sign-in__form__username__label">Username</label>
                                <Field className="sign-in__form__username__input" type="text" name="username"
                                       placeholder="Username"/>
                                <ErrorMessage className="sign-in__form__username__error-message" name="username"
                                              component="div"/>
                            </div>
                            <div className="sign-in__form__password">
                                <label className="sign-in__form__password__label">Password</label>
                                <Field className="sign-in__form__password__input" type="password" name="password"
                                       placeholder="Password"/>
                                <ErrorMessage className="sign-in__form__password__error-message" name="password"
                                              component="div"/>
                            </div>
                            <div className="sign-in__sign-in-error">
                                <ErrorMessage className="sign-in-form__sign-in-error__error-message"
                                              name="signInErr" component="div"/>
                            </div>
                            <div className="sign-in-form__sign-in">
                                <button className="sign-in-form__sign-in__button" type="submit" disabled={isSubmitting}>
                                    Sign In
                                </button>
                            </div>
                        </Form>
                    )}
                </Formik>
            )}
        </div>
    </>;
}

export default SignIn;
