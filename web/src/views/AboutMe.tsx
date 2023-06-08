import React from "react";

import {FaUniversity} from "react-icons/fa";
import {BsTelegram} from "react-icons/bs";
import {MdEmail} from "react-icons/md";

type Props = {};

type State = {};

export class AboutMe extends React.Component<Props, State> {
    render() {
        return (<div className="p-10 w-full flex justify-center items-center">
            <div
                className="p-2 flex flex-row justify-center border border-gray-600 rounded space-x-2 max-md:flex-col max-md:items-center max-md:space-x-0">
                <div className="w-2/5 border border-gray-600 rounded max-md:w-full">
                    <img className="rounded" src={require("../resources/img/about-me-avatar.jpg")} alt="Avatar"/>
                </div>
                <div
                    className="p-1 w-3/5 h-1/3 text-gray-600 space-y-2 border border-gray-600 rounded max-md:mt-2 max-md:w-full">
                    <h1 className="text-4xl border-b border-gray-600">Slota Nazar</h1>
                    <p className="flex items-center"><FaUniversity className="mr-1"/> University: Lviv Polytechnic
                        National University</p>
                    <p className="flex items-center"><BsTelegram className="mr-1"/>Telegram: <a
                        className="font-semibold hover:underline" href="https://t.me/nazarslota"> @nazarslota</a></p>
                    <p className="flex items-center"><MdEmail className="mr-1"/>Email: <a
                        className="font-semibold hover:underline"
                        href="mailto: nazarslota03@gmail.com"> nazarslota03@gmail.com</a></p>
                </div>
            </div>
        </div>);
    };
}
