import * as React from 'react';
import { ReactNode } from 'react';
import { BrowserRouter as Router, Route, Switch } from 'react-router-dom';
import Home from './home';
import NotFound from './notFound';
import Projects from './projects';
import ProjectsRepo from './projects/repository';

const projectsRepo = new ProjectsRepo();

export class AppRouter extends React.Component {
  render(): ReactNode {
    return (
      <Router>
        <Switch>
          <Route exact path="/" component={Home} />
          <Route
            path="/projects"
            render={(): ReactNode => <Projects projectRepo={projectsRepo} />}
          />
          <Route exact component={NotFound} />
        </Switch>
      </Router>
    );
  }
}
