import React from "react"
import {Route, Router, Switch} from "react-router"
import {history} from "./history"
import Dashboard from "./components/Dashboard"

const Routes = () => {
    return (
        <Router history={history}>
            <div className="App">
                <Switch>
                    <Route exact path='/' component={Dashboard}/>
                </Switch>
            </div>
        </Router>
    )
}

export default Routes