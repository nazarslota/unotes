import React from "react";
import {ErrorMessage, Field, Form, Formik, FormikErrors, FormikHelpers} from "formik";
import DatePicker from "react-datepicker";

import "react-datepicker/dist/react-datepicker.css";

import * as yup from "yup";

export default class NoteCreationForm extends React.Component<NoteCreationFormT.Props, NoteCreationFormT.State> {
    constructor(props: NoteCreationFormT.Props) {
        super(props);
        this.state = {};
    }

    public render = (): JSX.Element => (<>
        <div className={this.props.className}>
            <Formik
                initialValues={this.props.formInitial} validate={this.formValidate} onSubmit={this.props.formSubmit}
            >
                {({isSubmitting, values, handleChange, setFieldValue}) => (<Form>
                    <h1 className="w-full font-semibold text-center text-lg text-gray-600">New Note</h1>
                    <div className="mt-2">
                        <label className="w-full font-semibold text-base text-gray-600">Title</label>
                        <Field
                            className="w-full p-1.5 text-base text-gray-600 border border-gray-600 rounded focus:border-gray-600 focus:ring-gray-600"
                            type="text" name="title" placeholder="Title"/>
                        <ErrorMessage className="text-xs text-red-600" name="title" component="span"/>
                    </div>
                    <div className="mt-2">
                        <label className="w-full font-semibold text-base text-gray-600">Content</label>
                        <Field
                            className="w-full p-1.5 text-base text-gray-600 border border-gray-600 rounded focus:border-gray-600 focus:ring-gray-600 resize-none h-24 align-top"
                            type="text" name="content" placeholder="Content" as="textarea"/>
                    </div>
                    <div className="mt-2">
                        <h1 className="w-full font-semibold text-base text-gray-600">Completion Time</h1>
                        <DatePicker
                            className="w-full p-1.5 text-base text-gray-600 border border-gray-600 rounded focus:border-gray-600 focus:ring-gray-600"
                            isClearable
                            showTimeSelect
                            dateFormat="MMMM d, yyyy h:mm a"
                            minDate={new Date()}
                            selected={values.completionTime} onChange={date => setFieldValue("completionTime", date)}
                            placeholderText="Completion date and time"
                        />
                    </div>
                    <div className="mt-2">
                        <h1 className="w-full font-semibold text-base text-gray-600">Priority</h1>
                        <div className="w-full lg:flex">
                            <div className="flex">
                                <input className="accent-gray-600" type="radio" value="" name="priority"
                                       onChange={handleChange}
                                       checked={values.priority === undefined || values.priority === ""}/>
                                <label className="ml-2 text-gray-600" htmlFor="">Not Selected</label>
                            </div>
                            <div className="flex lg:ml-2">
                                <input className="accent-gray-600" type="radio" value="lo" name="priority"
                                       onChange={handleChange} checked={values.priority === "lo"}/>
                                <label className="ml-2 text-gray-600" htmlFor="lowPriority">Low</label>
                            </div>
                            <div className="flex lg:ml-2">
                                <input className="accent-gray-600" type="radio" value="md" name="priority"
                                       onChange={handleChange} checked={values.priority === "md"}/>
                                <label className="ml-2 text-gray-600" htmlFor="mediumPriority">Medium</label>
                            </div>
                            <div className="flex lg:ml-2">
                                <input className="accent-gray-600" type="radio" value="hi" name="priority"
                                       onChange={handleChange} checked={values.priority === "hi"}/>
                                <label className="ml-2 text-gray-600" htmlFor="highPriority">High</label>
                            </div>
                        </div>
                    </div>
                    <div className="mt-2 flex justify-center">
                        <button
                            className="w-full p-1.5 text-gray-600 hover:text-gray-100 border border-gray-600 rounded hover:bg-gray-800"
                            type="submit" disabled={isSubmitting}>Submit
                        </button>
                    </div>
                </Form>)}
            </Formik>
        </div>
    </>);

    private formValidate = (values: NoteCreationFormT.FormValues): FormikErrors<NoteCreationFormT.FormValues> => {
        const errors: { [key: string]: string } = {};
        try {
            yup.string().required("Title is required!").validateSync(values.title);
        } catch (e) {
            if (e instanceof yup.ValidationError) {
                errors.title = e.message;
            }
        }
        return errors;
    }
}

export module NoteCreationFormT {
    export type Props = {
        className?: string;

        formInitial: FormValues;
        formSubmit: (values: FormValues, actions: FormikHelpers<FormValues>) => void;
    };

    export type State = {};

    export type FormValues = {
        title: string;
        content: string;
        priority?: string;
        completionTime?: Date;
    };
}
