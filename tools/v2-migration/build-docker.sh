#!/usr/bin/env bash

###############################################################################
# Script options
#
# Reference: http://redsymbol.net/articles/unofficial-bash-strict-mode/
###############################################################################
set -e # exit on error
set -o pipefail # exit on pipe fails
set -u # fail on unset
#set -x # activate debugging from here down
#set +x # disable debugging from here down
# First set bash option to avoid
# unmatched patterns expand as result values
#shopt -s nullglob

###############################################################################
# Configuration variables
###############################################################################

declare -r container_owner="harold"
declare -r container_name="easycla-migration"
declare -r container_version="0.2.0"
declare -r container_fullname="${container_owner}/${container_name}:${container_version}"
declare -r cla_app_folder=${HOME}/projects/go/src/github.com/communitybridge/easycla/cla-backend/cla


if [[ ! -f Dockerfile ]]; then
  echo "Missing Dockerfile in current folder. Exiting..."
  exit 1
fi

if [[ ! -d ${cla_app_folder} ]]; then
  echo "Missing cla source folder: ${cla_app}. Exiting..."
  exit 1
fi

rm -Rf cla
cp -R ${cla_app_folder} cla
docker build -t ${container_fullname} .