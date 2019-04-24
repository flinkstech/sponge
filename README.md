# sponge
[![builds.sr.ht status](https://builds.sr.ht/~tslocum/sponge.svg)](https://builds.sr.ht/~tslocum/sponge)
[![Donate](https://img.shields.io/liberapay/receives/rocketnine.space.svg?logo=liberapay)](https://liberapay.com/rocketnine.space)

Soaks up all input from stdin and writes it to a file or stdout. Pipelines
reading from and writing to the same file may be safely constructed.

```
grep [...] log.txt | sponge log.txt
```

## Installation

```
go get git.sr.ht/~tslocum/sponge
```

## Usage

```
sponge [-a] [file]: soak up all input from stdin and write it to [file] or stdout
```


## Credits

sponge was originally written in C by Tollef Fog Heen. 
