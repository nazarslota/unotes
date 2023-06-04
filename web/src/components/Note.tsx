import React from "react";
import DeleteNoteButton from "./DeleteNoteButton";
import EditNoteButton from "./EditNoteButton";

import moment from "moment";

import * as models from "../models/models";

export default class Note extends React.Component<NoteT.Props, NoteT.State> {
    constructor(props: NoteT.Props) {
        super(props);
        this.state = {};
    }

    public render = (): JSX.Element => (<>
        <div className={this.props.className}>
            <div className="p-2 border border-gray-600 rounded">
                <div className="flex">
                    <div
                        className="w-full inline-flex items-center space-x-2 justify-between max-[1000px]:flex-col max-[1000px]:items-start">
                        <div className="flex items-center">
                            {this.props.note.priority && (
                                (this.props.note.priority.toLowerCase().includes('lo') &&
                                    <div className="mr-2 p-1.5 rounded-full bg-green-600"></div>
                                ) || (this.props.note.priority.toLowerCase().includes('md') &&
                                    <div className="mr-2 p-1.5 rounded-full bg-orange-600"></div>
                                ) || (this.props.note.priority.toLowerCase().includes('hi') &&
                                    <div className="mr-2 p-1.5 rounded-full bg-red-600"></div>
                                )
                            )}
                            <h1 className="my-auto font-semibold text-lg text-gray-600">{this.props.note.title}</h1>
                        </div>
                        <div className="ml-5 my-auto text-gray-600 max-[1000px]:ml-0" style={{marginLeft: 0}}>{
                            (this.props.note.completionTime && <span>{
                                `${moment(this.props.note.createdAt).format("MMMM D, YYYY h:MM A")} - ${moment(this.props.note.completionTime).format("MMMM D, YYYY h:MM A")}`
                            }</span>) || (<span>{
                                `${moment(this.props.note.createdAt).format("MMMM D, YYYY h:MM A")} - Not Specified`
                            }</span>)
                        }</div>
                    </div>
                    <div className="ml-2 border-l border-gray-600 border-dashed"/>
                    <div
                        className="flex items-center justify-between max-[1000px]:items-center max-[300px]:flex-col-reverse max-[300px]:justify-center">
                        <EditNoteButton
                            className="ml-2"
                            formInitial={({
                                newTitle: this.props.note.title,
                                newContent: this.props.note.content,
                                newPriority: this.props.note.priority,
                                newCompletionTime: this.props.note.completionTime
                            })}
                            onEdit={(update: {
                                newTitle: string,
                                newContent: string,
                                newPriority?: string,
                                newCompletionTime?: Date,
                            }) => this.props.onEdit?.(this.props.note.id, update)}
                        />
                        <DeleteNoteButton className="ml-2" onDelete={() => this.props.onDelete?.(this.props.note.id)}/>
                    </div>
                </div>
                <div className="my-1.5 border-b border-gray-600 border-dashed"></div>
                <p className="w-full text-base text-gray-600">{this.props.note.content}</p>
            </div>
        </div>
    </>);
}

export module NoteT {
    export type Props = {
        className?: string;
        note: models.Note;

        onEdit?: (id: string, update: {
            newTitle: string,
            newContent: string,
            newPriority?: string,
            newCompletionTime?: Date,
        }) => void;
        onDelete?: (id: string) => void;
    };

    export type State = {};
}
