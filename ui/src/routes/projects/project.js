import ProjectsRepo from "./repository";
import React from 'react'
import './style.scss'

export class Build {
    constructor(id, reason, status) {
        this.id = id;
        this.reason = reason;
        this.status = status;
    }
}

export class Image {
    constructor(tag, lastBuiltTag, builds) {
        this.tag = tag;
        this.lastBuiltTag = lastBuiltTag;
        this.builds = builds;
    }
}

export class Project {
    constructor(name, images) {
        this.name = name;
        this.images = images;
    }
}

function boxColor(build) {
    switch (build.status.toLowerCase()) {
        case 'succeeded':
            return "box-success";
        case 'failed':
            return "box-failure";
        case 'building':
            return "box-building";
        default:
            return "box-building";
    }
}

const BuildView = (props) => (
    <div className={boxColor(props.build)}>
        {props.build.id}
    </div>
);

const ImageView = (props) => (
    <div className="image">
        <div className="name">{props.image.tag}</div>
        <h6>Latest tag: {props.image.lastBuiltTag}</h6>
        <section className="boxes">
            {props.image.builds.map((build) => (
                <BuildView build={build}/>
            ))}
        </section>
    </div>
);

const ProjectView = (props) => (
    <div className="project">
        <div className="title">Project: {props.project.name}</div>
        <section className="images">
            {props.project.images.map((image) => (
                <ImageView image={image}/>
            ))}
        </section>
    </div>
);

class Projects extends React.Component {
    constructor(props) {
        super(props);
        this.repo = new ProjectsRepo();
        this.state = {
            listOfProjects: []
        };
        this.reloadProjects.bind(this);
    }

    componentDidMount() {
        this.reloadProjects();
    }

    reloadProjects() {
        this.repo.getProjects().then((listOfProjects) => {
            this.setState({listOfProjects: listOfProjects});
        });
    }

    componentWillUnmount() {
        clearInterval(this.interval);
    }

    render() {
        return <div>
            <button onClick={this.reloadProjects.bind(this)}>Refresh page</button>
            <section className="projects">
                {this.state.listOfProjects.map(project => (
                    <ProjectView project={project}/>
                ))}
            </section>
        </div>;
    }
}

export default Projects