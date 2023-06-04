import React from "react";
import {Link} from "react-router-dom";
import {Cookies, withCookies} from "react-cookie";
import {FormikHelpers} from "formik";

import axios from "axios";

import NoteCreationForm, {NoteCreationFormT} from "../components/NoteCreationForm";
import NoteList from "../components/NoteList";

import AuthService from "../services/AuthService";
import NoteService from "../services/NoteService";
import utils from "../utils/utils";


export default withCookies(class Home extends React.Component<HomeT.Props, HomeT.State> {
    constructor(props: HomeT.Props) {
        super(props);
        this.state = {
            isUserSignedIn: false,
            notes: [],
        };
    }

    public componentDidMount = async (): Promise<void> => {
        if (!await utils.isUserSignedIn()) {
            try {
                const {data} = await AuthService.refresh({
                    refresh_token: this.props.cookies.get("refresh_token"),
                });

                localStorage.setItem("access_token", data.access_token);
                this.props.cookies.set("refresh_token", data.refresh_token);

                window.location.reload();
            } catch (e) {
                if (!axios.isAxiosError(e) || !e.response)
                    return;
            }
        } else {
            try {
                const {data} = await NoteService.notes({});
                this.setState({notes: [...data.notes]});
            } catch (e) {
                if (!axios.isAxiosError(e) || !e.response)
                    return;
            }
            this.setState({isUserSignedIn: true});
        }
    };

    public render = (): JSX.Element => (<>
        {!this.state.isUserSignedIn ? (<div className="h-96 p-8 grid place-items-center">
            <div className="w-4/5 md:w-2/5">
                <p className="text-base text-gray-600"><b>Unotes</b> is a notes system that was written using
                    microservice architecture, with Go, React and technologies such as gRPC, MongoDB, Redis,
                    PostgreSQL and others.</p>
                <br/>
                <p className="text-base text-gray-600">
                    <Link className="text-base text-gray-800 hover:text-gray-900 hover:underline" to="/sign-in">
                        Sign In</Link><span> or </span>
                    <Link className="text-base text-gray-800 hover:text-gray-900 hover:underline" to="/sign-up">
                        Sign Up</Link><span> to use!</span>
                </p>
            </div>
        </div>) : (<div className="m-4 md:flex">
            <NoteCreationForm
                className="md:w-5/12"
                formInitial={this.noteCreationFormInitial}
                formSubmit={this.noteCreationFormSubmit}
            />
            <div className="mt-4 border-b border-gray-600 border-dashed md:my-0 md:border-b-0 md:mx-4 md:border-l"/>
            <NoteList
                className="mt-4 md:w-7/12"
                notes={this.state.notes}
                noteOnDelete={this.noteOnDelete}
            />
        </div>)}
    </>);

    private noteCreationFormInitial: NoteCreationFormT.FormValues = {
        title: "",
        content: "",
    };

    private noteCreationFormSubmit = (values: NoteCreationFormT.FormValues, actions: FormikHelpers<NoteCreationFormT.FormValues>): void => {
        const request = !values.priority?.length ? {
            title: values.title, content: values.content,
            completionTime: values.completionTime,
        } : {
            title: values.title, content: values.content,
            priority: values.priority, completionTime: values.completionTime,
        }

        NoteService.create(request).then(r => {
            const note = !values.priority?.length ? {
                id: r.data.id, title: values.title, content: values.content,
                createdAt: new Date(), completionTime: values.completionTime,
            } : {
                id: r.data.id, title: values.title, content: values.content,
                priority: values.priority, createdAt: new Date(), completionTime: values.completionTime,
            }

            this.setState({notes: [...this.state.notes, note]});
            actions.resetForm();
        }).catch(e => {
            actions.resetForm();
            if (!axios.isAxiosError(e) || !e.response) {
                return;
            }
        })
    }

    private noteOnDelete = (id: number): void => {
        NoteService.delete({id: id}).then(_ => {
            this.setState({notes: [...this.state.notes.filter(note => note.id !== id)]});
        }).catch(e => {
            if (!axios.isAxiosError(e) || !e.response) {
                return;
            }
        })
    }
});

export module HomeT {
    export type Props = {
        cookies: Cookies;
    };

    export type State = {
        isUserSignedIn: boolean;
        notes: Note[];
    };

    export type Note = {
        id: number;
        title: string;
        content: string;
        createdAt: Date;
        priority?: string;
        completionTime?: Date;
    };
}
