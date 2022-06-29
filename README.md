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
