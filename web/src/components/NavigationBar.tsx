import React, {FC, useState} from 'react';

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
                    <button className="navigation-bar__button" onClick={toggleCreateNoteModal}>Create Note</button>
                </div>
                <div className="navigation-bar__right">
                    <button className="navigation-bar__button" onClick={toggleSignInModal}>Sign In</button>
                    <button className="navigation-bar__button" onClick={toggleRegistrationModal}>Sign Up</button>
                </div>
            </nav>
        </>
    );
};

export default NavigationBar;