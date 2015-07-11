Mockingbird [![Build Status](https://travis-ci.org/lazywei/mockingbird.svg?branch=master)](https://travis-ci.org/lazywei/mockingbird)
===========

## Introduction

[Linguist](https://github.com/github/linguist)'s Classifier in Go.

Linguist can be used as a Go package by

```golang
import "github.com/lazywei/linguist"
```

and it also has a CLI (command line interface) in `cli/`

```
$ cd cli/
$ ./build.sh
$ ./mockingbird --help
```

## Command Line Interface Usage

### Train

```
usage: mockingbird train [<flags>]

Train Classifier

Flags:
  --help            Show help (also see --help-long and --help-man).
  --sample="samples.libsvm"
                    Path for samples (in libsvm format)
  --output="model"  Path for saving trained model
```

### Predict

```
usage: mockingbird predict

Predict via trained Classifier

Flags:
  --help  Show help (also see --help-long and --help-man).
```
