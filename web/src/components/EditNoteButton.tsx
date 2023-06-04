import React from "react";


export default class EditNoteButton extends React.Component<EditNoteButtonT.Props, EditNoteButtonT.State> {
    constructor(props: EditNoteButtonT.Props) {
        super(props);
        this.state = {};
    }

    public render = (): JSX.Element => (<>
        <button
            className={this.props.className + " p-1.5 inline-flex items-center text-base text-gray-600 hover:text-gray-800 hover:underline max-h-8"}
            onClick={this.props.onClick}
        >
            Edit
        </button>
    </>);
}

export module EditNoteButtonT {
    export type Props = {
        className?: string;
        onClick?: () => void;
    };

    export type State = {};
}
