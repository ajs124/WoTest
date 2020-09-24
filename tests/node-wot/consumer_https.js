WotCore = require("@node-wot/core")
Servient = require("@node-wot/core").Servient
HttpsClientFactory = require("@node-wot/binding-http").HttpsClientFactory

let servient = new Servient();
let wotHelper = new WotCore.Helpers(servient);
servient.addClientFactory(new HttpsClientFactory({
    "config": {
        "allowSelfSigned": true
    }
}));

wotHelper.fetch("https://127.0.0.1:8081/").then(async (td) => {
    let w = await servient.start();
    let thing = await w.consume(td);
});
