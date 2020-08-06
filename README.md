# WoTest
Testing framework for WoT implementations, namely:
- thingweb node-wot
- wot-py
- sane-city wot-servient

# Features
The following "protocols" are tested for the Consume and Producer:
- HTTP
- CoAP
- MQTT
- Websockets

For each of the protocols the following operations can be tested:
- readproperty
- writeproperty
- observeproperty
- unobserveproperty
- readallproperties
- writeallproperties
- readmultipleproperties
- writemultipleproperties
- invokeaction
- subscribeevent
- unsubscribeevent

With the following security schemes:
- none
- basic
- digest
- bearer
- apikey
- psk
- oauth2

# How to run
## Regular Linux distributions
Use `run.sh`.
It tries to do some best effor checks to see if you have all the relevant languages and runtimes installed and warns you if you do not.
How to install dependencies is up to you.
## NixOS
If `run.sh` detects that it's running on NixOS it will automatically download all the dependencies for you.
## Other operating systems (BSD, NT, darwin, ...)
You're on your own.
If the OS is easy enough to test on feel free to open a PR.
This excludes darwin, because of how it is distributed/sold.
