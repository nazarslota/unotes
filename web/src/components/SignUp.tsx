import {FC, useState} from 'react';
import {Navigate, useLocation} from 'react-router-dom';
import {ErrorMessage, Field, Form, Formik, FormikErrors, FormikHelpers} from 'formik';

import './SignUp.css';
import axios from 'axios';
import * as yup from 'yup';

type SignUpProps = {};

const SignUp: FC<SignUpProps> = () => {
    type Values = {
        username: string;
        password: string;
        confirmPassword: string;
        signUpErr: string;
    };

    const location = useLocation();
    const [redirectToSignIn, setRedirectToSignIn] = useState(false);

    const initial: Values = {
        username: '',
        password: '',
        confirmPassword: '',
        signUpErr: '',
    };

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

        if (values.password !== values.confirmPassword) {
            errors.confirmPassword = 'Passwords must match';
        }
        return errors;
    };

    const submit = async (values: Values, actions: FormikHelpers<Values>): Promise<void> => {
        const url = `${process.env.REACT_APP_AUTH_SERVICE_URL}/api/oauth2/sign-up`;
        const data = {
            username: values.username,
            password: values.password,
        };

        try {
            await axios.post(url, data, {
                headers: {'Content-Type': 'application/json', 'Accept': 'application/json'}
            });
            setRedirectToSignIn(true);
        } catch (err) {
            let message: string;
            if (axios.isAxiosError(err) && err.response) {
                switch (err.response.status) {
                    case 400: {
                        message = 'The username or password is incorrect';
                        break;
                    }
                    case 409: {
                        message = 'A user with this username already exists';
                        break;
                    }
                    case 500: {
                        message = 'Server error, please try again later';
                        break;
                    }
                    default: {
                        message = 'Unable to sign up. Unknown error';
                        break;
                    }
                }
            } else {
                message = 'Unable to sign up. Unknown error'
            }
            actions.setErrors({signUpErr: message});
        }
    };

    return <>
        {redirectToSignIn && <Navigate to="/sign-in" state={{from: location}}/>}
        <div className="sign-up">
            <Formik initialValues={initial} validate={validate} onSubmit={submit}>
                {({isSubmitting}) => (
                    <Form className="sign-up__form">
                        <h1 className="sign-up__form__title">Sign Up</h1>
                        <div className="sign-up__form__username">
                            <label className="sign-up__form__username__label">Username</label>
                            <Field className="sign-up__form__username__input" type="text" name="username"
                                   placeholder="Username"/>
                            <ErrorMessage className="sign-up__form__username__error-message" name="username"
                                          component="div"/>
                        </div>
                        <div className="sign-up__form__password">
                            <label className="sign-up__form__password__label">Password</label>
                            <Field className="sign-up__form__password__input" type="password" name="password"
                                   placeholder="Password"/>
                            <ErrorMessage className="sign-up__form__password__error-message" name="password"
                                          component="div"/>
                        </div>
                        <div className="sign-up__form__confirm-password">
                            <label className="sign-up__form__confirm-password__label">Confirm
                                password</label>
                            <Field className="sign-up__form__confirm-password__input" type="password"
                                   name="confirmPassword" placeholder="Confirm password"/>
                            <ErrorMessage className="sign-up__form__confirm-password__error-message"
                                          name="confirmPassword" component="div"/>
                        </div>
                        <div className="sign-up__form__sign-up-error">
                            <ErrorMessage className="sign-up__form__sign-up-error__error-message"
                                          name="signUpErr" component="div"/>
                        </div>
                        <div className="sign-up__form__sign-up">
                            <button className="sign-up-form__sign-up__button" type="submit" disabled={isSubmitting}>
                                Sign Up
                            </button>
                        </div>
                    </Form>
                )}
            </Formik>
        </div>
    </>;
}

export default SignUp;
