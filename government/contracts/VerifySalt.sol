pragma solidity ^0.8.9;

contract A {
    uint public a = 1;
}
contract VerifySalt {

    address public addr;


    function createA() external returns(address) {
        //bytes32 salt = keccak256(abi.encode("hello"));
        //addr = address(new A{salt: salt}());
        addr = address(new A());
        return addr;
    }

    function createB() external returns(bytes32 salt) {
        salt = keccak256(abi.encode("hello"));
    }

    function createC() external returns(address) {
        bytes32 salt = keccak256(abi.encode("hello"));
        addr = address(new A());
        return addr;
    }

    function createD() external returns(address) {
        bytes32 salt = keccak256(abi.encode("hello"));
        addr = address(new A{salt: salt}());
        return addr;
    }

    function createE(uint256 a) external {
        require(a==1, "not equal 1");
    }

}
