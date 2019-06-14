# mitohg

## Description

Определение гаплогрупп происходит на основе номенклатуры [PhyloTree build 17](http://www.phylotree.org) относительно референсной последовательности [RSRS](http://www.phylotree.org/resources/RSRS.fasta).

Для корректной работы программы переменная окружения `MUSCLE_BIN` должна указывать на бинарный файл к [muscle](http://drive5.com/muscle/downloads.htm). mitohg понимает `.env` файлы.

## Build

```
go build -o mitohg
```
