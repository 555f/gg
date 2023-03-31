# Golang code generation tool

[![Go Reference][go-reference-badge]][go-reference]

> _**Note**_: This tool is now available in **BETA**. Please try it out and
> give us feedback!

> _**Note**_: If you believe have found a critical issue in GG, _please responsibly
> disclose_ by contacting us at
> [vitalylobchuk@gmail.com](mailto:vitalylobchuk@gmail.com).

## Installation 

```sh
go mod download
```

```sh
go build -o gg cmd/gg
```

## Use

Create file `gg.yaml` in root project directory:

```yaml
packages:
  - ./internal/...
```

```shell
gg run
```

## Documentation

Coming soon.

## Examples 

Simple examples are located in the [examples](examples) directory.