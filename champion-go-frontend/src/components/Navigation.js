import React, {useState} from "react"
import {useDispatch, useSelector} from "react-redux"
import {logout} from "../store/modules/auth/actions/authActions"
import Default from "../assets/default.png"
import {
    Collapse,
    Navbar,
    NavbarToggler,
    NavbarBrand,
    Nav,
    NavItem,
    UncontrolledDropdown,
    DropdownToggle,
    DropdownMenu,
    DropdownItem
} from 'reactstrap';
import {Link, NavLink} from "react-router-dom";
import './Navigation.css'

const Navigation = () => {
    const [isOpen, setIsOpen] = useState(false)
    const currentState = useSelector(state => state)
    const {isAuthenticated, currentUser} = currentState.Auth

    const dispatch = useDispatch()

    const logoutUser = () => dispatch(logout())

    let imagePreview = null;
    if (currentUser && currentUser.avatar_path) {
        imagePreview = (<img className="img_style_nav" src={currentUser.avatar_path} alt="profile 1"/>);
    } else {
        imagePreview = (<img className="img_style_nav" src={Default} alt="profile 2"/>);
    }

    const logout = (e) => {
        e.preventDefault()
        logoutUser()
    }

    const userProfile = isAuthenticated ? `/profile/${currentState.Auth.currentUser.id}` : ""

    const signedInLinks = (
        <React.Fragment>
            <NavItem className="mt-2" style={{marginLeft: "15px"}}>
                <NavLink to="/createpost">Create Post</NavLink>
            </NavItem>
            <NavItem className="mt-2" style={{marginRight: "15px"}}>
                <NavLink to="/authposts">My Posts</NavLink>
            </NavItem>
            <UncontrolledDropdown nav inNavbar>
                <DropdownToggle nav caret>
                    {imagePreview}
                </DropdownToggle>
                <DropdownMenu right>
                    <DropdownItem>
                        <NavItem>
                            <NavLink to={userProfile}>Profile</NavLink>
                        </NavItem>
                    </DropdownItem>
                    <DropdownItem divider/>
                    <DropdownItem>
                        <a onClick={logout}>Logout</a>
                    </DropdownItem>
                </DropdownMenu>
            </UncontrolledDropdown>
        </React.Fragment>
    )

    const signedOutLinks = (
        <React.Fragment>
            <NavItem style={{marginRight: "20px"}}>
                <Link to='/login'>Login</Link>
            </NavItem>
            <NavItem>
                <Link to='/signup'>Signup</Link>
            </NavItem>
        </React.Fragment>
    )

    return (
        <div className="mb-3">
            <Navbar color="light" light expand="md">
                <NavbarBrand className="mx-auto" href="/">Seamflow</NavbarBrand>
                <NavbarToggler onClick={() => setIsOpen(!isOpen)}/>
                <Collapse isOpen={isOpen} navbar>
                    <Nav className="ml-auto" navbar>
                        {isAuthenticated ? signedInLinks : signedOutLinks}
                    </Nav>
                </Collapse>
            </Navbar>
        </div>
    )
}

export default Navigation