# mitohg

## Description

Определение гаплогрупп происходит на основе номенклатуры [PhyloTree build 17](http://www.phylotree.org) относительно референсной последовательности [RSRS](http://www.phylotree.org/resources/RSRS.fasta).

## Build

```
cd wfa_bridge
gcc -O2 -I/path/to/WFA -c wfa_bridge.c
ar q libwfabridge.a wfa_bridge.o

cd ..
go build -o mitohg
```
