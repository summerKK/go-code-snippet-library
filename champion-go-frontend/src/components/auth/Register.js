import {useDispatch, useSelector} from "react-redux"
import React, {useState} from "react"
import {Redirect} from "react-router"
import Navigation from "../Navigation"
import {Card, CardHeader, CardBody, FormGroup, Input, Label, Button} from "reactstrap"
import {Link} from "react-router-dom"
import {signUp} from "../../store/modules/auth/actions/authActions"
import "./Index.css"

const Register = () => {
    const authState = useSelector(state => state.Auth)
    const [user, setUser] = useState({
        username: "",
        email: "",
        password: "",
    })

    const dispatch = useDispatch()

    const registerHandle = credentials => dispatch(signUp(credentials))

    const handleChange = e => {
        setUser({
            ...user,
            [e.target.name]: e.target.value,
        })
    }

    const submitUser = e => {
        e.preventDefault()
        registerHandle({
            username: user.username,
            email: user.email,
            password: user.password,
        })
    }

    if (authState.isAuthenticated) {
        return <Redirect to='/'/>
    }

    return (
        <div className="app">
            <div>
                <Navigation/>
            </div>
            <div className="container auth">
                <Card className="card-style">
                    <CardHeader>Welcome To SeamFlow</CardHeader>
                    <CardBody>
                        <form onSubmit={submitUser}>
                            <FormGroup>
                                <Label>User Name</Label>
                                <Input
                                    type="text"
                                    name="username"
                                    placeholder="Enter Username"
                                    onChange={handleChange}
                                />
                                {authState.signupError && authState.signupError.Required_username ? (
                                    <small className="color-red">{authState.signupError.Required_username}</small>
                                ) : (
                                    ""
                                )}
                                {authState.signupError && authState.signupError.Taken_username ? (
                                    <small className="color-red">{authState.signupError.Taken_username}</small>
                                ) : (
                                    ""
                                )}
                            </FormGroup>

                            <FormGroup>
                                <Label>Email</Label>
                                <Input type="email" name="email" placeholder="Enter email" onChange={handleChange}/>
                                {authState.signupError && authState.signupError.Required_email ? (
                                    <small className="color-red">{authState.signupError.Required_email}</small>
                                ) : (
                                    ""
                                )}
                                {authState.signupError && authState.signupError.Invalid_email ? (
                                    <small className="color-red">{authState.signupError.Invalid_email}</small>
                                ) : (
                                    ""
                                )}
                                {authState.signupError && authState.signupError.Taken_email ? (
                                    <small className="color-red">{authState.signupError.Taken_email}</small>
                                ) : (
                                    ""
                                )}
                            </FormGroup>

                            <FormGroup>
                                <Label>Password</Label>
                                <Input type="password" name="password" placeholder="Enter password"
                                       onChange={handleChange}/>
                                {authState.signupError && authState.signupError.Required_password ? (
                                    <small className="color-red">{authState.signupError.Required_password}</small>
                                ) : (
                                    ""
                                )}
                                {authState.signupError && authState.signupError.Invalid_password ? (
                                    <small className="color-red">{authState.signupError.Invalid_password}</small>
                                ) : (
                                    ""
                                )}
                            </FormGroup>

                            {authState.isLoading ? (
                                <Button
                                    color="primary"
                                    type="submit"
                                    block
                                    disabled
                                >
                                    Registering...
                                </Button>
                            ) : (
                                <Button
                                    color="primary"
                                    type="submit"
                                    block
                                    disabled={user.username === "" || user.email === "" || user.password === ""}
                                >
                                    Register
                                </Button>
                            )}
                        </form>
                        <div className="mt-2">
                            <small>Have an account? <Link to="/login">Please login</Link></small>
                        </div>
                    </CardBody>
                </Card>
            </div>
        </div>
    )

}

export default Register