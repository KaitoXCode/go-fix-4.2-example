# Executing shell for make cmds
SHELL := /bin/bash

acceptor:
	go run ./acceptor/.

initiator:
	go run ./initiator/.

# Phony targets are not files
.PHONY: initiator acceptor
