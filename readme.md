A simple steganographic algorithm, exposed in library routines, a
command-line interface, as well as a web proxy demonstration
implementation.

# [GoDoc Documentation](https://godoc.org/chrispennello.com/go/steg)

# I/O Throughput
    % sysctl hw.model hw.machine hw.ncpu
    hw.model: AMD Phenom(tm) II X4 955 Processor
    hw.machine: amd64
    hw.ncpu: 4

The following test was performed.  Each benchmark run represents a test
of muxing an appropriate number of message bytes into three million
carrier bytes.

    % go test -bench . -benchtime 60s
    PASS
    Benchmark1           100         683815634 ns/op
    Benchmark2            50        1501107422 ns/op
    Benchmark3            50        1738694850 ns/op

This yields the following throughput statistics.

    1 byte  per atom    4.387MB/s
    2 bytes per atom    1.999MB/s
    3 bytes per atom    1.725MB/s

# Installation
    go get chrispennello.com/go/steg

The command-line interface:

    go get chrispennello.com/go/steg/cmd/steg

The web proxy demo:

    go get chrispennello.com/go/steg/cmd/stegserve

`stegserve` can also be run as an App Engine application.
