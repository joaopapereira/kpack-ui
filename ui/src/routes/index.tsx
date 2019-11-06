import * as React from 'react';
import { ReactNode } from 'react';
import { BrowserRouter as Router, Route, Switch } from 'react-router-dom';
import Home from './home';
import NotFound from './notFound';
import Projects from './projects';

export class AppRouter extends React.Component {
  render(): ReactNode {
    return (
      <Router>
        <Switch>
          <Route exact path="/" component={Home} />
          <Route path="/projects" render={(): ReactNode => <Projects />} />
          <Route exact component={NotFound} />
        </Switch>
      </Router>
    );
  }
}
