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

PLUGINS=()

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

function getPluginData() {

    local sourceDir="$1"
    local pluginsDir="${sourceDir%/}/plugins"
    local packageJson="extension/package.json"

    for plugin in $pluginsDir/*.vsix
    do

        [[ -f $plugin ]] || continue

        info "Plugin: ${plugin##*/} found."

        local publisher=$(unzip -p $plugin $packageJson | jq -r ".publisher")
        local name=$(unzip -p $plugin $packageJson | jq -r ".name")
        local pluginVersion=$(unzip -p $plugin $packageJson | jq -r ".version")
        local major="${pluginVersion%%.*}"
        local minor="${pluginVersion%.*}"; minor="${minor#*.}"
        local patch="${pluginVersion##*.}"; patch="${patch%%-*}"
        PLUGINS+=( "$plugin:$major.$minor.$patch:$publisher:$name" )

    done
}

function processPlugins() {

    local openVSXJson="/openvsx-server/openvsx-sync.json"
    local openVSXVsixDir="/openvsx-server/vsix"

    for plugin in ${PLUGINS[@]}
    do

        IFS=":" read -r -a chunk <<< "${plugin}"

        local source="${chunk[0]}"
        local version="${chunk[1]}"
        local publisher="${chunk[2]}"
        local name="${chunk[3]}"
        local pluginId="$publisher.$name"
        local name_version="$pluginId-$version.vsix"

        [[ -f $openVSXJson ]] && {

            # Handle OpenVSX
            local openVSXVersion=$(jq -r --arg id "$pluginId" '(.[] | select(.id == $id) | .version)' "$openVSXJson")

            [[ -z $openVSXVersion ]] && {
                tmpfile=$(mktemp)
                jq --arg id "$pluginId" --arg ver "$version" '. += [{ "id": $id, "version": $ver }]' "$openVSXJson" > "$tmpfile"
                mv -fv "$tmpfile" "$openVSXJson"
            }

            [[ "$openVSXVersion" != "$version" ]] && {
                tmpfile=$(mktemp)
                jq --arg id "$pluginId" --arg ver "$version" '(.[] | select(.id == $id) | .version) |= $ver' "$openVSXJson" > "$tmpfile"
                mv -fv "$tmpfile" "$openVSXJson"
                rm -rfv "$openVSXVsixDir/$pluginId-$openVSXVersion.vsix"
            }
            
            mv -vf "$source" "$openVSXVsixDir/$name_version"

        } # [[ -f $openVSXJson ]]

    done
}

# Main Function
function main() {

    getPluginData "${@}"
    processPlugins
}

# --- Run ---
main "${@}"

# First parameter - relative or absolute folder location for the v3 directory

# Notes:
# The source comes from OpenVSX originally, therefore if there is major deviation then this script needs to be reworked
# yq tool is installed via python and wrapped with jq, which converts yaml to json for processing
