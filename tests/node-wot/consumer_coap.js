WotCore = require("@node-wot/core")
Servient = require("@node-wot/core").Servient
CoapClientFactory = require("@node-wot/binding-coap").CoapClientFactory

let servient = new Servient();
let wotHelper = new WotCore.Helpers(servient);
servient.addClientFactory(new CoapClientFactory());

wotHelper.fetch("coap://127.0.0.1:8081/").then(async (td) => {
    let thing = await WoT.consume(td);
});
