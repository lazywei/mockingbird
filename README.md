Mockingbird [![Build Status](https://travis-ci.org/lazywei/mockingbird.svg?branch=master)](https://travis-ci.org/lazywei/mockingbird)
===========

# Introduction

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

# Command Line Interface Usage

## Preparing LIBSVM format dataset

### Collect Rosetta Code

1. Clone the [RosettaCodeData](https://github.com/acmeism/RosettaCodeData)

  ```
  git clone git@github.com:acmeism/RosettaCodeData.git
  ```

2. Build this `cli` executable

  ```
  cd cli/
  ./build.sh
  ```

3. Run the `collectRosetta` according to the cloned RosettaCodeData, and collect
   files to `../samples`

  ```
  ./mockingbird collectRosetta path/to/clones/RosettaCodeData ../samples
  ```

### Build Bag-of-Words and Convert Samples to Libsvm

Build from scratch
```
./mockingbird convertLibsvm ../samples ../
```
This will save `libsvm.samples` and `bow.gob` to `../`. The `bow.gob` is the
parameters for constructing bag-of-words. This can be used afterward:

```
./mockingbird convertLibsvm ../samples ../ --bowPath ../bow.gob
```

## Train and Predict

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
