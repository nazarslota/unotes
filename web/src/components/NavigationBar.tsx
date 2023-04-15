import React, {FC} from 'react';
import {Link} from 'react-router-dom';
import {useCookies} from 'react-cookie';

import AuthService from '../service/AuthService'

type NavigationBarProps = {
    signed: boolean
}

const NavigationBar: FC<NavigationBarProps> = ({signed}) => {
    const [, , removeCookie] = useCookies<string, string>(['refresh_token'])

    const handleSignOut = async () => {
        try {
            await AuthService.signOut({access_token: localStorage.getItem('access_token') || ''})
            localStorage.removeItem('access_token')
            removeCookie('access_token')

            window.location.reload()
        } catch (err) {
            throw err
        }
    }

    return <>
        <header className="flex items-center justify-between px-3 border-b-2">
            <Link className="font-bold text-2xl text-indigo-500 hover:text-indigo-800 hover:underline"
                  to='/'>Unotes</Link>
            <nav>
                {signed ? (<ul className="text-gray-500 font-semibold inline-flex items-center">
                    <li>
                        <button className="inline-block py-3 px-2 text-xl hover:text-indigo-800 hover:underline"
                                onClick={handleSignOut}>Sign Out
                        </button>
                    </li>
                </ul>) : (<ul className="text-gray-500 font-semibold inline-flex items-center">
                    <li><Link className="inline-block py-3 px-2 text-xl hover:text-indigo-800 hover:underline"
                              to='/sign-in'>Sign In</Link></li>
                    <li><Link className="inline-block py-3 px-2 text-xl hover:text-indigo-800 hover:underline"
                              to='/sign-up'>Sign Up</Link></li>
                </ul>)}
            </nav>
        </header>
    </>;
};

export default NavigationBar;
