import React, {FC, useState} from 'react';
import {Link} from 'react-router-dom';

import './NavigationBar.css';

type NavigationBarProps = {};

const NavigationBar: FC<NavigationBarProps> = () => {
    const [isRegistrationOpen, setIsSignUpOpen] = useState(false);
    const toggleSignUpModal = () => setIsSignUpOpen(!isRegistrationOpen);

    const [isSignInOpen, setIsSignInOpen] = useState(false);
    const toggleSignInModal = () => setIsSignInOpen(!isSignInOpen);

    const [isHomeModalOpen, setIsHomeModalOpen] = useState(false);
    const toggleHomeModal = () => setIsHomeModalOpen(!isHomeModalOpen);

    return <>
        <nav className="navigation-bar">
            <div className="navigation-bar__left">
                <Link className="navigation-bar__left__button navigation-bar__left__button__link" to='/'
                      onClick={toggleHomeModal}>Home</Link>
            </div>
            <div className="navigation-bar__right">
                <Link className="navigation-bar__right__button navigation-bar__right__button__link" to='/sign-in'
                      onClick={toggleSignInModal}>Sign In</Link>
                <Link className="navigation-bar__right__button navigation-bar__right__button__link" to='/sign-up'
                      onClick={toggleSignUpModal}>Sign Up</Link>
            </div>
        </nav>
    </>;
};

export default NavigationBar;
