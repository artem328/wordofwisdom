# Word Of Wisdom Server (Proof Of Work Test)

## Requirements
Design and implement “Word of Wisdom” tcp server.
- TCP server should be protected from DDOS attacks with the Proof of Work (https://en.wikipedia.org/wiki/Proof_of_work), 
the challenge-response protocol should be used.
- The choice of the POW algorithm should be explained.
- After Proof Of Work verification, server should send one of the quotes from “word of wisdom” book 
or any other collection of the quotes.
- Docker file should be provided both for the server and for the client that solves the POW challenge

## Run

First, build the docker images of the server and client with the following command
```shell
docker build -f Dockerfile.server -t wow-server .; \
docker build -f Dockerfile.client -t wow-client .
```

Create docker network for the server-client communication
```shell
docker network create wow-net
```

Run the server with default options or modify command options as per your needs 
(run the server with `-h` flag to see available options). 

```shell
docker run --rm --net wow-net --name wow-server wow-server -debug
```

When server is launched, start the client in other terminal (use `-h` flag to see available options of the client)
```shell
docker run --rm --net wow-net wow-client -addr wow-server:9000 -verbose
```

Alternatively, run the docker compose setup
```shell
docker compose up
```

## PoW Algorithm

To perform proof-of-work the interactive Hashcash algorithm is used. It has unbounded probabilistic cost, 
which theoretically can take forever to solve the challenge, however, the more time is spend for solving the challenge,
there are more chances to solve it.

The algorithm flow is the following
1. Client connects to the server 
2. Server sends a 4 bytes challenge (_c_) and target difficulty (_d_).
3. Client finds a 4 bytes solution (_s_), such that sha256(_c_ & _s_) produces hash sum with _d_ leading zero bits
(& means concatenation)
4. Client sends the solution to the server
5. Server verifies that challenge concatenated with solution indeed produce hash with target leading zero bits
6. In case verification successful, the server waits for the request, otherwise closes connection

Hashcash is chosen since
- it has configurable difficulty, which may reduce or increase computation time of the client dynamically 
depending on different conditions, such as server load, or number of connections from IP/subnet
- the verification process takes constant time and significantly faster than solving a challenge
- the algorithm uses sha256, which is available in most programming languages, thus the algorithm is easy to implement
- the algorithm is CPU bound
- it is nearly impossible to store the solutions for all challenges in memory
- it is not impossible to find a solution depending on previous challenge and solution
