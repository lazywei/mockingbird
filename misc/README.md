Preparing Code Samples
================


## Collect Rosetta Code

1. Clone the [RosettaCodeData](https://github.com/acmeism/RosettaCodeData)
  ```
  git clone git@github.com:acmeism/RosettaCodeData.git
  ```
2. Build this `misc` executable
  ```
  go build
  ```
3. Run the `collectRosetta` according to the cloned RosettaCodeData, and collect
   files to `../samples`
  ```
  ./misc collectRosetta path/to/clones/RosettaCodeData ../samples
  ```

## Build Bag-of-Words and Convert Samples to Libsvm

```
./misc convertLibsvm ../samples
```
