import React, {useState} from 'react'
import {ErrorMessage, Field, Form, Formik, FormikHelpers} from 'formik'
import {Link} from 'react-router-dom'

import NoteService from '../service/NoteService'

type Note = {
    id: number
    title: string
    content: string
}

type FormValues = {
    title: string
    content: string
}

type HomeProps = {
    signed: boolean

    notes: Note[]
    setNotes: React.Dispatch<React.SetStateAction<Note[]>>
}

const Home: React.FC<HomeProps> = ({signed, notes, setNotes}) => {
    const [searchTerm, setSearchTerm] = useState('')
    const filteredNotes = notes
        .filter((note) => {
            return note.title.toLowerCase().includes(searchTerm.toLowerCase())
        })
        .sort((noteA, noteB): number => {
            if (noteA.title < noteB.title) {
                return -1
            } else if (noteA.title > noteB.title) {
                return 1
            }
            return 0
        })

    const formInitialValues: FormValues = {title: '', content: ''}
    const formSubmit = async (values: FormValues, actions: FormikHelpers<FormValues>): Promise<void> => {
        try {
            const {data} = await NoteService.create({title: values.title, content: values.content})
            setNotes([...notes, {id: data.id, title: values.title, content: values.content}])
            actions.resetForm()
        } catch (err) {
            // TODO
        }
    }

    const noteDelete = (id: number) => {
        return async () => {
            try {
                await NoteService.delete({id: id})
                setNotes([...notes.filter(note => note.id !== id)])
            } catch (err) {
                // TODO
            }
        }
    }

    return <>
        {signed ? (<div className="flex flex-col md:flex-row">
            <div className="flex justify-center items-start w-full p-4">
                <Formik initialValues={formInitialValues} onSubmit={formSubmit}>
                    {({isSubmitting}) => (<Form className="w-full">
                        <label className="font-semibold w-full text-lg text-center text-gray-700">New Note</label>
                        <div className="mt-3.5">
                            <label className="font-semibold text-gray-700">Title</label>
                            <Field
                                className="text-md text-gray-700 bg-gray-100 p-2 w-full border rounded hover:border-gray-600"
                                type="text" name="title" placeholder="Note title"/>
                            <ErrorMessage className="text-sm text-red-700" name="title" component="div"/>
                        </div>
                        <div className="mt-3.5">
                            <label className="font-semibold text-gray-700">Content</label>
                            <Field
                                className="resize-none text-md text-gray-700 bg-gray-100 p-2 w-full border rounded hover:border-gray-600"
                                as="textarea" type="text" name="content" placeholder="Note content here..."/>
                            <ErrorMessage className="text-sm text-red-700" name="content" component="div"/>
                        </div>
                        <div className="mt-3.5">
                            <button
                                className="w-100 px-3 py-1 border-2 rounded border-indigo-600 text-indigo-600 hover:text-white hover:bg-indigo-500"
                                type="submit" disabled={isSubmitting}>New
                            </button>
                        </div>
                    </Form>)}
                </Formik>
            </div>
            <div className="border-b-2 border-dashed"></div>
            <div className="w-full p-4">
                <label className="font-semibold w-full text-lg text-center text-gray-700">Your Notes</label>
                {notes.length ? (<div className="flex flex-col mt-3.5">
                    <label className="font-semibold text-gray-700">Search</label>
                    <input
                        className="resize-none text-md text-gray-700 bg-gray-100 p-2 w-full border rounded hover:border-gray-600"
                        type="text" placeholder="Search term" value={searchTerm}
                        onChange={(event: React.ChangeEvent<HTMLInputElement>) => setSearchTerm(event.target.value)}/>
                    <ul className="mt-3.5">
                        {filteredNotes.map(note => (<li className="my-2 p-2 border rounded-md" key={note.id}>
                            <div className="flex justify-between">
                                <div>
                                    <span className="text-gray-700 font-semibold">{note.title}</span>
                                </div>
                                <div>
                                    <button className="text-red-700 hover:underline hover:text-red-900"
                                            onClick={noteDelete(note.id)}>
                                        <svg className="h-6 w-6" xmlns="http://www.w3.org/2000/svg" fill="none"
                                             viewBox="0 0 24 24" stroke="currentColor" aria-hidden="true">
                                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2"
                                                  d="M6 18L18 6M6 6l12 12"/>
                                        </svg>
                                    </button>
                                </div>
                            </div>
                            <div className="border-b border-dashed mb-1"></div>
                            <p className="text-gray-700">{note.content}</p>
                        </li>))}
                    </ul>
                </div>) : (<div>
                    <label className="font-semibold w-full text-lg text-center text-red-700">You don't have any notes
                        yet. Add them!</label>
                </div>)}
            </div>
        </div>) : (<div className="flex justify-center items-center vh-100 max-h-96">
            <div>
                <h1 className="font-semibold text-2xl text-gray-700">Unotes</h1>
                <p className="w-80 text-base text-gray-700">Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed
                    eget orci massa. In varius elit nec tortor blandit mollis. Donec tristique, velit non elementum
                    aliquam, velit mauris facilisis leo, quis efficitur metus lectus ut ipsum. <Link
                        className="text-indigo-500 hover:text-indigo-800 hover:underline" to="/sign-in">Sign
                        In</ Link> or <Link className="text-indigo-500 hover:text-indigo-800 hover:underline"
                                            to="/sign-up"> Sign Up</ Link> to use!</p>
            </div>
        </div>)}
    </>
};

export default Home;
