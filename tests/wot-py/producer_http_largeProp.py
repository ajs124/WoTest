import json
import sys
import logging
import random
import math
import tornado
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
    "properties": {
        "bench": {
            "title": "download this",
            "type": "string"
        }
    }
}

buf = ""

async def bench_handler():
    global buf
    return buf

@tornado.gen.coroutine
def main():
    global buf
    LOGGER.info("Creating HTTP server on: {}".format(HTTP_PORT))
    http_server = HTTPServer(port=HTTP_PORT)

    LOGGER.info("Creating servient with TD catalogue on: {}".format(CATALOGUE_PORT))
    servient = Servient(catalogue_port=CATALOGUE_PORT)
    servient.add_server(http_server)

    while len(buf) < int(sys.argv[1]):
        r = random.random()
        c = math.floor(r*256)
        buf += format(c, 'x')

    LOGGER.info("Starting servient")
    wot = yield servient.start()

    LOGGER.info("Exposing and configuring Thing")
    exposed_thing = wot.produce(json.dumps(DESCRIPTION))
    exposed_thing.set_property_read_handler("bench", bench_handler)
    exposed_thing.expose()

if __name__ == "__main__":
    LOGGER.info("Starting loop")
    IOLoop.current().add_callback(main)
    IOLoop.current().start()
