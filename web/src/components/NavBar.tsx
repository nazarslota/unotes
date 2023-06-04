import React from "react";
import {Link} from "react-router-dom";
import {Cookies, withCookies} from "react-cookie";

import utils from "../utils/utils";

export default withCookies(class NavBar extends React.Component<NavBarT.Props, NavBarT.State> {
    constructor(props: NavBarT.Props) {
        super(props);
        this.state = {
            isUserSignedIn: false,
        };
    }

    public componentDidMount = async (): Promise<void> => {
        try {
            const isSignedIn = await utils.isUserSignedIn();
            this.setState({isUserSignedIn: isSignedIn});
        } catch (e) {
            this.setState({isUserSignedIn: false});
        }
    };

    public render = (): JSX.Element => (<>
        <div className={this.props.className}>
            <header className="px-2 py-1.5 flex justify-between border-b border-gray-600 rounded-b">
                <div className="my-auto">
                    <Link className="font-bold text-2xl text-gray-600 hover:text-gray-800 hover:underline"
                          to="/">Unotes</Link>
                </div>
                <nav className="my-auto">
                    {this.state.isUserSignedIn ? (<ul className="flex justify-between space-x-4">
                        <li>
                            <button
                                className="text-lg text-gray-600 hover:text-gray-800 hover:underline"
                                onClick={this.onSignOut}>Sign Out
                            </button>
                        </li>
                    </ul>) : (<ul className="flex justify-between space-x-4">
                        <li><Link className="text-lg text-gray-600 hover:text-gray-800 hover:underline"
                                  to="/sign-in">Sign In</Link></li>
                        <li><Link className="text-lg text-gray-600 hover:text-gray-800 hover:underline"
                                  to="/sign-up">Sign Up</Link></li>
                    </ul>)}
                </nav>
            </header>
        </div>
    </>);

    private onSignOut = (): void => {
        localStorage.removeItem("access_token");
        this.props.cookies.remove('refresh_token');

        window.location.reload();
    }
});

export module NavBarT {
    export type Props = {
        className?: string;
        cookies: Cookies;
    };

    export type State = {
        isUserSignedIn: boolean;
    };
}
