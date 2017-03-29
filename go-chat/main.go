package main

// A trivial and unscientific benchmark test case modeled around UDP
// communication patterns, namely to compare mutex versus channel
// performance when attempting to distribute messages to many connections
// while measuring the performance characteristics as the number of
// connections increases

// I have two possible conclusions from the benchmarks:
// 1. Either the performance of channels outperforms mutexes at greater scale
// 2. Or the operational overhead of distributing messages outweighs the costs
// of mutex or channel distribution.
// On my local system, once we use 20 connections the performance is near equal

func main() {}
