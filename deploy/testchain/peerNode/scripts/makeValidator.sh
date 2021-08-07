#NOW WE MAKE OUR FULL NODE A VALIDATOR NODE
#!/bin/bash
set -eu

echo "building environment"
apt-get install wget nano jq -yq
# Initial dir
CURRENT_WORKING_DIR=~
# Name of the network to bootstrap
CHAINID="testchain"
# Name of the gravity artifact
GRAVITY=gravity
# The name of the gravity node
GRAVITY_NODE_NAME="gravity"
# The address to run gravity node
GRAVITY_HOST="0.0.0.0"
# Home folder for gravity config
GRAVITY_HOME="$CURRENT_WORKING_DIR/$CHAINID/$GRAVITY_NODE_NAME"
# Home flag for home folder
$GRAVITY_HOME_FLAG="--home $GRAVITY_HOME"
# Config directories for gravity node
GRAVITY_HOME_CONFIG="$GRAVITY_HOME/config"
# Config file for gravity node
GRAVITY_NODE_CONFIG="$GRAVITY_HOME_CONFIG/config.toml"
# App config file for gravity node
GRAVITY_APP_CONFIG="$GRAVITY_HOME_CONFIG/app.toml"
# Keyring flag
GRAVITY_KEYRING_FLAG="--keyring-backend test"
# Chain ID flag
GRAVITY_CHAINID_FLAG="--chain-id $CHAINID"
# The name of the gravity validator
GRAVITY_VALIDATOR_NAME=val2
# The name of the gravity orchestrator/validator
GRAVITY_ORCHESTRATOR_NAME=orch2
# Gravity chain demons
STAKE_DENOM="stake"
#NORMAL_DENOM="samoleans"
NORMAL_DENOM="footoken"
# Moniker of orchestrator
MONIKER_ORCH="popular grant rural draft unhappy equal service expire evoke topple ozone lens chapter female soda fun hair clock century rail student robot prize mosquito"

# Recover the orchestrator to take some token from it
$GRAVITY_HOME_FLAG keys add orch1 --recover $$GRAVITY_KEYRING_FLAG <<< \"$MONIKER_ORCH\"

# Transfer some stake token to new validator
$GRAVITY_HOME_FLAG tx bank send $($GRAVITY_HOME_FLAG keys show -a orch1 $GRAVITY_KEYRING_FLAG) $($GRAVITY_HOME_FLAG keys show -a val2 $GRAVITY_KEYRING_FLAG) 10000000stake $GRAVITY_CHAINID_FLAG $GRAVITY_KEYRING_FLAG -y

# Transfer some footoken to new validator
$GRAVITY_HOME_FLAG tx bank send $($GRAVITY_HOME_FLAG keys show -a orch1 $GRAVITY_KEYRING_FLAG) $($GRAVITY_HOME_FLAG keys show -a val2 $GRAVITY_KEYRING_FLAG) 10000000footoken $GRAVITY_CHAINID_FLAG $GRAVITY_KEYRING_FLAG -y

# Stor the public key of validator
PUB_KEY=$($GRAVITY_HOME_FLAG tendermint show-validator)

# Do the create validator transaction
$GRAVITY_HOME_FLAG tx staking create-validator \
--amount=100000000$STAKE_DENOM \
--pubkey=\"$PUB_KEY\" \
--moniker=\"$GRAVITY_VALIDATOR_NAME\" \
--chain-id=$CHAINID \
--commission-rate="0.10" \
--commission-max-rate="0.20" \
--commission-max-change-rate="0.01" \
--min-self-delegation="10" \
--gas="auto" \
--gas-adjustment=1.5 \
--gas-prices=\"1$NORMALDENOM\" \
--from=$GRAVITY_VALIDATOR_NAME \
$GRAVITY_KEYRING_FLAG