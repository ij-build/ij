#!/bin/sh -e

SCRATCH_DIR=/scratch-keys

# Move volume-mounted keys to another directory.
# This is necessary as we can't directly change 
# the permissions of files in a Windows volume.
mkdir $SCRATCH_DIR
cp /root/.ssh/* $SCRATCH_DIR

# Set permissions of private keys so that they
# can bge accepted by ssh-add.
chmod 700 $SCRATCH_DIR
chmod 400 $SCRATCH_DIR/*

# Add each private key
find $SCRATCH_DIR -type f -exec grep -l "PRIVATE" {} \; | xargs ssh-add &> /dev/null
