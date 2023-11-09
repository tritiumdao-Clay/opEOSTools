
yiran地址:
0xC9DB2f861dec8d63bC2203E6CC727b6cff233d74
0x5AEcBcEc7Dc59E3D89fe1d53e3d00F486530F3eF

- golang服务:
  - 实现目的:
    - 提供prove参数
    - 提供finalize参数
    - 获取L2跨链状态(提现完成, 跨链完成)
    - 提供address ,txHash
    - 获取address, txHash 
      - txStatus = 0 --> 提现成功;
        txStatus = 1 ---> prove成功;
        txStatus = 2 ----> finalize成功(最终状态)
        txStatus = -1 ---> 提现失败;