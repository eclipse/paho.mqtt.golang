Docker example
==============

This folder contains all that is needed to build an environment with a publisher, broker (mosquitto) and subscriber
using docker (ideally `docker-compose`). While it provides an end-to-end example its primary purpose is to act as a
starting point for producing reproducible examples (when logging an issue with the library).

Because the publisher (`pub`), broker (`mosquitto`) and subscriber (`sub`) run in separate containers this setup closely
simulates a real deployment. One thing to bear in mind is that the network between the containers is very fast and
reliable (but there are some techniques that can be used to simulate failures etc).

# Usage

Ensure that you have [docker](https://docs.docker.com/get-docker/) and
[docker-compose](https://docs.docker.com/compose/install/) installed.

To start everything up change into the `cmd/docker` folder and run:

```
docker-compose up --build --detach
```

This will start everything up in the background. You can see what is happening by running:

```
docker-compose logs --follow
```

This will display a lot of information (mosquitto is running with debug level logging). To see the subscriber logs:

```
docker-compose logs --follow sub
```

Note: Messages received by the subscriber will be written to `shared/receivedMessages` (you may want to delete the
contents of this file from time to time!).

To stop everything run:

```
docker-compose down
```

Feel free to copy the folder and modify the publisher/subscriber to work as you want them to!

Note: The `pub` and `sub` containers connect to mosquitto via the internal network (`test-net`) but mosquitto should
also be available on the host port `8883` if you wish to connect to it. This will not work if you have mosquitto
installed locally (edit the `docker-compose.yml` and change the `published` port).

# Simulating Network Connection Loss

You can simulate the loss of network connectivity by disconnecting the network adapter within a container. e.g.

```
docker network disconnect lostpackets_test-net lostpackets_pub_1
docker network connect lostpackets_test-net lostpackets_pub_1
```
  