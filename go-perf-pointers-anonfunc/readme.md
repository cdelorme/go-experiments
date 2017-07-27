
# go-perf

The original intention was to get a basic idea as to the overhead of using accessors and pointers.

After getting those results I began to add to the list to try and identify what behaviors differ, and what the cost of basic operations was.


## usage

To run simply:

	go test -run=X -bench .


## results

My results come from a Core i7 IvyBridge 4790 with 16GB DDR3 1600Mhz Corsair XMS memory, and a pair of 256GB Samsung Evo SSD's:

	$ go test -run=X -bench .
	testing: warning: no tests to run
	BenchmarkProperty-8          	2000000000	         0.35 ns/op
	BenchmarkAccessor-8          	2000000000	         0.35 ns/op
	BenchmarkPropertyPointer-8   	2000000000	         0.69 ns/op
	BenchmarkAccessorPointer-8   	2000000000	         0.69 ns/op
	BenchmarkMethod-8            	2000000000	         0.26 ns/op
	BenchmarkMethodParam-8       	2000000000	         0.26 ns/op
	BenchmarkFunc-8              	2000000000	         0.26 ns/op
	BenchmarkFuncParam-8         	2000000000	         0.26 ns/op
	BenchmarkAnonFunc-8          	2000000000	         1.55 ns/op
	BenchmarkAnonFuncParam-8     	2000000000	         1.55 ns/op
	PASS
	ok  	_/home/cdelorme/Desktop/go-perf	13.086s


## analysis

These results strongly indicate that either the compiler is performing optimzations, or that the overhead when accessing properties from functions is non-existent.  _The only overhead is then lines of code to follow suite with the idea of encapsulation._

**The benefit of this knowledge is that it becomes easier to satisfy the needs of interfaces.**

Next, the cost of pointer access is roughly double without, which makes sense as you'd hit a cache-miss on first-pass.

Next I tried benchmarking functions by package then per struct, which appear to have an identical but tiny cost, which also does not appear to vary when passing one parameter (_I assume that more parameters increases the chance of increasing the cost, but I imagine it would take a lot of parameters for it to begin to show_).

Finally, I also tested anonymous functions, which surprisingly cost nearly 6x the cost of named functions.


## conclusion

The main conclusion here is that the overhead of parameters in functions and accessor/mutator vs direct property access are very effectively non-existent.

These numbers, or the benchmark specifically, might be useful when building certain types of applications.  _For example, with a game you know the number of properties, and whether they are accessed per 16ms loop, so you can identify at what number of entities you max out resources on a system._
