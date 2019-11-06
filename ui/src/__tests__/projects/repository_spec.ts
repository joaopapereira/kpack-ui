import ProjectsRepo from "../../routes/projects/repository";
import axios from 'axios';
import MockAdapter from "axios-mock-adapter";
import {Build, Project} from "../../routes/projects/project";

describe('getProjects', () => {
    let subject: ProjectsRepo;
    let axiosMock = new MockAdapter(axios);

    beforeEach(() => {
        subject = new ProjectsRepo();
    });

    afterEach(() => {
        axiosMock.reset()
    });

    it('returns empty list when no projects are fetched', async () => {
        axiosMock.onGet('/images').reply(200, []);

        await subject.getProjects().then((projects: Project[]) => {
            expect(projects).toHaveLength(0);
        })
    });

    it('returns single project with images when no images exist for a project', async () => {
        axiosMock.onGet('/images').reply(200, [{
            name: 'some-name'
        }]);

        await subject.getProjects().then((projects: Project[]) => {
            expect(projects).toHaveLength(1);
            expect(projects[0].name).toEqual('some-name');
            expect(projects[0].images).toEqual([]);
        })
    });

    it('returns project with an image an no builds', async () => {
        axiosMock.onGet('/images').reply(200, [{
            name: 'some-name',
            images: [
                {
                    tag: 'some-tag',
                    lastBuiltTag: 'last-built-tag'
                }
            ]
        }]);

        await subject.getProjects().then((projects: Project[]) => {
            expect(projects).toHaveLength(1);
            expect(projects[0].name).toEqual('some-name');
            expect(projects[0].images).toHaveLength(1);
            expect(projects[0].images[0].tag).toEqual('some-tag');
            expect(projects[0].images[0].lastBuiltTag).toEqual('last-built-tag');
            expect(projects[0].images[0].builds).toEqual([]);
        })
    });

    it('returns project with an image and builds', async () => {
        axiosMock.onGet('/images').reply(200, [{
            name: 'some-name',
            images: [
                {
                    tag: 'some-tag',
                    lastBuiltTag: 'last-built-tag',
                    builds: [
                        {
                            id: '1',
                            reason: 'some-reason',
                            status: 'some-status'
                        }
                    ]
                }
            ]
        }]);

        await subject.getProjects().then((projects: Project[]) => {
            expect(projects).toHaveLength(1);
            expect(projects[0].name).toEqual('some-name');
            expect(projects[0].images).toHaveLength(1);
            expect(projects[0].images[0].tag).toEqual('some-tag');
            expect(projects[0].images[0].lastBuiltTag).toEqual('last-built-tag');
            expect(projects[0].images[0].builds).toEqual([new Build('1', 'some-reason', 'some-status')]);
        })
    });
});