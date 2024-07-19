# Multithreaded TCP Server written in Go

- one thread (goroutine in this case) per client model
- spawns a goroutine thread for every connecting client
- server close down gracefully handled
