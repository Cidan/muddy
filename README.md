# muddy
Muddy is a side project of mine and an evolution of my previous Go based MUD,
gomud. While gomud was clever, it was too clever for it's own good and had
too many antipatterns for Go, specifically around locks and the interpreter.
Muddy will instead focus on correct message passing between goroutines via
channels to avoid locks when updating/writing objects, a supervisor based
goroutine manager using the wonderful Suture library, and a more simplified
interpreter system.

Much of the code that comprises muddy will come from gomud, with most of the
low level construct handling being refactored.


## Why?

I use this project both to work on a passion of mine, text based games, and
to further my skills as a developer. This project ends up being where I try
a lot of different things/expand my software architecture skills as a result.


## Overall Archtecture

The idea mehind muddy is to emulate an erlang-like actor system, but much
more lightweight and flexible, while maintaining a Go-like programming style.

Each object in the world is an actor running under a supervisor, which includes
players, rooms, objects, etc. These objects can only be modified by messages
passed into the object via a channel. The channel implementation is hidden
from the caller, and just appears as normal methods on the object type.

This ensures that only a single goroutine is responsible for updating an objects
properties, simplifying the issue of complex object interaction and reentrant
locks, which Go does not support. Reading a property of an object (i.e. player
health, etc) is a simple read of the variable that defines it, wrapped in a read
lock. This is also hidden from the called in various helper functions.

Player input is handled by the player's actor, and is passed to an interpreter
via a channel. This interpreter acts as a nexus for player commands. For now,
all commands will be interpreted in the same goroutine, effectively making all
commands across all players serial in nature. In the future, this may be done
in parallel across users once a strategy for dealing with deadlocks has been
worked through. 