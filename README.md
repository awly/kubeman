# kubeman - kubernetes management interface

kubeman is a [termbox](https://github.com/nsf/termbox) based tool for realtime monitoring and management of your kubernetes cluster.

It is intended to be used as a replacement for `kubectl` although it has a long way to go to match the functionality. First version of this tool will essentially be read-only.

## Installation

    go get github.com/alytvynov/kubeman

Yes, you need a working and relatively up to date version of go for this.

## But.. why?

Pretty much the difference between this:

![kubectl](kubectl.png)

And this:

![kubeman](kubeman.png)

TODO:

- status management (create, stop, resize, etc.)
- stats charts
- help page
