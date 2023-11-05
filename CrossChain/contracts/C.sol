pragma solidity ^0.8.9;

contract C {

    bytes32 public a;
    bytes32 public b;

    function proposeL2Output(
        bytes32 _outputRoot,
        uint256 _l2BlockNumber,
        bytes32 _l1BlockHash,
        uint256 _l1BlockNumber
    ) external payable {
        a = blockhash(_l1BlockNumber);
        b = _l1BlockHash;
        //if (_l1BlockHash != bytes32(0)) {
        //    require(
        //        blockhash(_l1BlockNumber) == _l1BlockHash,
        //        "L2OutputOracle: block hash does not match the hash at the expected height"
        //    );
        //}
    }

}