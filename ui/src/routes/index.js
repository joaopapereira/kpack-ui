import React from 'react'
import { BrowserRouter as Router, Route, Switch } from 'react-router-dom'
import Home from './home'
import NotFound from './notFound'
import Projects from './projects'
class AppRouter extends React.Component {
    render() {
        return (
            <Router>
                <Switch>
                    <Route exact path='/' component={Home} />
                    <Route
                        path="/projects"
                        render={props =>
                            <Projects />
                        }
                    />
                    <Route exact component={NotFound} />
                </Switch>
            </Router>
        )
    }
}
export default AppRouter;