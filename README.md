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
geth_repo: "https://github.com/ethereum/go-ethereum"
geth_branch: "v1.14.4"
path: "./"
build_flags: ""
output_dir: "./build"
```

Note: in this example paths are relative to the location of the config file.

### Run

Once you have the configuration set-up, there are commands to build the geth binary, or create an archive file for distribution, or run the tracer tests.

#### Build

To bundle the plugin together with geth run:

```terminal
./bin/geth-builder build --config simple/geth-builder.yaml
```

It will output the geth binary at `./simple/build/geth`. The configuration fields can also be passed as flags, in which case they will override the values in the config file.

```terminal
./bin/geth-builder build --path ./simple --output ./build
```

This will build the default geth repository and output the binary to `./build/geth`.

#### Test

If you have Go tests in the package geth-builder will allow you to run those against the same go-ethereum version as configured. To run the tests:

```terminal
./bin/geth-builder -v test --config simple/geth-builder.yaml
```

#### Archive

Once you have tested the plugin and are ready to distribute it, you can create an archive file:

```terminal
./bin/geth-builder archive --config simple/geth-builder.yaml --type [zip|tar]
```

This will create a zip or tar archive in the output directory. It is possible to cross-build for different platforms:

```terminal
./bin/geth-builder archive --config simple/geth-builder.yaml --arch amd64
```

The `--arch` flag follows the same naming convention as `GOARCH`.