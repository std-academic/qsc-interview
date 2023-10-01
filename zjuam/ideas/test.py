import requests
import math
from Crypto.PublicKey import RSA

req = requests.get("https://zjuam.zju.edu.cn/cas/v2/getPubKey")
pubkey_info = req.json()

pubkey_info['modulus'] = 'e0d3c43d75f316cdec3be4201981662f2be1cc609da4126c6fc1caaab4ebdeeee019b57916113151c3144afad168fd0d2168f7d737f9a8f9af80d7db55899c23'
key = RSA.construct((
    int(pubkey_info['modulus'], 16),
    int(pubkey_info['exponent'], 16)
))

def enc_pass(password: str):
    # code from https://github.com/pyca/cryptography/issues/2735
    # zero padding RSA encrypt
    e = key.e
    n = key.n
    keylen = math.ceil(n.bit_length() / 8)
    input_nr = int.from_bytes(password[::-1].encode('ascii'), byteorder='big')
    crypted_nr = pow(input_nr, e, n)
    crypted_data = crypted_nr.to_bytes(keylen, byteorder='big')
    return crypted_data.hex()

print(enc_pass('testtest'))

# a89804bb826caa5006ae17d59c2da81e05de99bdbe925c2835d72fc1bd145405b88f60b40f029d44e860f7827663e8904230f16f515e0749619cb2e618932fa7