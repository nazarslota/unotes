import React from "react";
import {Link} from "react-router-dom";
import {ErrorMessage, Field, Form, Formik, FormikErrors, FormikHelpers} from "formik";

import axios from "axios";
import * as yup from "yup";

import AuthService from "../services/AuthService";
import utils from "../utils/utils";

export default class SignUp extends React.Component<SignUpT.Props, SignUpT.State> {
    constructor(props: SignUpT.Props) {
        super(props);
        this.state = {
            isUserSignedIn: false,
            isUserSignedUp: false,
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
            {this.state.isUserSignedUp ? (<div className="h-96 p-8 grid place-items-center">
                <div>
                    <span className="text-lg text-gray-600">You can go to the <Link
                        className="text-lg text-gray-800 hover:text-gray-900 hover:underline"
                        to="/sign-in">Sign In</Link> page now!</span>
                </div>
            </div>) : (<Formik
                initialValues={this.formInitial} validate={this.formValidate} onSubmit={this.fromSubmit}
            >
                {({isSubmitting}) => (<div className="mt-40 p-4 w-full flex justify-center"><Form
                    className="p-1.5 border border-gray-600 rounded">
                    <h1 className="font-semibold text-center text-lg text-gray-600">Sign Up</h1>
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
                    <div className="mt-1.5">
                        <label className="font-semibold text-base text-gray-600">Confirm Password</label>
                        <Field
                            className="w-full p-1.5 text-base text-gray-600 hover:text-gray-800 border border-gray-600 hover:border-gray-800 rounded focus:ring-gray-600"
                            type="password" name="confirm" placeholder="Confirm Password"/>
                        <ErrorMessage className="text-sm text-red-600" name="confirm" component="div"/>
                    </div>
                    <ErrorMessage className="text-sm text-red-600" name="err" component="div"/>
                    <button
                        className="w-full mt-2.5 p-1.5 text-center text-gray-600 hover:text-gray-300 border border-gray-600 hover:border-gray-800 hover:bg-gray-800 rounded"
                        type="submit" disabled={isSubmitting}>Sign Up
                    </button>
                </Form></div>)}
            </Formik>)}
        </>);
    }

    private formInitial: SignUpT.FormValues = {
        username: "",
        password: "",
        confirm: "",
        err: ""
    }

    private formValidate = (values: SignUpT.FormValues): FormikErrors<SignUpT.FormValues> => {
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

        if (values.password !== values.confirm) {
            errors.confirm = "Passwords must match!";
        }
        return errors;
    }

    private fromSubmit = (values: SignUpT.FormValues, actions: FormikHelpers<SignUpT.FormValues>): void => {
        AuthService.signUp({username: values.username, password: values.password}).then(_ => {
            this.setState({isUserSignedUp: true});
        }).catch(e => {
            const errors: { [key: string]: string } = {};
            if (!axios.isAxiosError(e) || !e.response) {
                errors.err = "Unknown error!";
            } else if (e.response.status === 400) {
                errors.err = "Username or password is in invalid format!";
            } else if (e.response.status === 409) {
                errors.err = "User with this username already exist!";
            } else if (e.response.status === 500) {
                errors.err = "Server error! Please try again later!";
            } else {
                errors.err = "Unknown error!";
            }

            actions.setSubmitting(false);
            actions.setErrors({err: errors.err});
        })
    }
}

export module SignUpT {
    export type Props = {};

    export type State = {
        isUserSignedIn: boolean;
        isUserSignedUp: boolean;
    };

    export type FormValues = {
        username: string;
        password: string;
        confirm: string;
        err: string;
    };
}
