import React from "react";
import Note from "./Note";

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
                        id={note.id}
                        title={note.title}
                        content={note.content}
                        createdAt={note.createdAt}
                        priority={note.priority}
                        completionTime={note.completionTime}
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
        noteOnDelete: (id: number) => void;
    };

    export type State = {};

    export type Note = {
        id: number;
        title: string;
        content: string;
        createdAt: Date;
        priority?: string;
        completionTime?: Date;
    };
}
