// SPDX-License-Identifier: Apache-2.0
pragma solidity ^0.8.15;
// import "./IBridge.sol";

import "./ProtobufLib.sol";

library ReportCodec {
    struct Report {
        string qdata;
        uint64 value;
        uint64 timestamp;
    }

    function encode(
        Report memory instance
    ) internal pure returns (bytes memory) {
        bytes memory finalEncoded;

        if (bytes(instance.qdata).length > 0) {
            finalEncoded = abi.encodePacked(
                finalEncoded,
                ProtobufLib.encode_key(
                    1,
                    uint64(ProtobufLib.WireType.LengthDelimited)
                ),
                ProtobufLib.encode_uint64(uint64(bytes(instance.qdata).length)),
                bytes(instance.qdata)
            );
        }

        if (uint64(instance.value) != 0) {
            finalEncoded = abi.encodePacked(
                finalEncoded,
                ProtobufLib.encode_key(2, uint64(ProtobufLib.WireType.Varint)),
                ProtobufLib.encode_uint64(instance.value)
            );
        }

        if (uint64(instance.timestamp) != 0) {
            finalEncoded = abi.encodePacked(
                finalEncoded,
                ProtobufLib.encode_key(3, uint64(ProtobufLib.WireType.Varint)),
                ProtobufLib.encode_uint64(instance.timestamp)
            );
        }

        return finalEncoded;
    }
}
