import json
import logging
import tornado.gen
from tornado.ioloop import IOLoop, PeriodicCallback

from wotpy.protocols.http.server import HTTPServer
from wotpy.wot.servient import Servient

CATALOGUE_PORT = 9090
HTTP_PORT = 9494

logging.basicConfig()
LOGGER = logging.getLogger()
LOGGER.setLevel(logging.INFO)

ID_THING = "urn:http_plain_empty"

DESCRIPTION = {
    "title": ID_THING,
    "id": ID_THING,
}

@tornado.gen.coroutine
def main():
    LOGGER.info("Creating HTTP server on: {}".format(HTTP_PORT))
    http_server = HTTPServer(port=HTTP_PORT)

    LOGGER.info("Creating servient with TD catalogue on: {}".format(CATALOGUE_PORT))
    servient = Servient(catalogue_port=CATALOGUE_PORT)
    servient.add_server(http_server)

    LOGGER.info("Starting servient")
    wot = yield servient.start()

    LOGGER.info("Exposing and configuring Thing")
    exposed_thing = wot.produce(json.dumps(DESCRIPTION))
    exposed_thing.expose()

if __name__ == "__main__":
    LOGGER.info("Starting loop")
    IOLoop.current().add_callback(main)
    IOLoop.current().start()
