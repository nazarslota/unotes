import React from "react";
import ReactModal from "react-modal";
import {ErrorMessage, Field, Form, Formik, FormikErrors, FormikHelpers} from "formik";
import DatePicker from "react-datepicker";

import * as yup from "yup";

export default class EditNoteButton extends React.Component<EditNoteButtonT.Props, EditNoteButtonT.State> {
    constructor(props: EditNoteButtonT.Props) {
        super(props);
        this.state = {modalIsOpen: false};
    }

    public render = (): JSX.Element => (<>
        <button
            className={this.props.className + " p-1.5 inline-flex items-center text-base text-gray-600 hover:text-gray-800 hover:underline max-h-8"}
            onClick={this.openModal}>Edit
        </button>
        <ReactModal
            style={EditNoteButton.modalStyles}
            isOpen={this.state.modalIsOpen}
            onRequestClose={this.closeModal}
            ariaHideApp={false}
        >
            <Formik
                initialValues={this.props.formInitial} validate={this.formValidate}
                onSubmit={this.formSubmit}
            >
                {({isSubmitting, values, handleChange, setFieldValue}) => (<Form>
                    <h1 className="w-full font-semibold text-center text-lg text-gray-600">Edit Note</h1>
                    <div className="mt-2">
                        <label className="w-full font-semibold text-base text-gray-600">New title</label>
                        <Field
                            className="w-full p-1.5 text-base text-gray-600 border border-gray-600 rounded focus:border-gray-600 focus:ring-gray-600"
                            type="text" name="newTitle" placeholder="New title"/>
                        <ErrorMessage className="text-sm text-red-600" name="newTitle" component="span"/>
                    </div>
                    <div className="mt-2">
                        <label className="w-full font-semibold text-base text-gray-600">New content</label>
                        <Field
                            className="w-full p-1.5 text-base text-gray-600 border border-gray-600 rounded focus:border-gray-600 focus:ring-gray-600 resize-none h-24 align-top"
                            type="text" name="newContent" placeholder="New content" as="textarea"/>
                    </div>
                    <div className="mt-2">
                        <h1 className="w-full font-semibold text-base text-gray-600">New completion time</h1>
                        <DatePicker
                            className="w-full p-1.5 text-base text-gray-600 border border-gray-600 rounded focus:border-gray-600 focus:ring-gray-600"
                            isClearable showTimeSelect dateFormat="MMMM d, yyyy h:mm a" minDate={new Date()}
                            selected={values.newCompletionTime}
                            onChange={date => setFieldValue("newCompletionTime", date)}
                            placeholderText="Completion date and time"/>
                    </div>
                    <div className="mt-2">
                        <h1 className="w-full font-semibold text-base text-gray-600">New priority</h1>
                        <div className="w-full lg:flex">
                            <div className="flex">
                                <input className="accent-gray-600" type="radio" value="" name="newPriority"
                                       onChange={handleChange}
                                       checked={values.newPriority === undefined || values.newPriority === ""}/>
                                <label className="ml-2 text-gray-600" htmlFor="">Not Selected</label>
                            </div>
                            <div className="flex lg:ml-2">
                                <input className="accent-gray-600" type="radio" value="lo" name="newPriority"
                                       onChange={handleChange} checked={values.newPriority === "lo"}/>
                                <label className="ml-2 text-gray-600" htmlFor="lowPriority">Low</label>
                            </div>
                            <div className="flex lg:ml-2">
                                <input className="accent-gray-600" type="radio" value="md" name="newPriority"
                                       onChange={handleChange} checked={values.newPriority === "md"}/>
                                <label className="ml-2 text-gray-600" htmlFor="mediumPriority">Medium</label>
                            </div>
                            <div className="flex lg:ml-2">
                                <input className="accent-gray-600" type="radio" value="hi" name="newPriority"
                                       onChange={handleChange} checked={values.newPriority === "hi"}/>
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

    private formValidate = (values: EditNoteButtonT.FormValues): FormikErrors<EditNoteButtonT.FormValues> => {
        const errors: { [key: string]: string } = {};
        try {
            yup.string().required("Title is required!").validateSync(values.newTitle);
        } catch (e) {
            if (e instanceof yup.ValidationError) {
                errors.newTitle = e.message;
            }
        }
        return errors;
    }

    private formSubmit = (values: EditNoteButtonT.FormValues, actions: FormikHelpers<EditNoteButtonT.FormValues>) => {
        const update = values.newPriority?.length ? {...values} : {
            newTitle: values.newTitle,
            newContent: values.newContent,
            newCompletionTime: values.newCompletionTime
        };
        this.props.onEdit?.({...update});

        this.closeModal();
        actions.resetForm();
    }
}

export module EditNoteButtonT {
    export type Props = {
        className?: string;
        formInitial: FormValues;

        onEdit?: (update: {
            newTitle: string,
            newContent: string,
            newPriority?: string,
            newCompletionTime?: Date,
        }) => void;
    };

    export type State = {
        modalIsOpen: boolean,
    };

    export type FormValues = {
        newTitle: string;
        newContent: string;
        newPriority?: string;
        newCompletionTime?: Date;
    };
}
