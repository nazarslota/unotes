import React from "react";
import {Link} from "react-router-dom";

type Props = {};

type State = {};

export class Footer extends React.Component<Props, State> {
    render() {
        return (<div className="w-full p-1.5 flex justify-center border-t border-gray-600">
            <span className="text-gray-600">© <Link className="font-semibold text-gray-800 hover:underline" to="/about-me">Slota Nazar</Link> 2023</span>
        </div>);
    };
}
