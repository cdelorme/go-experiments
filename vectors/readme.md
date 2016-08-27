
# vectors

I had considered programming a game with 2D vectors and in the process of following a guide I was distracted in thinking about how go implementation varied from C/C++ and what some of the performance characteristics would be.

Since go makes it too damn easy to identify these characteristics, I thought I would create a showcase of the process.


## premise

I was following a book that suggested using C floats, which I believe is a float32 whereas C's double is a float64.

The standard `math` package in go has a Sqrt function that returns a float64, so to keep things simple I started with this for the values, but I realize the overhead of extra memory consumed by float64 vs float32 given how used vectors are in games might be substantial, as well as the performance of related operations.

Then I noticed two types of overloads, the + and +=, but since golang is garbage collected it may make more sense to only support += to reduce the number of allocations and frees happening on the stack.


## comparison

Thus two layers of comparison were born:

- float32 vs float64
- add vs add-to-self

Since this comparison has layers I had to create a total of four implementations (times two or more for every additional layer).


## howto

A simple benchmark file was included with each implementation to test the performance of the four variations.

We can run these benchmarks to produce cpu profile information for `pprof` and capture memory consumption with `time`:

	/usr/bin/time --format '%Uu %Ss %er %MkB' go test -cpuprofile cpu.out -bench .

_In theory we could use the `-memprofile mem.out` as well, but I haven't had any luck in reviewing the output in `pprof`._

The cpu profile can be used to dive into which part of the code took the longest, letting us identify potential bottlenecks:

	go tool pprof -top 10 cpu.out

_The interactive mode may give us greater details if we prefer to sift through them._


## conclusion

Setting up the tooling took less than 5 minutes, which showcases the primary intention of this document; the fact is go makes gathering useful results on program execution very easy, so we can pursue trivial things like this minimal effort.

I crafted a simple script to [collect](`collect.sh`) all of the results:

	testing: warning: no tests to run
	BenchmarkVector-8   	100000000	        14.4 ns/op
	PASS
	ok  	_/home/cdelorme/tmp/vectors/float32/add	1.460s
	1.61u 0.02s 1.62r 45960kB
	1450ms of 1450ms total (  100%)
	      flat  flat%   sum%        cum   cum%
	     990ms 68.28% 68.28%      990ms 68.28%  _/home/cdelorme/tmp/vectors/float32/add.(*Vector).Length
	     460ms 31.72%   100%     1450ms   100%  _/home/cdelorme/tmp/vectors/float32/add.BenchmarkVector
	         0     0%   100%     1450ms   100%  runtime.goexit
	         0     0%   100%     1450ms   100%  testing.(*B).launch
	         0     0%   100%     1450ms   100%  testing.(*B).runN
	testing: warning: no tests to run
	BenchmarkVector-8   	100000000	        11.2 ns/op
	PASS
	ok  	_/home/cdelorme/tmp/vectors/float32/addto	1.138s
	1.29u 0.02s 1.29r 46332kB
	1130ms of 1130ms total (  100%)
	      flat  flat%   sum%        cum   cum%
	     730ms 64.60% 64.60%      730ms 64.60%  _/home/cdelorme/tmp/vectors/float32/addto.(*Vector).Length
	     400ms 35.40%   100%     1130ms   100%  _/home/cdelorme/tmp/vectors/float32/addto.BenchmarkVector
	         0     0%   100%     1130ms   100%  runtime.goexit
	         0     0%   100%     1130ms   100%  testing.(*B).launch
	         0     0%   100%     1130ms   100%  testing.(*B).runN
	testing: warning: no tests to run
	BenchmarkVector-8   	100000000	        16.8 ns/op
	PASS
	ok  	_/home/cdelorme/tmp/vectors/float64/add	1.695s
	1.85u 0.02s 1.85r 46084kB
	1680ms of 1680ms total (  100%)
	      flat  flat%   sum%        cum   cum%
	     990ms 58.93% 58.93%      990ms 58.93%  _/home/cdelorme/tmp/vectors/float64/add.(*Vector).Length
	     690ms 41.07%   100%     1680ms   100%  _/home/cdelorme/tmp/vectors/float64/add.BenchmarkVector
	         0     0%   100%     1680ms   100%  runtime.goexit
	         0     0%   100%     1680ms   100%  testing.(*B).launch
	         0     0%   100%     1680ms   100%  testing.(*B).runN
	testing: warning: no tests to run
	BenchmarkVector-8   	100000000	        12.8 ns/op
	PASS
	ok  	_/home/cdelorme/tmp/vectors/float64/addto	1.292s
	1.44u 0.03s 1.45r 46292kB
	1280ms of 1280ms total (  100%)
	      flat  flat%   sum%        cum   cum%
	     770ms 60.16% 60.16%      770ms 60.16%  _/home/cdelorme/tmp/vectors/float64/addto.(*Vector).Length
	     510ms 39.84%   100%     1280ms   100%  _/home/cdelorme/tmp/vectors/float64/addto.BenchmarkVector
	         0     0%   100%     1280ms   100%  runtime.goexit
	         0     0%   100%     1280ms   100%  testing.(*B).launch
	         0     0%   100%     1280ms   100%  testing.(*B).runN

From the test results we can clearly see that the `float32` implementation performs fast, but a more significant performance gain comes from using `addto` where we don't allocate new structures and assign the results.

The memory results collected by the `time` command tells another interesting story.  It seems that the overall consumed memory was nearly identical, but we had better values when we used `add` and were allocating new structures.

On a final note, _it's worth noting that these results are biased by my platform, a 64bit linux operating system._  We may find one implementation to be optimal, or we may find that it varies by platform.  We then have the choice of picking whichever option is the best for our target market, _or we could even use build constraints that leverages the best implementation per-platform._
