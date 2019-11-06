import * as React from 'react';
import Projects from "../../routes/projects";
import ProjectsRepo from "../../routes/projects/repository";
import {shallow} from "enzyme";

describe('Display projects', () => {
    let projectRepo: ProjectsRepo
    beforeEach(() => {
        projectRepo = new ProjectsRepo();
    });

    it('no projects are displayed when cannot retrieve any', () => {
        projectRepo.getProjects = jest.fn(() => Promise.resolve([]));
        let subject = shallow(<Projects projectRepo={projectRepo}/>);
        expect(subject.text()).toContain('Refresh page');
        expect(projectRepo.getProjects).toHaveBeenCalledTimes(1);
    });

    it('calls the backend when refresh page button is pressed', () => {
        projectRepo.getProjects = jest.fn(() => Promise.resolve([]));
        let subject = shallow(<Projects projectRepo={projectRepo}/>);
        expect(subject.text()).toContain('Refresh page');
        expect(subject.find('button').at(0).text()).toEqual('Refresh page');
        subject.find('button').simulate('click');
        expect(projectRepo.getProjects).toHaveBeenCalledTimes(2);
    });
});