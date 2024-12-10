#!/bin/bash

# Script variables
scriptName=$(basename "$0")

# Common functions
err_exit() {
    echo >&2 "[ERROR] $1"
}

info() {
    echo "Info: $1"
}

# Builds the manifest using podman for the desired architectures
build_manifest() {
    which podman 2>/dev/null || err_exit "Podman not installed, please install podman and proceed";
    local manifest_name="$1"   # $1 == IMAGE MANIFEST NAME
    local architectures="${2:-linux/amd64}"   # $2 == REQUIRED ARCHITECTURES
    local dockerfile_path="${3:-./Dockerfile}" # $3 == PATH TO DOCKERFILE

    [[ $manifest_name ]]   || err_exit "Image manifest name not set"

    podman manifest create ${manifest_name}
    if [ $? -ne 0 ]; then
        err_exit "podman manifest create failed with exit code $?"
    fi

    podman build --platform ${architectures} --manifest localhost/${manifest_name} -f ${dockerfile_path}
    if [ $? -ne 0 ]; then
        err_exit "podman build manifest failed with exit code $?"
    fi

    inspect_manifest ${manifest_name}
}

inspect_manifest(){
    which podman 2>/dev/null || err_exit "Podman not installed, please install podman and proceed";
    local manifest_name="$1"   # $1 == IMAGE MANIFEST NAME
    [[ $manifest_name ]]   || err_exit "Image manifest name not set"

    podman manifest inspect localhost/${manifest_name}:latest

    if [ $? -ne 0 ]; then
        err_exit "podman manifest inspect failed with exit code $?"
    fi
}

push_manifest(){
    which podman 2>/dev/null || err_exit "Podman not installed, please install podman and proceed";
    local manifest_name="$1"   # $1 == IMAGE MANIFEST NAME
    local registry_path="$2"   # $2 == REGISTRY TO WHICH THE IMAGE MUST BE PUSHED
    [[ $manifest_name ]]   || err_exit "Image manifest name not set"
    [[ $manifest_name ]]   || err_exit "Image registry path not set"

   image_registry="$(echo $registry_path | awk -F[/:] '{print $4}')"
   echo "logging into $image_registry" && podman login ${image_registry}

   if [ $? -ne 0 ]; then
        err_exit "podman login into $image_registry failed, please login and continue"
   fi

   image_path="$(echo $registry_path | sed 's/https:\/\///')"
   podman manifest push localhost/$manifest_name:latest $image_path:latest

   if [ $? -ne 0 ]; then
        err_exit "failed to push to $registry_path"
   fi

}

# Print usage information
function print_usage() {

  echo "
USAGE: ./${scriptName} [OPTIONS]
OPTIONS:
"

  echo "
  --build-manifest | -b                                                      : Creates the Java string from all files in a given directory.
  --push-manifest  | -p                                                      : Deletes the Java string files created using --create flag.
  --help                                                                     : This usage information
EXAMPLES:
./${scriptName} -h
"
}

#loader function that parses various options
function evaluate_flags(){
      while [[ "$#" -gt 0 ]]; do
        case $1 in
            '--help') print_usage; exit 0; shift 0;;
            '-b'|'--build-manifest') build_manifest "$2" "$3" "$4"; exit $?;;
            '-p'|'--push-manifest') push_manifest "$2" "$3"; exit $?;;
            '-i'|'--inspect-manifest') inspect_manifest $2; exit $?;;
            *) err_exit "Invalid Option $2"; shift 0;;
        esac
        shift 1
    done

}

# entrypoint function
function main(){
   evaluate_flags "${@}"
}

# Run
main "${@}"
