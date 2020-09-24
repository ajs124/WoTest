import json
import logging
import tornado.gen
from tornado.ioloop import IOLoop, PeriodicCallback

from wotpy.protocols.coap.server import CoAPServer
from wotpy.wot.servient import Servient

CATALOGUE_PORT = 9090
COAP_PORT = 9393

logging.basicConfig()
LOGGER = logging.getLogger()
LOGGER.setLevel(logging.INFO)

ID_THING = "urn:coap_plain_empty"

DESCRIPTION = {
    "title": ID_THING,
    "id": ID_THING,
}

@tornado.gen.coroutine
def main():
    LOGGER.info("Creating COAP server on: {}".format(COAP_PORT))
    coap_server = CoAPServer(port=COAP_PORT)

    LOGGER.info("Creating servient with TD catalogue on: {}".format(CATALOGUE_PORT))
    servient = Servient(catalogue_port=CATALOGUE_PORT)
    servient.add_server(coap_server)

    LOGGER.info("Starting servient")
    wot = yield servient.start()

    LOGGER.info("Exposing and configuring Thing")
    exposed_thing = wot.produce(json.dumps(DESCRIPTION))
    exposed_thing.expose()

if __name__ == "__main__":
    LOGGER.info("Starting loop")
    IOLoop.current().add_callback(main)
    IOLoop.current().start()
