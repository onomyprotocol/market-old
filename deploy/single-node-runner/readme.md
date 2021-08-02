# Single node runner

Table of Contents
=================
*  [Description](#Description)
*  [Entrypoints](#Entrypoints)
*  [Build](#Build&Run)
*  [Accounts](#Accounts)

## Description

The "single-node-runner" is a docker image that contains already prebuilt market home files and bootstraps 
the testnet of the market chain with one validator. 

## Entrypoints

0.0.0.0 should be changed to a host of a deployed container

- gravity swagger: [http://0.0.0.0:1317/swagger/](http://0.0.0.0:1317/swagger/)
- gravity rpc: [http://0.0.0.0:1317/](http://0.0.0.0:1317/)
- gravity grpc: [http://0.0.0.0:9090/](http://0.0.0.0:9090/)

## Build&Run

### Build locally

  ```
  docker build -t onomy/market-single-node-runner:local  -f Dockerfile ../../
  ```
### Run locally in docker-compose

  ```
  docker-compuse up  
  ```


### Run with docker (image from the dockerhub)

  ```
  docker run --name market-single-node-runner \
              -p 26656:26656 -p 26657:26657 -p 1317:1317 -p 61278:61278 -p 9090:9090 \
              -v /mnt/volume_sfo3_01:/root/.market/data/. \
              -it --restart on-failure onomy/market-single-node-runner:latest
  ```

  **latest** here is a tag of the runner. You can get the full list on the page [tags](https://hub.docker.com/repository/docker/onomy/market-single-node-runner/tags?page=1&ordering=last_updated)

  The docker command uses local "/mnt/volume_sfo3_01" directory to save market db files, market_home/data/priv_validator_state.json
  file should be there before the first run of container.

## Accounts

The generated accounts with their keys are localed in the folder [market_home](./market_home)




