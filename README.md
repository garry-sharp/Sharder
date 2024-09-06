<div align="center">

![Bitcoin](https://img.shields.io/badge/Bitcoin-000?style=for-the-badge&logo=bitcoin&logoColor=white)
![Ethereum](https://img.shields.io/badge/Ethereum-3C3C3D?style=for-the-badge&logo=Ethereum&logoColor=white)
![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)

[![codecov](https://codecov.io/github/garry-sharp/Sharder/graph/badge.svg?token=IWGRIG1DDF)](https://codecov.io/github/garry-sharp/Sharder)
[![build_and_test](https://github.com/garry-sharp/Sharder/actions/workflows/build_and_test.yml/badge.svg)](https://github.com/garry-sharp/Sharder/actions/workflows/build_and_test.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

</div>

# Sharder

Sharder is a lightweight command-line tool designed to simplify the process of sharding mnemonics. With Sharder, you can effortlessly split your mnemonics into smaller, more manageable chunks, making it easier to organize and store sensitive information securely. Whether you're working with cryptographic keys, passwords, or any other mnemonic-based data, Sharder provides a seamless solution for dividing and distributing your mnemonics across multiple locations. Take control of your mnemonic sharding process with Sharder and ensure the utmost security for your valuable data.

This tool is free for anyone to use. If you like it please consider donating [bitcoin](bitcoin:bc1qvt37xsc3980zk3nvg44dn92vg2whq73xzsxlna) or [ethereum](ethereum:0x61ae64504549432a94D09E0C258c981698253F7A)

## Main Use Case & Background

The inspiration for this project was simply the question "what would happen to your crypto if you suddenly died?", well the answer is quite clear, you can either trust other people with your keys or your loved ones will lose their crypto.

![image](docs/Shamir.png)

The shamir secret sharing algorithm is an elegant technique that uses coordinates along a polynomial to encode the secret. 2 properties are defined, a threshold and a total which generate the `shares`. So long as the threshold number of shares are presented back to the application, a secret can be recovered, without disclosing any information of the secret in the shares themselves.

## CLI

[See Here](docs/cryptosharder.md) for CLI documentation

## Not supported

The following are a list of features that aren't supported

-   Mnemonics with passphrases
