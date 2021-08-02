#!/bin/sh

set -eu

DAEMON=marketd

# Start the node
$DAEMON start --pruning=nothing
