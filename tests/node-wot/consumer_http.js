WotCore = require("@node-wot/core")
Servient = require("@node-wot/core").Servient
HttpClientFactory = require("@node-wot/binding-http").HttpClientFactory

let servient = new Servient();
let wotHelper = new WotCore.Helpers(servient);
servient.addClientFactory(new HttpClientFactory());

wotHelper.fetch("http://127.0.0.1:8081/").then(async (td) => {
    let thing = await WoT.consume(td);
});
