import React, {useState} from 'react';

import './Home.css';

type Note = {
    id: number;
    title: string;
    content: string;
};

type HomeProps = {};

const Home: React.FC<HomeProps> = () => {
    const [notes, setNotes] = useState<Note[]>([]);
    const [searchTerm, setSearchTerm] = useState('');
    const [newNoteTitle, setNewNoteTitle] = useState('');
    const [newNoteContent, setNewNoteContent] = useState('');

    const handleSubmit = (event: React.FormEvent<HTMLFormElement>) => {
        event.preventDefault();
        const newNote: Note = {
            id: notes.length + 1,
            title: newNoteTitle,
            content: newNoteContent,
        };
        setNotes([...notes, newNote]);
        setNewNoteTitle('');
        setNewNoteContent('');
    };

    const handleSearch = (event: React.ChangeEvent<HTMLInputElement>) => setSearchTerm(event.target.value);
    const filteredNotes = notes.filter(note => note.title.toLowerCase().includes(searchTerm.toLowerCase()));

    return (
        <>
            <div className="home">
                <div className="home__description">
                    <h1 className="home__description-title">Note App</h1>
                    <p className="home__description-text">
                        Welcome to the Note App! Use the form below to add a new note, and use the search field to the
                        right to filter through your notes.
                    </p>
                </div>
                <form onSubmit={handleSubmit} className="home__form">
                    <input
                        type="text"
                        placeholder="Title"
                        value={newNoteTitle}
                        onChange={event => setNewNoteTitle(event.target.value)}
                        className="home__form-title"
                    />
                    <textarea
                        placeholder="Content"
                        value={newNoteContent}
                        onChange={event => setNewNoteContent(event.target.value)}
                        className="home__form-content"
                    />
                    <button type="submit" className="home__form-button">
                        Add Note
                    </button>
                </form>
                <div className="home__list">
                    <input
                        type="text"
                        placeholder="Search"
                        value={searchTerm}
                        onChange={handleSearch}
                        className="home__list-search"
                    />
                    {filteredNotes.map(note => (
                        <div key={note.id} className="home__list-item">
                            <h3 className="home__list-item-title">{note.title}</h3>
                            <p className="home__list-item-content">{note.content}</p>
                        </div>
                    ))}
                </div>
            </div>
        </>
    );
};

export default Home;
