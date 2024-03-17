# Executing shell for make cmds
SHELL := /bin/bash

acceptor:
	go run ./acceptor/.

initiator:
	go run ./initiator/.

all:
	go run ./acceptor/.
	go run ./initiator/.
