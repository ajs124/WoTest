Servient = require("@node-wot/core").Servient
HttpServer = require("@node-wot/binding-http").HttpServer
Helpers = require("@node-wot/core").Helpers

let servient = new Servient();
servient.addServer(new HttpServer({
    "port": "8080",
    "security": {
        "scheme": "basic"
    }
}));

servient.addCredentials({
    "urn:uuid:4e9d9ff9-02ee-444c-805d-6618f4151e60": {
        "username": "Ooz8aideeMuuRahwae3u",
        "password": "ih4ieZ7Pe7wigi8aengu"
    }
});

servient.start().then((WoT) => {
    WoT.produce({
        "@context": "https://www.w3.org/2019/wot/td/v1",
        "title": "http_basic_empty",
        "id": "urn:uuid:4e9d9ff9-02ee-444c-805d-6618f4151e60"
    }).then((thing) => {
        console.log("Produced " + thing.getThingDescription().title);

        thing.expose().then(() => {
            console.info(thing.getThingDescription().title + " ready");
            console.info("TD: " + JSON.stringify(thing.getThingDescription()));
        });
    });
});
