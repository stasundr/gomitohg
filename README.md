# mitohg

## Description

Определение гаплогрупп происходит на основе номенклатуры [PhyloTree build 17](http://www.phylotree.org) относительно референсной последовательности [RSRS](http://www.phylotree.org/resources/RSRS.fasta). Для выравнивания митогеномов относительно RSRS используется библиотека [WFA](https://github.com/smarco/WFA).

## Build

See `Dockerfile` for build details.

```
docker build -t mitohg .
```
