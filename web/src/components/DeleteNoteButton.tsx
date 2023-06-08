import React from "react";
import ReactModal from "react-modal";

export default class DeleteNoteButton extends React.Component<DeleteNoteButtonT.Props, DeleteNoteButtonT.State> {
    constructor(props: DeleteNoteButtonT.Props) {
        super(props);
        this.state = {modalIsOpen: false};
    }

    public render = (): JSX.Element => (<>
        <button
            className={this.props.className + " p-1.5 inline-flex items-center text-base text-gray-600 hover:text-gray-100 hover:bg-red-800 border border-red-600 hover:border-red-800 rounded max-h-8"}
            onClick={this.openModal}>Delete
        </button>
        <ReactModal style={DeleteNoteButton.modalStyles} isOpen={this.state.modalIsOpen}
                    onRequestClose={this.closeModal} ariaHideApp={false}>
            <h1 className="w-full font-semibold text-base text-gray-600">Are you sure you want to delete the current
                note?</h1>
            <div className="mt-2 flex justify-between">
                <button
                    className="w-full  p-1.5 inline-flex items-center text-base text-gray-600 hover:text-gray-100 hover:bg-red-800 border border-red-600 hover:border-red-800 rounded max-h-8"
                    onClick={this.props.onDelete}><p className="w-full text-center">Yes</p>
                </button>
                <button
                    className="w-full text-center ml-2 p-1.5 inline-flex items-center text-base text-gray-600 hover:text-gray-100 hover:bg-gray-800 border border-gray-600 hover:border-gray-800 rounded max-h-8"
                    onClick={this.closeModal}><p className="w-full text-center">Cancel</p>
                </button>
            </div>
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
}

export module DeleteNoteButtonT {
    export type Props = {
        className?: string;

        onDelete?: () => void;
    };

    export type State = {
        modalIsOpen: boolean;
    };
}
