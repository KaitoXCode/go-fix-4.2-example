# bi-dir-fix-example

project including a fix producer and consumer which communicate simple messages
to eachother

## deployment

manual -> start acceptor: `go run ./consumer/.` start initiator:
`go run ./producer/.` \n make -> start acceptor: `make consumer` start
initiator: `make producer` \n make -> start all: `make all`
