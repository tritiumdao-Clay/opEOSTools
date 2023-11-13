import time

import requests

content = [
    {"to":"0xcBB2e4aaa97dEA80dba687e42a455034A39954e6","chain":"jungle4"},
    {"to":"0xe7659B6dc602D920119AB612c6C3Dd9DB7D6A2EC","chain":"jungle4"},
    {"to":"0x1aEB0963d2A169FCD802582dE06682CA43C3a045","chain":"jungle4"}
]
url = 'https://faucet.testnet.evm.eosnetwork.com/api/send'

count = 0
while True:
    result = requests.post(url, json=content[count % len(content)])
    if result.status_code == 200:
        count += 1
    else:
        print(result.status_code, result.text)
    time.sleep(12)
