import {FC} from 'react'
import {useCookies} from 'react-cookie'
import {useNavigate} from 'react-router-dom'
import {ErrorMessage, Field, Form, Formik, FormikErrors, FormikHelpers} from 'formik'

import axios from 'axios'
import * as yup from 'yup'

import AuthService from '../service/AuthService'

type FormValues = {
    username: string
    password: string
    err: string
}

type SignInProps = {};

const SignIn: FC<SignInProps> = () => {
    const navigate = useNavigate()
    const [, setCookie] = useCookies<string, string>(['refresh_token'])

    const formInitialValues: FormValues = {username: '', password: '', err: ''}
    const formValidate = async (values: FormValues): Promise<FormikErrors<FormValues>> => {
        const errors: { [key: string]: string } = {}
        try {
            await yup.string()
                .min(4, "Min username length 4")
                .max(32, "Max username length 32")
                .required("Required")
                .validate(values.username)
        } catch (err) {
            if (err instanceof yup.ValidationError) {
                errors.username = err.message
            } else {
                throw err
            }
        }

        try {
            await yup.string()
                .min(8, "Min password length 8")
                .max(64, "Max password length 64")
                .required("Required")
                .validate(values.password)
        } catch (err) {
            if (err instanceof yup.ValidationError) {
                errors.password = err.message
            } else {
                throw err
            }
        }
        return errors
    }

    const formSubmit = async (values: FormValues, actions: FormikHelpers<FormValues>): Promise<void> => {
        try {
            const {data} = await AuthService.signIn({username: values.username, password: values.password})
            localStorage.setItem('access_token', data.access_token)
            setCookie('refresh_token', data.refresh_token)

            navigate('/')
            window.location.reload()
        } catch (err) {
            if (!axios.isAxiosError(err) || !err.response) {
                throw err
            }

            switch (err.response.status) {
                case 400:
                    actions.setErrors({err: 'Username or password is in invalid format'})
                    break
                case 404:
                    actions.setErrors({err: 'User with this username dose not exist'})
                    break
                case 500:
                    actions.setErrors({err: 'Server error, please try again later'})
                    break
                default:
                    actions.setErrors({err: 'Unable to sign up. Unknown error'})
                    break
            }
        }
    }

    return <>
        <div className="flex justify-center p-4">
            <Formik initialValues={formInitialValues} validate={formValidate} onSubmit={formSubmit}>
                {({isSubmitting}) => (<Form>
                    <div className="flex justify-center">
                        <h1 className="font-bold text-xl text-gray-700">Sign In</h1>
                    </div>
                    <div className="mt-3.5">
                        <label className="font-semibold text-gray-700">Username</label>
                        <Field
                            className="text-md text-gray-700 bg-gray-100 p-2 w-full border rounded hover:border-gray-600"
                            type="text" name="username" placeholder="Username"/>
                        <ErrorMessage className="text-sm text-red-700" name="username" component="div"/>
                    </div>
                    <div className="mt-3.5">
                        <label className="font-semibold text-gray-700">Password</label>
                        <Field
                            className="text-md text-gray-700 bg-gray-100 p-2 w-full border rounded hover:border-gray-600"
                            type="password" name="password" placeholder="Password"/>
                        <ErrorMessage className="text-sm text-red-700" name="password" component="div"/>
                    </div>
                    <div>
                        <ErrorMessage className="text-sm text-red-700" name="signInErr" component="div"/>
                    </div>
                    <div className="flex justify-center mt-3.5">
                        <button
                            className="w-100 px-3 py-1 border-2 rounded border-indigo-600 text-indigo-600 hover:text-white hover:bg-indigo-500"
                            type="submit" disabled={isSubmitting}>Sign In
                        </button>
                    </div>
                </Form>)}
            </Formik>
        </div>
    </>;
}

export default SignIn;
