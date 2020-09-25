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
        "title": "https_largeProp",
        "properties": {
            "bench": {
                "title": "download this",
                "type": "string"
            }
        }
    }).then((thing) => {
        console.log("Produced " + thing.getThingDescription().title);
        let goalSize = parseInt(process.argv[2]);
        let buf = Buffer.alloc(goalSize);
        for (i=0; i < goalSize; ++i) {
            let r = Math.floor(Math.random() * Math.floor(255));
            buf.write(r.toString(16), i);
        }
        console.log("Generated " + buf.length + " character long random string")
        thing.setPropertyReadHandler("bench", (options) => {
            return new Promise ((resolve, reject) => {
                resolve(buf.toString());
            });
        });

        thing.expose().then(() => {
            console.info(thing.getThingDescription().title + " ready");
            console.info("TD: " + JSON.stringify(thing.getThingDescription()));
        });
    });
});
