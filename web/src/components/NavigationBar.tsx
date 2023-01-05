import React, {FC, useState} from 'react';
import {Link} from 'react-router-dom';

import './NavigationBar.css';

type NavigationBarProps = {};

const NavigationBar: FC<NavigationBarProps> = () => {
    const [isSignInOpen, setIsSignInOpen] = useState(false);
    const [isRegistrationOpen, setIsRegistrationOpen] = useState(false);
    const [isCreateNoteOpen, setIsCreateNoteOpen] = useState(false);

    const toggleSignInModal = () => setIsSignInOpen(!isSignInOpen);
    const toggleRegistrationModal = () => setIsRegistrationOpen(!isRegistrationOpen);
    const toggleCreateNoteModal = () => setIsCreateNoteOpen(!isCreateNoteOpen);

    return (
        <>
            <nav className="navigation-bar">
                <div className="navigation-bar__left">
                    <button className="navigation-bar__button" onClick={toggleCreateNoteModal}>
                        <Link className="navigation-bar__button__link" to='/'>Home</Link>
                    </button>
                    <button className="navigation-bar__button" onClick={toggleCreateNoteModal}>
                        <Link className="navigation-bar__button__link" to='/create-note'>Create Note</Link>
                    </button>
                </div>
                <div className="navigation-bar__right">
                    <button className="navigation-bar__button" onClick={toggleSignInModal}>
                        <Link className="navigation-bar__button__link" to='/sign-in'>Sign In</Link>
                    </button>
                    <button className="navigation-bar__button" onClick={toggleRegistrationModal}>
                        <Link className="navigation-bar__button__link" to='/sign-up'>Sign Up</Link>
                    </button>
                </div>
            </nav>
        </>
    );
};

export default NavigationBar;
