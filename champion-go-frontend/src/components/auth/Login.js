import React, {useState} from "react"
import {useDispatch, useSelector} from "react-redux"
import {signIn} from "../../store/modules/auth/actions/authActions"
import {Redirect} from "react-router"
import Navigation from "../Navigation"
import {Card, CardBody, CardHeader, FormGroup, Label, Button, Input} from "reactstrap"
import {Link} from "react-router-dom"

const Login = () => {
    const authState = useSelector(state => state.Auth)

    const [user, setUser] = useState({
        email: "",
        password: "",
    });

    const dispatch = useDispatch();

    const loginHandler = loginInfo => dispatch(signIn(loginInfo))

    const handleChange = e => {
        setUser({
            ...user,
            [e.target.name]: e.target.value,
        })
    }

    const submitUser = e => {
        e.preventDefault()
        loginHandler({
            email: user.email,
            password: user.password,
        })
    }

    if (authState.isAuthenticated) {
        return <Redirect to="/"/>
    }

    return (
        <div className="app">
            <div>
                <Navigation/>
            </div>
            <div className="container auth">
                <Card className="card-style">
                    <CardHeader>Login</CardHeader>
                    <CardBody>
                        <form onSubmit={submitUser}>
                            <div className="mb-2">
                                {authState.loginError && authState.loginError.Incorrect_details ? (
                                    <small className="color-red">{authState.loginError.Incorrect_details}</small>
                                ) : (
                                    ""
                                )}
                                {authState.loginError && authState.loginError.No_record ? (
                                    <small className="color-red">{authState.loginError.No_record}</small>
                                ) : (
                                    ""
                                )}
                            </div>

                            <FormGroup>
                                <Label>Email</Label>
                                <Input
                                    type="text"
                                    name="email"
                                    onChange={handleChange}
                                    placeholder="Enter email"
                                />
                                {authState.loginError && authState.loginError.Required_email ? (
                                    <small className="color-red">{authState.loginError.Required_email}</small>
                                ) : (
                                    ""
                                )}
                                {authState.loginError && authState.loginError.Invalid_email ? (
                                    <small className="color-red">{authState.loginError.Invalid_email}</small>
                                ) : (
                                    ""
                                )}
                            </FormGroup>

                            <FormGroup>
                                <Label>Password</Label>
                                <Input
                                    type="text"
                                    name="password"
                                    onChange={handleChange}
                                    placeholder="Enter password"
                                />
                                {authState.loginError && authState.loginError.Required_password ? (
                                    <small className="color-red">{authState.loginError.Required_password}</small>
                                ) : (
                                    ""
                                )}
                                {authState.loginError && authState.loginError.Invalid_password ? (
                                    <small className="color-red">{authState.loginError.Invalid_password}</small>
                                ) : (
                                    ""
                                )}
                                {authState.loginError && authState.loginError.Incorrect_password ? (
                                    <small className="color-red">{authState.loginError.Incorrect_password}</small>
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
                                    Login...
                                </Button>
                            ) : (
                                <Button
                                    color="primary"
                                    type="submit"
                                    block
                                    disabled={user.email === "" || user.password === ""}
                                >
                                    Login
                                </Button>
                            )}

                        </form>

                        <div className="mt-2" style={{display: "flex", justifyContent: "space-between"}}>
                            <div>
                                <small><Link to="/signup">Sign Up</Link></small>
                            </div>
                            <div>
                                <small><Link to="/forgotpassword">Forgot Password?</Link></small>
                            </div>
                        </div>

                    </CardBody>
                </Card>
            </div>
        </div>
    )
}

export default Login