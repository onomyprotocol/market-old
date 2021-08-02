# Single node runner

Table of Contents
=================
*  [Description](#Description)
*  [Entrypoints](#Entrypoints)
*  [Build](#Build&Run)

## Description

The "single-node-runner" is a docker image that contains already prebuilt market home files and bootstraps 
the testnet of the market chain with one validator. 

## Entrypoints

0.0.0.0 should be change to a host of a deployed container

- gravity swagger: [http://0.0.0.0:1317/swagger/](http://0.0.0.0:1317/swagger/)
- gravity rpc: [http://0.0.0.0:1317/](http://0.0.0.0:1317/)
- gravity grpc: [http://0.0.0.0:9090/](http://0.0.0.0:9090/)

## Build&Run

### Build locally

  ```
  docker build -t onomy/market-single-node-runner:local  -f Dockerfile ../../
  ```
### Run in docker-compose

  ```
  docker-compuse up  
  ```

## Accounts

The generated accounts with their keys are localed in the folder [market_home](./market_home)




