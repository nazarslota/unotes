import React from "react";
import ReactModal from "react-modal";
import {ErrorMessage, Field, Form, Formik, FormikErrors, FormikHelpers} from "formik";
import DatePicker from "react-datepicker";

import * as yup from "yup";

export default class EditNoteButton extends React.Component<EditNoteButtonT.Props, EditNoteButtonT.State> {
    constructor(props: EditNoteButtonT.Props) {
        super(props);
        this.state = {
            modalIsOpen: false,
        };
    }

    public render = (): JSX.Element => (<>
        <button
            className={this.props.className + " p-1.5 inline-flex items-center text-base text-gray-600 hover:text-gray-800 hover:underline max-h-8"}
            onClick={this.openModal}
        >Edit
        </button>
        <ReactModal
            style={EditNoteButton.modalStyles}
            isOpen={this.state.modalIsOpen}
            onRequestClose={this.closeModal}
            ariaHideApp={false}
        >
            <Formik
                initialValues={this.formInitial} validate={this.formValidate} onSubmit={this.formSubmit}
            >
                {({isSubmitting, values, handleChange, setFieldValue}) => (<Form>
                    <h1 className="w-full font-semibold text-center text-lg text-gray-600">Edit Note</h1>
                    <div className="mt-2">
                        <label className="w-full font-semibold text-base text-gray-600">New title</label>
                        <Field
                            className="w-full p-1.5 text-base text-gray-600 border border-gray-600 rounded focus:border-gray-600 focus:ring-gray-600"
                            type="text" name="title" placeholder="Title"/>
                        <ErrorMessage className="text-sm text-red-600" name="title" component="span"/>
                    </div>
                    <div className="mt-2">
                        <label className="w-full font-semibold text-base text-gray-600">New content</label>
                        <Field
                            className="w-full p-1.5 text-base text-gray-600 border border-gray-600 rounded focus:border-gray-600 focus:ring-gray-600 resize-none h-24 align-top"
                            type="text" name="content" placeholder="Content" as="textarea"/>
                    </div>
                    <div className="mt-2">
                        <h1 className="w-full font-semibold text-base text-gray-600">Completion time</h1>
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
                        <h1 className="w-full font-semibold text-base text-gray-600">New priority</h1>
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
                            type="submit" disabled={isSubmitting}>Edit
                        </button>
                    </div>
                </Form>)}
            </Formik>
        </ReactModal>
    </>);

    private static modalStyles = {
        content: {
            top: '50%',
            left: '50%',
            right: 'auto',
            bottom: 'auto',
            marginRight: '-50%',
            transform: 'translate(-50%, -50%)',
            borderWidth: '1px',
        },
    };

    private openModal = () => {
        this.setState({modalIsOpen: true});
    }

    private closeModal = () => {
        this.setState({modalIsOpen: false});
    }

    private formInitial: EditNoteButtonT.FormValues = {
        title: this.props.noteTitle,
        content: this.props.noteContent,
        priority: this.props.notePriority,
        completionTime: this.props.noteCompletionTime,
    }

    private formValidate = (values: EditNoteButtonT.FormValues): FormikErrors<EditNoteButtonT.FormValues> => {
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

    private formSubmit = (values: EditNoteButtonT.FormValues, actions: FormikHelpers<EditNoteButtonT.FormValues>) => {
        this.props.onClick?.(this.props.noteId, values.title, values.content, values.priority, values.completionTime);

        actions.resetForm();
        this.closeModal();
    }
}

export module EditNoteButtonT {
    export type Props = {
        noteId: number;
        noteTitle: string;
        noteContent: string;
        notePriority?: string;
        noteCompletionTime?: Date;

        className?: string;
        onClick?: (id: number, title: string, content: string, priority?: string, completionTime?: Date) => void;
    };

    export type State = {
        modalIsOpen: boolean,
    };

    export type FormValues = {
        title: string;
        content: string;
        priority?: string;
        completionTime?: Date;
    };
}
