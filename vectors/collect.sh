#!/bin/bash

cd float32/add
/usr/bin/time --format '%Uu %Ss %er %MkB' go test -cpuprofile cpu.out -bench .
go tool pprof -top 10 cpu.out

cd ../addto
/usr/bin/time --format '%Uu %Ss %er %MkB' go test -cpuprofile cpu.out -bench .
go tool pprof -top 10 cpu.out

cd ../../float64/add
/usr/bin/time --format '%Uu %Ss %er %MkB' go test -cpuprofile cpu.out -bench .
go tool pprof -top 10 cpu.out

cd ../addto
/usr/bin/time --format '%Uu %Ss %er %MkB' go test -cpuprofile cpu.out -bench .
go tool pprof -top 10 cpu.out
