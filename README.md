# WoTest
A black box testing framework for WoT implementations, namely:
- [thingweb node-wot](https://github.com/eclipse/thingweb.node-wot)
- [wot-py](https://github.com/agmangas/wot-py)
- [sane-city wot-servient](https://github.com/sane-city/wot-servient)

# Features
The following protocols are tested for the Consumer and Producer:
- [x] HTTP
- [x] CoAP (UDP only)
- [ ] MQTT

For each of the protocols the following operations can be tested:
- [ ] readproperty
- [ ] writeproperty
- [ ] observeproperty
- [ ] unobserveproperty
- [ ] readallproperties
- [ ] writeallproperties
- [ ] readmultipleproperties
- [ ] writemultipleproperties
- [ ] invokeaction
- [ ] subscribeevent
- [ ] unsubscribeevent

With the following security schemes:
| Scheme | HTTP | CoAP | MQTT |
|--------|:----:|:----:|:----:|
| none   |      |      |      |
| basic  |      |      |      |
| digest |      |      |      |
| bearer |      |      |      |
| apikey |      |      |      |
| psk    |      |      |      |
| oauth2 |      |      |      |

There are also some quantitative tests:
- How many connections per second can be handled including latency percentiles
- How much throughput can be reached on a single connection (e.g. read and writeProperty with large data)
- How big of a property can be read and written without error (up to a maximum of ? MB/GB)

# How to run
## Regular Linux distributions
Use `run.sh`.
It tries to do some best effor checks to see if you have all the relevant languages and runtimes installed and warns you if you do not.
How to install dependencies is up to you.
## With docker installed
If `run.sh` finds a docker binary, it will use the included Dockerfiles, which should pull all dependencies for you.
## NixOS
If `run.sh` detects that it's running on NixOS it will automatically download all the dependencies for you.
## Other operating systems (BSD, NT, darwin, ...)
You're on your own.
If the OS is easy enough to test on feel free to open a PR.
This excludes darwin, because of how it is distributed/sold.

# Notes
## Protocol implementations
The following client and server implementations are used:
- http: golang standard library
- coap: https://github.com/plgd-dev/go-coap
- mqtt: https://github.com/eclipse/paho.mqtt.golang
## Design decisions
### Programming language
You might be wondering why this is implemented in go, if the implementations being tested are in node.js, java and python.
The reason is simple, this way the risk of reusing code from them and therefore being biased is mitigated.
The protocol implementations are also likely do be different, so the likelyhood of being bug-compatible to other implementations is lower.
### Overall
The idea was to not write another WoT implementation, but to have a blackbox testing framework which operates in given constraints and executes predefined tests.
This is to be able to find bugs, strengths and weaknesses of existing WoT implementations and give them an opportunity to validate changes e.g. in a CI pipeline.
