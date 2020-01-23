# sponge
[![GoDoc](https://godoc.org/gitlab.com/tslocum/sponge?status.svg)](https://godoc.org/gitlab.com/tslocum/sponge)
[![CI status](https://gitlab.com/tslocum/sponge/badges/master/pipeline.svg)](https://gitlab.com/tslocum/sponge/commits/master)
[![Donate](https://img.shields.io/liberapay/receives/rocketnine.space.svg?logo=liberapay)](https://liberapay.com/rocketnine.space)

Soaks up all input from stdin and writes it to a file or stdout. Pipelines
reading from and writing to the same file may be safely constructed.

```
grep [...] log.txt | sponge log.txt
```

## Installation

```
go get gitlab.com/tslocum/sponge
```

## Usage

```
sponge [-a] [file]: soak up all input from stdin and write it to [file] or stdout
```

## Support

Please share issues and suggestions [here](https://gitlab.com/tslocum/sponge).

## Credits

sponge was originally written in C by Tollef Fog Heen. 
