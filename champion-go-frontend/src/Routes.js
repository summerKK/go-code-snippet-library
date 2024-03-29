import React from "react"
import {Route, Router, Switch} from "react-router"
import {history} from "./history"
import Dashboard from "./components/Dashboard"
import Register from "./components/auth/Register"
import Login from "./components/auth/Login"
import Profile from "./components/users/Profile"

const Routes = () => {
    return (
        <Router history={history}>
            <div className="App">
                <Switch>
                    <Route exact path="/" component={Dashboard}/>
                    <Route path="/signup" component={Register}/>
                    <Route path="/login" component={Login}/>
                    <Route path="/profile" component={Profile}/>
                </Switch>
            </div>
        </Router>
    )
}

export default Routes