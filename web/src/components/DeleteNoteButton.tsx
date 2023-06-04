import React from "react";


export default class DeleteNoteButton extends React.Component<DeleteNoteButtonT.Props, DeleteNoteButtonT.State> {
    constructor(props: DeleteNoteButtonT.Props) {
        super(props);
        this.state = {};
    }

    public render = (): JSX.Element => (<>
        <button
            className={this.props.className + " p-1.5 inline-flex items-center text-base text-gray-600 hover:text-gray-100 hover:bg-red-800 border border-red-600 hover:border-red-800 rounded max-h-8"}
            onClick={this.props.onClick}
        >
            Delete
        </button>
    </>);
}

export module DeleteNoteButtonT {
    export type Props = {
        className?: string;
        onClick?: () => void;
    };

    export type State = {};
}
