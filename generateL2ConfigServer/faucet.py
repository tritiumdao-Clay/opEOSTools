import time
import requests

content = {"to":"0x17F42C02Ca0094f53e813A0F3e0ACfBf22b732Ef","chain":"jungle4"}
url = 'https://faucet.testnet.evm.eosnetwork.com/api/send'

count = 0
while True:
    result = requests.post(url, json=content)
    if result.status_code == 200:
        count += 1
        print('领取成功, 次数:', count)
    else:
        print(result.status_code, result.text)
    time.sleep(13)
