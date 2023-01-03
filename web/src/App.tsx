import React from 'react';
import './App.css';

import {Route, Routes} from 'react-router-dom';
import SignUp from "./components/SignUp";
import SignIn from "./components/SignIn";
import NavigationBar from "./components/NavigationBar";
import Home from "./components/Home";
import CreateNote from "./components/CreateNote";

const App = () => {
    return (
        <>
            <NavigationBar/>
            <Routes>
                <Route path='/' element={<Home/>}/>
                <Route path='/create-note' element={<CreateNote/>}/>
                <Route path='/sign-up' element={<SignUp/>}/>
                <Route path='/sign-in' element={<SignIn/>}/>
            </Routes>
        </>
    )
}

export default App;
