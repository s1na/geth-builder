# geth-builder

Geth-builder builds geth (surprise!) together with [native](https://geth.ethereum.org/docs/developers/evm-tracing/custom-tracer#custom-go-tracing) or [live](https://geth.ethereum.org/docs/developers/evm-tracing/live-tracing) tracers. It allows developers to have a clean repository for their tracing package without forking go-ethereum or having it as a submodule.

## Install

Fetch the source code from github. Then run:

```terminal
make build
# Binary will be in ./bin/geth-builder
```

## Usage

Create a package which will contain the tracer file and configuration.

```
mkdir simple
cd simple
```

Initialize a config file:

```terminal
./bin/geth-builder init
```

This will dump default configuration to `geth-builder.yaml`. In the simplest case the complete package structure will be as follows:

```
simple/
  - tracer.go
  - geth-builder.yaml
```

Note: The go files should have the same package name as the directory.

### Configuration

The config files will determine where to fetch the go-ethereum source from, where to locate the tracing package and the output directory:

```yaml
eth_repo: "https://github.com/ethereum/go-ethereum"
geth_branch: "v1.14.4"
path: "./"
build_flags: ""
output_dir: "./build"
```

Note: in this examples paths are relative to the location of the config file.

### Run

Now to bundle the plugin together with geth run:

```terminal
./bin/geth-builder --config simple/geth-builder.yaml
```

It will output the geth binary at `./simple/build/geth`.
