import React from "react";
import {Link} from "react-router-dom";
import {ErrorMessage, Field, Form, Formik, FormikErrors, FormikHelpers} from "formik";
import {Cookies, withCookies} from "react-cookie";

import axios from "axios";
import * as yup from "yup";

import AuthService from "../services/AuthService";
import utils from "../utils/utils";

export default withCookies(class SignIn extends React.Component<SignInT.Props, SignInT.State> {
    constructor(props: SignInT.Props) {
        super(props);
        this.state = {
            isUserSignedIn: false,
        };
    }

    public componentDidMount = async (): Promise<void> => {
        try {
            const isSignedIn = await utils.isUserSignedIn();
            this.setState({isUserSignedIn: isSignedIn});
        } catch (e) {
            this.setState({isUserSignedIn: false});
        }
    };

    public render(): JSX.Element {
        return (<>
            {this.state.isUserSignedIn ? (<div className="h-96 p-8 grid place-items-center">
                <p className="text-lg text-gray-600">You are signed in. You can go to the <Link
                    className="text-lg text-gray-800 hover:text-gray-900 hover:underline"
                    to="/">Home</Link> page now!</p>
            </div>) : (<Formik
                initialValues={this.formInitial} validate={this.formValidate} onSubmit={this.fromSubmit}
            >
                {({isSubmitting}) => (<div className="mt-40 p-4 w-full flex justify-center"><Form
                    className="p-1.5 border border-gray-600 rounded">
                    <h1 className="font-semibold text-center text-lg text-gray-600">Sign In</h1>
                    <div>
                        <label className="font-semibold text-base text-gray-600">Username</label>
                        <Field
                            className="w-full p-1.5 text-base text-gray-600 hover:text-gray-800 border border-gray-600 hover:border-gray-800 rounded focus:ring-gray-600"
                            type="text" name="username" placeholder="Username"/>
                        <ErrorMessage className="text-sm text-red-600" name="username" component="div"/>
                    </div>
                    <div className="mt-1.5">
                        <label className="font-semibold text-base text-gray-600">Password</label>
                        <Field
                            className="w-full p-1.5 text-base text-gray-600 hover:text-gray-800 border border-gray-600 hover:border-gray-800 rounded focus:ring-gray-600"
                            type="password" name="password" placeholder="Password"/>
                        <ErrorMessage className="text-sm text-red-600" name="password" component="div"/>
                    </div>
                    <ErrorMessage className="text-sm text-red-600" name="err" component="div"/>
                    <button
                        className="w-full mt-2.5 p-1.5 text-center text-gray-600 hover:text-gray-300 border border-gray-600 hover:border-gray-800 hover:bg-gray-800 rounded"
                        type="submit" disabled={isSubmitting}>Sign In
                    </button>
                </Form></div>)}
            </Formik>)}
        </>);
    }

    private formInitial: SignInT.FormValues = {
        username: "",
        password: "",
        err: ""
    }

    private formValidate = (values: SignInT.FormValues): FormikErrors<SignInT.FormValues> => {
        const errors: { [key: string]: string } = {};
        try {
            yup.string()
                .min(4, "Min username length 4")
                .max(32, "Max username length 32")
                .required("Required")
                .validateSync(values.username);
        } catch (e) {
            if (e instanceof yup.ValidationError) {
                errors.username = e.message;
            }
        }

        try {
            yup.string()
                .min(8, "Min password length 8")
                .max(64, "Max password length 64")
                .required("Required")
                .validateSync(values.password);
        } catch (e) {
            if (e instanceof yup.ValidationError) {
                errors.password = e.message;
            }
        }
        return errors;
    }

    private fromSubmit = (values: SignInT.FormValues, actions: FormikHelpers<SignInT.FormValues>): void => {
        AuthService.signIn({username: values.username, password: values.password}).then(r => {
            localStorage.setItem("access_token", r.data["access_token"]);
            this.props.cookies.set("refresh_token", r.data["refresh_token"]);

            window.location.reload();
        }).catch(e => {
            const errors: { [key: string]: string } = {};
            if (!axios.isAxiosError(e) || !e.response) {
                errors.err = "Unknown error!";
            } else if (e.response.status === 400) {
                errors.err = "Invalid username or password!";
            } else if (e.response.status === 404) {
                errors.err = "User with this username does not exist!";
            } else if (e.response.status === 500) {
                errors.err = "Server error! Please try again later!";
            } else {
                errors.err = "Unknown error!";
            }

            actions.setSubmitting(false);
            actions.setErrors({err: errors.err});
        })
    }
});

export module SignInT {
    export type Props = {
        cookies: Cookies
    };

    export type State = {
        isUserSignedIn: boolean;
    };

    export type FormValues = {
        username: string;
        password: string;
        err: string;
    };
}
