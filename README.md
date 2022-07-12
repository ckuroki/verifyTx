# Merkle proof — Verify Transaction Integrity

### Requirements

- make

- go compiler

### Build

```
make
```

##### Usage

```sh
verify_tx <txHash>
```

Example

```
./cmd/verify_tx/verify_tx 0x0b41fc4c1d8518cdeda9812269477256bdc415eb39c4531885ff9728d6ad096b
Found : 4 transactions on block : 10593417
0xab41f886be23cd786d8a69a72b0f988ea72e0b2e03970d0798f5e03763a442cc : 0xab41f886be23cd786d8a69a72b0f988ea72e0b2e03970d0798f5e03763a442cc
Verification successful
```

#### Docs

![Docs](docs/trie.md)
