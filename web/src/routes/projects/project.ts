export class Build {
    id: number;
    reason: string;
    status: string;

    constructor(id, reason, status) {
        this.id = id;
        this.reason = reason;
        this.status = status;
    }
}

export class Image {
    tag: string;
    lastBuiltTag: string;
    builds: Build[];

    constructor(tag, lastBuiltTag, builds) {
        this.tag = tag;
        this.lastBuiltTag = lastBuiltTag;
        this.builds = builds;
    }
}

export class Project {
    name: string;
    images: Image[];

    constructor(name, images) {
        this.name = name;
        this.images = images;
    }
}
