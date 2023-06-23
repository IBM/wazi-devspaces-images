#!/bin/bash

###############################################################################
# Licensed Materials - Property of IBM.
# Copyright IBM Corporation 2023. All Rights Reserved.
# U.S. Government Users Restricted Rights - Use, duplication or disclosure
# restricted by GSA ADP Schedule Contract with IBM Corp.
#
# Contributors:
#  IBM Corporation - initial API and implementation
###############################################################################

set -e

# ---------- Command arguments ----------
scriptName=$(basename "$0")
scriptDir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
arch="$(uname -m)"
system="$(uname -s)"

# ---------- Global variables ----------
DS_CONFIG="ds_config.json"
CHE_DEVWORKSPACE_GENERATOR="che-devworkspace-generator"
VERSION="$(cat VERSION)"

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

function generateLocalDevworkspaces() {
    
    for dir in ./devfiles/*/
    do

        devworkspaces=("${dir}"devworkspace*)
        [[ -e "${devworkspaces[0]}" ]] && continue
        ! [[ -f "${dir}"devfile.yaml ]] && continue

        name=$(yq -r '.projects[0].name' "${dir}"devfile.yaml)
        project="${name}={{_INTERNAL_URL_}}/resources/v2/${name}.zip"

        # Generate devworkspace-che-code-insiders.yaml
        npm_config_yes=true npx @eclipse-che/${CHE_DEVWORKSPACE_GENERATOR} \
        --devfile-path:"${dir}"devfile.yaml \
        --editor-entry:che-incubator/che-code/insiders \
        --plugin-registry-url:https://redhat-developer.github.io/devspaces/che-plugin-registry/"${VERSION}"/"${arch}"/v3 \
        --output-file:"${dir}"devworkspace-che-code-insiders.yaml \
        --project."${project}"

        # Generate devworkspace-che-code-latest.yaml
        npm_config_yes=true npx @eclipse-che/${CHE_DEVWORKSPACE_GENERATOR} \
        --devfile-path:"${dir}"devfile.yaml \
        --editor-entry:che-incubator/che-code/latest \
        --plugin-registry-url:https://redhat-developer.github.io/devspaces/che-plugin-registry/"${VERSION}"/"${arch}"/v3 \
        --output-file:"${dir}"devworkspace-che-code-latest.yaml \
        --project."${project}"

        # Generate devworkspace-che-idea-latest.yaml
        npm_config_yes=true npx @eclipse-che/${CHE_DEVWORKSPACE_GENERATOR} \
        --devfile-path:"${dir}"devfile.yaml \
        --editor-entry:che-incubator/che-idea/latest \
        --plugin-registry-url:https://redhat-developer.github.io/devspaces/che-plugin-registry/"${VERSION}"/"${arch}"/v3 \
        --output-file:"${dir}"/devworkspace-che-idea-latest.yaml \
        --project."${project}"
    done
}

function generateLocalProjects() {

    [[ -e $DS_CONFIG ]] && {

        mkdir -p ./resources/v2/
        source "$scriptDir/clone_and_zip.sh"
        local projects=$(cat $DS_CONFIG | jq -c '.projects[]' 2> /dev/null)
        
        for project in ${projects[@]}
        do
            local projectName=$(jq -r '.projectName' <<< $project)
            local projectURL=$(jq -r '.projectUrl' <<< $project)
            if [[ $projectName != null ]] && [[ $projectURL != null ]]; then
                
                projectRepo="${projectURL%/tree*}.git"
                projectBranch=${projectURL##*/}
                projectArchive="$(pwd)/resources/v2/${projectName}.zip"
                
                clone_and_zip "${projectRepo}" "${projectBranch}" "${projectArchive}"
            fi

        done
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
  --devworkspaces            : Generate local Devworkspaces
  --projects                 : Generate local Projects
  --help | -h                : This usage information

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
            '--devworkspaces') generateLocalDevworkspaces; exit $?; shift 0;;
            '--projects') generateLocalProjects; exit $?; shift 0;;
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
