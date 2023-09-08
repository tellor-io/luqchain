package bridge

var relayFormat = []byte(`
[
    {
        "components": [
            {
                "internalType": "bytes32",
                "name": "luqchainIAVLStateHash",
                "type": "bytes32"
            },
            {
                "internalType": "bytes32",
                "name": "mintStoreMerkleHash",
                "type": "bytes32"
            },
            {
                "internalType": "bytes32",
                "name": "icacontrollerToIcahostMerkleHash",
                "type": "bytes32"
            },
            {
                "internalType": "bytes32",
                "name": "feegrantToIbcMerkleHash",
                "type": "bytes32"
            },
            {
                "internalType": "bytes32",
                "name": "accToEvidenceMerkleHash",
                "type": "bytes32"
            },
            {
                "internalType": "bytes32",
                "name": "paramsToVestingMerkleHash",
                "type": "bytes32"
            }
        ],
        "internalType": "struct Bridge.MultistoreData",
        "name": "self",
        "type": "tuple"
    },
    {
        "components": [
            {
                "internalType": "bytes32",
                "name": "versionAndChainIdHash",
                "type": "bytes32"
            },
            {
                "internalType": "uint64",
                "name": "height",
                "type": "uint64"
            },
            {
                "internalType": "uint64",
                "name": "timeSecond",
                "type": "uint64"
            },
            {
                "internalType": "uint32",
                "name": "timeNanoSecondFraction",
                "type": "uint32"
            },
            {
                "internalType": "bytes32",
                "name": "lastBlockIdCommitMerkleHash",
                "type": "bytes32"
            },
            {
                "internalType": "bytes32",
                "name": "nextValidatorConsensusMerkleHash",
                "type": "bytes32"
            },
            {
                "internalType": "bytes32",
                "name": "lastResultsHash",
                "type": "bytes32"
            },
            {
                "internalType": "bytes32",
                "name": "evidenceProposerMerkleHash",
                "type": "bytes32"
            }
        ],
        "internalType": "struct Bridge.BlockHeaderMerklePartsData",
        "name": "self",
        "type": "tuple"
    },
    {
        "components": [
            {
                "internalType": "bytes",
                "name": "signedDataPrefix",
                "type": "bytes"
            },
            {
                "internalType": "bytes",
                "name": "signedDataSuffix",
                "type": "bytes"
            }
        ],
        "internalType": "struct Bridge.CommonPartsData",
        "name": "commonEncodedVotePart",
        "type": "tuple"
    },
    {
        "components": [
            {
                "internalType": "bytes32",
                "name": "r",
                "type": "bytes32"
            },
            {
                "internalType": "bytes32",
                "name": "s",
                "type": "bytes32"
            },
            {
                "internalType": "uint8",
                "name": "v",
                "type": "uint8"
            },
            {
                "internalType": "bytes",
                "name": "encodedTimestamp",
                "type": "bytes"
            }
        ],
        "internalType": "struct Bridge.TMSignatureData[]",
        "name": "signatures",
        "type": "tuple[]"
    }
]
`)

var verifyFormat = []byte(`
[
    {
		"internalType": "uint256",
		"name": "blockHeight",
		"type": "uint256"
	},
    {
        "components": [
            {
                "internalType": "string",
                "name": "qdata",
                "type": "string"
            },
            {
                "internalType": "uint64",
                "name": "value",
                "type": "uint64"
            },
            {
                "internalType": "uint64",
                "name": "timestamp",
                "type": "uint64"
            }
        ],
        "internalType": "struct ReportCodec.Report",
        "name": "",
        "type": "tuple"
    },
	{
		"internalType": "uint256",
		"name": "version",
		"type": "uint256"
	  },
    {
        "components": [
            {
                "internalType": "bool",
                "name": "isDataOnRight",
                "type": "bool"
            },
            {
                "internalType": "uint8",
                "name": "subtreeHeight",
                "type": "uint8"
            },
            {
                "internalType": "uint256",
                "name": "subtreeSize",
                "type": "uint256"
            },
            {
                "internalType": "uint256",
                "name": "subtreeVersion",
                "type": "uint256"
            },
            {
                "internalType": "bytes32",
                "name": "siblingHash",
                "type": "bytes32"
            }
        ],
        "internalType": "struct Bridge.IAVLData[]",
        "name": "self",
        "type": "tuple[]"
    }
]
`)
