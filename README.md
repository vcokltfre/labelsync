# labelsync

Sync GitHub repository labels to a YAML file.

## Installation and Usage

Installing labelsync from source:

```sh
git clone https://github.com/vcokltfre/labelsync
cd labelsync
go build && go install
```

or

```sh
go install github.com/vcokltfre/labelsync@latest
```

Installing labelsync from binary:

Visit the [releases page](https://github.com/vcokltfre/labelsync/releases) and download the relevant file for your OS and architecture. Extract the downloaded file and place the binary somewhere in your PATH.

Running labelsync:

To run labelsync you'll need a .env file in the current working directory with a `GITHUB_TOKEN` variable.

```sh
labelsync [schema=.gitlabels.yml]
```

## Examples and Usage

Please see the `examples/` directory for usage examples.
