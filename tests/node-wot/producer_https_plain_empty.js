Servient = require("@node-wot/core").Servient
HttpServer = require("@node-wot/binding-http").HttpServer
Helpers = require("@node-wot/core").Helpers

let servient = new Servient();
servient.addServer(new HttpServer({
    "port": "8080",
    "allowSelfSigned": true,
    "serverKey": "../rsa.key",
    "serverCert": "../rsa.crt"
}));

servient.start().then((WoT) => {
    WoT.produce({
        "@context": "https://www.w3.org/2019/wot/td/v1",
        "title": "https_plain_empty"
    }).then((thing) => {
        console.log("Produced " + thing.getThingDescription().title);

        thing.expose().then(() => {
            console.info(thing.getThingDescription().title + " ready");
            console.info("TD: " + JSON.stringify(thing.getThingDescription()));
        });
    });
});
