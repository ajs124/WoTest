import sys
import asyncio

from wotpy.wot.servient import Servient
from wotpy.wot.wot import WoT

async def main():
    wot = WoT(servient=Servient())
    consumed_thing = await wot.consume_from_url(sys.argv[1])

    print("ConsumedThing: {}".format(consumed_thing))

if __name__ == "__main__":
    loop = asyncio.get_event_loop()
    loop.run_until_complete(main())
