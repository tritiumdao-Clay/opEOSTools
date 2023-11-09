// SPDX-License-Identifier: MIT
pragma solidity 0.8.20;


contract TestDeposit {

    uint32 public ethFlag;
    uint32 public erc20Flag;

    event EventDepositETH(uint32, bytes);
    event EventDepositERC20(address indexed, address indexed, uint256 indexed, uint32, bytes);

    function depositETH(uint32 gasLimit,bytes memory data) external payable {

        ethFlag += 1;
        emit EventDepositETH(gasLimit, data);
    }

    function depositERC20(address l1Addr,address l2Addr,uint256 value,uint32 gasLimit,
        bytes memory data) external {
        erc20Flag +=1;
        emit EventDepositERC20(l1Addr, l2Addr, value, gasLimit, data);
    }

}
