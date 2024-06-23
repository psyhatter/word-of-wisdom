# Word of Wisdom TCP Server

The Word of Wisdom TCP Server is a server that provides quotes after completing a Proof of Work challenge. It is
designed to protect against DDoS attacks using a challenge-response protocol. Hashcash was chosen as the POW algorithm
because it is the most well-known. The TCP server has been dockerized for ease of deployment.

## Task

Design and implement “Word of Wisdom” tcp server.
- TCP server should be protected from DDOS attacks with the [Proof of Work](https://en.wikipedia.org/wiki/Proof_of_work),
the challenge-response protocol should be used.
- The choice of the POW algorithm should be explained.
- After Proof Of Work verification, server should send one of the quotes from “word of wisdom” book or any other
collection of the quotes.
- Docker file should be provided both for the server and for the client that solves the POW challenge

## Requirements:

- Go 1.22+ installed (to run tests, start server or client without Docker)
- Docker installed (to run docker-compose)

## Getting started

```bash
git clone https://github.com/psyhatter/word-of-wisdom.git
cd word-of-wisdom
docker compose -f docker-compose.yml up
```
