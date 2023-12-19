Hot reloading with plugin system in golang.

Using [ebiten engine](https://ebitengine.org/) and [fsnotify](https://github.com/fsnotify/fsnotify) as example

How to run:

0. Make sure you are on os that support plugin (Linux/MacOS)
1. Clone this repository
3. Run `go run build.go`. this will automatically rebuild file in src/draw/draw.go
4. Run `go run src/main.go`. Main application

Good: No need restart the system.
Bad: Increase complexity.

Possible implementation:
- https://github.com/slok/reload
- https://github.com/edwingeng/hotswap
