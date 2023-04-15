import React, {useEffect, useState} from 'react'
import {Route, Routes} from 'react-router-dom'

import axios from 'axios'

import NavigationBar from './components/NavigationBar'
import Home from './components/Home'
import SignUp from './components/SignUp'
import SignIn from './components/SignIn'

import NoteService from './service/NoteService'

type Note = {
    id: number
    title: string
    content: string
}

const App = () => {
    const [signed, setSigned] = useState(false)
    const [notes, setNotes] = useState<Note[]>([])
    useEffect(() => {
        NoteService.notes({}).then(r => {
            setNotes([...r.data.notes])
            setSigned(true)
        }).catch(err => {
            if (!axios.isAxiosError(err) || !err.response) {
                throw err
            }

            switch (err.response.status) {
                case 404:
                    setSigned(true)
                    break
                default:
                    throw err
            }
        })
    }, [])


    return (
        <>
            <NavigationBar signed={signed}/>
            <Routes>
                <Route path='/' element={<Home signed={signed} notes={notes} setNotes={setNotes}/>}/>
                <Route path='/sign-up' element={<SignUp/>}/>
                <Route path='/sign-in' element={<SignIn/>}/>
            </Routes>
        </>
    );
}

export default App;
