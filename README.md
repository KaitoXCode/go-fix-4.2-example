# bi-dir-fix-example

project including a fix initiator and acceptor which communicate a simple
execution report to eachother every 10s

## deployment

manual -> start acceptor: `go run ./acceptor/.` start initiator:
`go run ./initiator/.` \n make -> start acceptor: `make acceptor` start
initiator: `make initiator`
