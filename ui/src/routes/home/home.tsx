import React, {ReactNode} from 'react';
import { Link } from 'react-router-dom';

const Home = (): ReactNode => (
  <main className="flex">
    <div className="column-main tile">
      <h1>kpack metrics</h1>
      <p>Simple portal with metrics for kpack</p>
      <p>
        <Link to="/projects">Show current projects</Link>
      </p>
    </div>
  </main>
);
export default Home;
