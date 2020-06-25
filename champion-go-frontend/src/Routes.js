import React from "react"
import {Route, Router, Switch} from "react-router"
import {history} from "./history"
import Dashboard from "./components/Dashboard"
import Register from "./components/auth/Register";

const Routes = () => {
    return (
        <Router history={history}>
            <div className="App">
                <Switch>
                    <Route exact path='/' component={Dashboard}/>
                    <Route path='/signup' component={Register}/>
                </Switch>
            </div>
        </Router>
    )
}

export default Routes