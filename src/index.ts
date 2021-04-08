import * as ws from "ws";
/*
const server = new ws.Server({
    port: +process.env.PORT
});

server.on('connection', function connection(ws) {
    ws.send('Hello, world3');
    ws.close();
});
*/

type Video = {
    id: number;
    title: string;
    rating: number;
    description: string;
};

const video: Video = {
    id: 69,
    title: "Foo Bar",
    rating: 7,
    description: "This is about foo bar",
}

console.log(JSON.stringify(video, null, 4));


class VideoImpl {
    constructor(private data: (string | number)[], private offset: number) { }

    getTitle() {
        this.data[this.offset + 1];
    }

    getId() {
        this.data[this.offset + 0];
    }

    getRating() {
        this.data[this.offset + 2];
    }

    getDescription() {
        this.data[this.offset + 3];
    }

    static toArray(video: Video): (string | number)[] {
        return [
            video.id,
            video.title,
            video.rating,
            video.description,
        ];
    }
}

