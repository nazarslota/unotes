import React from "react";
import Note from "./Note";

import * as models from "../models/models";

export default class NoteList extends React.Component<NoteListT.Props, NoteListT.State> {
    constructor(props: NoteListT.Props) {
        super(props);
        this.state = {};
    }

    public render = (): JSX.Element => (<>
        <div className={this.props.className}>
            {!this.props.notes.length ? (<div className="w-full font-semibold text-center text-sm text-red-600">
                <p>You don't have any notes yet. Add them!</p>
            </div>) : (<div>
                <h1 className="w-full font-semibold text-center text-lg text-gray-600">Your Notes</h1>
                <ul className="mt-2">
                    {this.props.notes.map(note => (<li className="mt-2" key={note.id}><Note
                        note={({
                            id: note.id,
                            title: note.title,
                            content: note.content,
                            createdAt: note.createdAt,
                            priority: note.priority,
                            completionTime: note.completionTime,
                        })}
                        onEdit={this.props.noteOnEdit}
                        onDelete={this.props.noteOnDelete}
                    /></li>))}
                </ul>
            </div>)}
        </div>
    </>);
}

export module NoteListT {
    export type Props = {
        className?: string;
        notes: Note[];

        noteOnEdit?: (id: string, update: {
            newTitle: string,
            newContent: string,
            newPriority?: string,
            newCompletionTime?: Date,
        }) => void;
        noteOnDelete?: (id: string) => void;
    };

    export type State = {};

    export type Note = models.Note;
}
