import httpApi from "../../utils/http";
import {Build, Image, Project} from "./project";

class ProjectsRepo {
    getProjects() {
        return httpApi.request('get', '/images').then(data => {
            let allProjects: Project[] = [];
            data.forEach((proj) => {
                let allImages: Image[] = [];
                if (proj['images'] != null) {
                    proj['images'].forEach((img) => {
                        let image = new Image(img['tag'], img['lastBuiltTag'], []);
                        if (img['builds'] != null) {
                            img['builds'].forEach((build) => {
                                image.builds.push(new Build(build['id'], build['reason'], build['status']))
                            });
                        }
                        allImages.push(image);
                    });
                }
                let project = new Project(proj['name'], allImages);
                allProjects.push(project)
            });
            return allProjects;
        });
    }
}

export default ProjectsRepo;