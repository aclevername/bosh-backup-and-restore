#!/bin/bash

set -eu

cd "bbr-systest-releases/${RELEASE_NAME}"
bosh -n create release --force

bosh -n target "$BOSH_HOST"
bosh login "$BOSH_CLIENT" "$BOSH_CLIENT_SECRET"
bosh upload release --rebase

bosh \
  -n \
  --ca-cert=../../bosh-backup-and-restore-meta/certs/lite-bosh-uaa.backup-and-restore.cf-app.com.crt \
  target "$BOSH_UAA_HOST"

export BOSH_CLIENT_SECRET=$BOSH_UAA_CLIENT_SECRET
bosh login
bosh upload release --rebase
