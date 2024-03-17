# Executing shell for make cmds
SHELL := /bin/bash

consumer:
	go run ./consumer/.

producer:
	go run ./producer/.

all:
	go run ./consumer/.
	go run ./producer/.
