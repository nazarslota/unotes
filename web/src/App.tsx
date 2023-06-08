import React from "react";
import {Route, Routes} from "react-router-dom";

import NavBar from "./components/NavBar";

import Home from "./views/Home";
import SignUp from "./views/SignUp";
import SignIn from "./views/SignIn";
import {AboutMe} from "./views/AboutMe";
import {Footer} from "./components/Footer";

export default class App extends React.Component {
    public render = (): JSX.Element => (<div className="flex flex-col h-screen justify-between">
        <NavBar/>
        <main className="mb-auto">
            <Routes>
                <Route path='/' element={<Home/>}/>
                <Route path='/sign-up' element={<SignUp/>}/>
                <Route path='/sign-in' element={<SignIn/>}/>
                <Route path='/about-me' element={<AboutMe/>}/>
            </Routes>
        </main>
        <Footer/>
    </div>);
}
