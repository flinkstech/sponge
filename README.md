# sponge
This repository was forked from [this repository](https://code.rocketnine.space/tslocum/sponge)

Soaks up all input from stdin and writes it to a file or stdout. Pipelines
reading from and writing to the same file may be safely constructed.

```grep [...] log.txt | sponge log.txt```

## Installation

```go get gitlab.com/tslocum/sponge```

Or download the binary in the [release page](https://github.com/flinkstech/sponge/releases)

## Usage

```sponge [-a] [file]: soak up all input from stdin and write it to [file] or stdout```

## Support

Please share issues and suggestions [here](https://github.com/flinkstech/sponge/issues).

## Credits

sponge was originally written in C by Tollef Fog Heen.

This Go port was initially written by [Trevor Slocum](https://code.rocketnine.space/tslocum/)