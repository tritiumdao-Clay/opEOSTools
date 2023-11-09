// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "./Hashing.sol";

//0x68F44f4e37c2217cC5b5b337dB43DEE4B225d825
//cast call --rpc-url $L1_RPC 0x68F44f4e37c2217cC5b5b337dB43DEE4B225d825 "hashOutputRootProof((bytes32,bytes32,bytes32,bytes32))" (
//cast call --rpc-url $L1_RPC 0x68F44f4e37c2217cC5b5b337dB43DEE4B225d825 "hashWithdrawal((bytes32,bytes32,bytes32,bytes32))" (
contract Hashing2 {

    bytes32 public outputRoot;
    bytes32 public version;
    bytes32 public stateRoot;
    bytes32 public messagePasserStorageRoot;
    bytes32 public latestBlockhash;

    function hashOutputRootProof(Types.OutputRootProof memory _outputRootProof) external pure returns (bytes32) {
        return Hashing.hashOutputRootProof(_outputRootProof);
    }

    function hashOutputRootProof2(Types.OutputRootProof memory _outputRootProof) external {
        version = _outputRootProof.version;
        stateRoot = _outputRootProof.stateRoot;
        messagePasserStorageRoot = _outputRootProof.messagePasserStorageRoot;
        latestBlockhash = _outputRootProof.latestBlockhash;
        outputRoot = Hashing.hashOutputRootProof(_outputRootProof);
    }

    function hashWithdrawal(Types.WithdrawalTransaction memory _tx) external pure returns (bytes32) {
        return Hashing.hashWithdrawal(_tx);
    }

}
