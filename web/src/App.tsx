import React from "react";
import {Route, Routes} from "react-router-dom";

import NavBar from "./components/NavBar";

import Home from "./views/Home";
import SignUp from "./views/SignUp";
import SignIn from "./views/SignIn";

export default class App extends React.Component {
    public render = (): JSX.Element => (<>
        <NavBar/>
        <Routes>
            <Route path='/' element={<Home/>}/>
            <Route path='/sign-up' element={<SignUp/>}/>
            <Route path='/sign-in' element={<SignIn/>}/>
        </Routes>
    </>);
}
