WotCore = require("@node-wot/core")
Servient = require("@node-wot/core").Servient
CoapClientFactory = require("@node-wot/binding-coap").CoapClientFactory

let servient = new Servient();
let wotHelper = new WotCore.Helpers(servient);
servient.addClientFactory(new CoapClientFactory());

wotHelper.fetch("coap://127.0.0.1:5683/").then(async (td) => {
    let w = await servient.start();
    let thing = await w.consume(td);
});
