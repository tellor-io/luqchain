# Steps taken to create chain and its modules

```sh
# scaffold chain
ignite scaffold chain luqchain --address-prefix luq
```

```sh
# scaffold a report type
# qid: key to this type
ignite scaffold type report qdata value:uint timestamp:uint --signer reporter
```

```sh
# scaffold submitVal message
ignite scaffold message submitVal qdata value:uint -d "params are query data in hex string and value as uint"
```

```sh
# scaffold retrieveVal query
ignite scaffold query retrieveVal qid timestamp:uint --response report:Report -d "params are query id (hash of query data) and timestamp of when the report was submitted"
```

```sh
# scaffold retrieveAll query
ignite scaffold query retrieveAll qid --response report:Report --paginated -d "Fetch all the reports for a given query id"
```

Things to consider editing in `config.yml` is the denom of the staking token.
Switch the pubkey type to secp256k1

- Edit the Initcmd function by overriding the cosmos sdk function

- Edit the ```InitializeNodeValidatorFilesFromMnemonic``` that is a dependency of InitCmd function

- Add to ```x/keeper/types/expected_keepers.go``` StakingKeeper interface that adopts methods from staking module

- Add key prefix for the report store in ```x/keeper/types/keys.go```

- Add ```repeated``` and ```[(gogoproto.nullable) = false]``` to the .proto file for getAll query message for a list response then run ```ignite generate proto-go```

<!-- Proof grpc API -->
```sh
# add 2 blockheader type one regular and the other evm
ignite scaffold type blockHeaderMerkle versionChainidHash height:uint timeSecond:uint timeNanosecond:uint lastblockidCommitHash nextvalidatorConsensusHash lastresultsHash evidenceProposerHash

ignite scaffold type blockHeaderMerkleEvm versionChainidHash height:uint64 timeSecond:uint64 timeNanosecond:uint32 lastblockidCommitHash nextvalidatorConsensusHash lastresultsHash evidenceProposerHash
```
