import React, {FC, useState} from 'react';
import axios from 'axios';

import './Home.css';

type Note = {
    id: string;
    title: string;
    content: string;
};

type HomeProps = {};

const Home: FC<HomeProps> = () => {
    // Initialize an empty list of notes
    const [notes, setNotes] = useState<Note[]>([]);

    // Initialize an empty search term
    const [searchTerm, setSearchTerm] = useState<string>('');

    // Fetch the list of notes from an API or a database
    const fetchNotes = () => {
        axios.get('/api/notes').then(response => {
            setNotes(response.data);
        });
    };

    // Filter the list of notes based on the search term
    const filteredNotes = notes.filter(note => {
        return note.title.includes(searchTerm) || note.content.includes(searchTerm);
    });

    return (
        <div className="home">
            <h1>Welcome to the Notes App</h1>
            <p>Here you can create, view, and edit notes.</p>
            <input
                type="text"
                placeholder="Search notes"
                value={searchTerm}
                onChange={e => setSearchTerm(e.target.value)}
            />
            <ul>
                {filteredNotes.map(note => (
                    <li key={note.id}>
                        <a href={`/view-note/${note.id}`}>{note.title}</a>
                    </li>
                ))}
            </ul>
        </div>
    );
};

export default Home;
