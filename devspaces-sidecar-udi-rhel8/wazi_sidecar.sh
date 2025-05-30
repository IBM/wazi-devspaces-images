#!/bin/bash

###############################################################################
# Licensed Materials - Property of IBM.
# Copyright IBM Corporation 2023, 2024. All Rights Reserved.
# U.S. Government Users Restricted Rights - Use, duplication or disclosure
# restricted by GSA ADP Schedule Contract with IBM Corp.
#
# Contributors:
#  IBM Corporation - initial API and implementation
###############################################################################

set -e

# ***** GLOBALS *****

# ---------- Command arguments ----------
scriptName=$(basename "$0")
scriptDir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
arch="$(uname -m)"
system="$(uname -s)"

# ---------- Global variables ----------


# ---------- Script functions -----------

# Info reporting functions
info() {
    echo >&1 "[INFO] $1"
}

# Warning reporting functions
warn() {
    echo >&1 "[WARNING] $1"
}

# Error reporting functions
err() {
    echo >&2 "[ERROR] $1"
}

err_exit() {
    echo >&2 "[ERROR] $1"
}

# ---------------------------------------

function cleanup() {
    npm cache clean --force
}

function setup_permissions() {

    for resource in "${HOME}"; do
        info "Setting permissions for ${resource}...(please wait)"
        chmod -R go+rwX ${resource}
    done
}

function setup_npmrc() {

    local npmrc_file="$1"   # $1 == NPM RC File
    local npm_uri="$2"      # $2 == NPM URI
    local npm_registry="$3" # $3 == NPM registry
    local npm_username="$4" # $4 == NPM username
    local npm_password="$5" # $5 == NPM password

    [[ $npmrc_file ]]   || err_exit "Value for NPM RC File is not set"
    [[ $npm_uri ]]      || err_exit "Value for NPM URI is not set"
    [[ $npm_registry ]] || err_exit "Value for NPM registry is not set"
    [[ $npm_username ]] || err_exit "Value for NPM user is not set"
    [[ $npm_password ]] || err_exit "Value for NPM password is not set"

    npm_registry="$npm_uri/$npm_registry"

    IFS= read -r -d '' npmrc_data <<-EOF || :
			@ibm:registry=https://${npm_registry}/
			//${npm_registry}/:_password=${npm_password}
			//${npm_registry}/:username=${npm_username}
			//${npm_registry}/:email=${npm_username}
			//${npm_registry}/:always-auth=true
			EOF

    [[ $npmrc_data ]] && {
        info "Setting up ${npmrc_file} for ${npm_registry}"
        echo "${npmrc_data}" > "${npmrc_file}"
    }
}

# ---------- Base functions ----------

# Print usage information
function print_usage() {

  echo "

USAGE: ./${scriptName} [OPTIONS]

OPTIONS:
"

  echo "
  --cleanup                                                          : Runs a cleanup to remove files/folders
  --permissions                                                      : Setup permissions for various well-known locations
  --npmrc 'npmrcfile' 'npmuri' 'npmregistry' 'username' 'password'   : Setup npm rc file for @ibm eu.artifactory
  --help | -h                                                        : This usage information

EXAMPLES:
./${scriptName} -h
"
}

# ---------- Main functions ----------

# Parses the CLI arguments
function parse_arguments() {

    while [[ "$#" -gt 0 ]]; do
        case $1 in
            '-h'|'--help') print_usage; exit 0; shift 0;;
            '-d'|'--debug') set -x; shift 0;;
            '--cleanup') cleanup; exit $?; shift 0;;
            '--permissions') setup_permissions; exit $?; shift 0;;
            '--npmrc') setup_npmrc "$2" "$3" "$4" "$5" "$6"; exit $?; shift 0;;
            *) err_exit "Invalid Option $2"; shift 0;;
        esac
        shift 1
    done
}

# Main Function
function main() {

    # parse command arguments
    parse_arguments "${@}"
}

# --- Run ---
main "${@}"
