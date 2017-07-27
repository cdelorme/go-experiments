
# separating methods using go build directives

This is a trivial demonstration of how you can use build directives to deal with multiple architectures without writing architecture-aware code.

A build directive in go is the separation of files with a suffix matching the `GOOS`, or an explicit comment at the start of a file such as `// +build linux darwin`.

**More importantly, this demonstration shows that you can define functions with structure receivers defined in other files within the same package.**

This sort of behavior gives you an enormous amount of flexibility.
