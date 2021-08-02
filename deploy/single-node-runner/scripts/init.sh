#!/bin/sh

set -eu

DAEMON=marketd

CHAINID="market-testnet"
CHAIN_VALIDATOR_ACC="validator"

CHAIN_HOME=/root/.market
CHAIN_NODE_CONFIG=$CHAIN_HOME/config/config.toml
CHAIN_APP_CONFIG=$CHAIN_HOME/config/app.toml

CHAIN_HOST=0.0.0.0

KEYRING_FLAG="--keyring-backend test"

# Build genesis file incl account for passed address
coins="10000000000stake,100000000000samoleans"
$DAEMON init --chain-id $CHAINID $CHAINID

# generate validator and 5 more accounts
for account in $CHAIN_VALIDATOR_ACC 'acc1' 'acc2' 'acc3' 'acc4' 'acc5' ;
  do
    echo "creating account $account"
    $DAEMON keys add $account $KEYRING_FLAG --output json >> $CHAIN_HOME/"$account"_key.json
    $DAEMON add-genesis-account $($DAEMON keys show $account -a $KEYRING_FLAG) $coins
  done

$DAEMON gentx $CHAIN_VALIDATOR_ACC 5000000000stake $KEYRING_FLAG --chain-id $CHAINID

$DAEMON collect-gentxs

# Change ports
sed -i "s#\"tcp://127.0.0.1:26656\"#\"tcp://$CHAIN_HOST:26656\"#g" $CHAIN_NODE_CONFIG
sed -i "s#\"tcp://127.0.0.1:26657\"#\"tcp://$CHAIN_HOST:26657\"#g" $CHAIN_NODE_CONFIG
sed -i 's#addr_book_strict = true#addr_book_strict = false#g' $CHAIN_NODE_CONFIG
sed -i 's#external_address = ""#external_address = "tcp://'$CHAIN_HOST:26656'"#g' $CHAIN_NODE_CONFIG
sed -i 's#enable = false#enable = true#g' $CHAIN_APP_CONFIG
sed -i 's#swagger = false#swagger = true#g' $CHAIN_APP_CONFIG
