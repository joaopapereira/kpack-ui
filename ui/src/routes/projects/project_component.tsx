import * as React from "react";
import ProjectsRepo from "./repository";
import {Project} from "./project";
import './style.scss'


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

class Projects extends React.Component<ProjectsProps, ProjectsState> {

    constructor(props) {
        super(props);
        this.state = {
            listOfProjects: []
        };
        this.reloadProjects.bind(this);
    }

    componentDidMount() {
        this.reloadProjects();
    }

    reloadProjects() {
        this.props.projectRepo.getProjects().then((listOfProjects) => {
            this.setState({listOfProjects: listOfProjects});
        });
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

class ProjectsProps {
    projectRepo: ProjectsRepo
}

class ProjectsState {
    listOfProjects: Project[]
}

export default Projects;