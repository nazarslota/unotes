import React from 'react';
import {Route, Routes} from 'react-router-dom';

import NavigationBar from './components/NavigationBar';
import Home from './components/Home';
import SignUp from './components/SignUp';
import SignIn from './components/SignIn';

import './App.css';

const App = () => {
    return (
        <>
            <NavigationBar/>
            <Routes>
                <Route path='/' element={<Home/>}/>
                <Route path='/sign-up' element={<SignUp/>}/>
                <Route path='/sign-in' element={<SignIn/>}/>
            </Routes>
        </>
    );
}

export default App;
