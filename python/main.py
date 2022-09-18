"""
Decode transactions for WETH-USDC uniswap pool
https://etherscan.io/address/0x88e6a0c2ddd26feeb64f039a2c41296fcb3f5640#tokentxns
"""

import os
from web3 import Web3
from web3.logs import DISCARD
from abi import UNISWAP_ABI


def fetch_swap_transactions(tx_id: str):
    """Get transactions"""

    infura_node = os.getenv("INFURA_NODE")
    if infura_node  is None:
        print("Error: Need to set INFURA_NODE in env")

    web3_client = Web3(Web3.HTTPProvider(infura_node))

    receipts = web3_client.eth.get_transaction_receipt(tx_id)
    uniswap_v3_usdc = "0x88e6A0c2dDD26FEEb64F039a2c41296FcB3f5640"
    contract = web3_client.eth.contract(address=uniswap_v3_usdc, abi=UNISWAP_ABI)
    events = contract.events.Swap().processReceipt(receipts, DISCARD)


    for event in events:
        # print(event)
        args = event["args"]
        amount_0 = args["amount0"]
        amount_1 = -args["amount1"]
        print(f"Swapped {amount_0 * 0.000001} USDC for {amount_1 * 0.000000000000000001} ETH")


fetch_swap_transactions("0x5e3eb8b97677395319973153013db75b76fc453d3f5bfd459cc3bc5d4f44a56e")
