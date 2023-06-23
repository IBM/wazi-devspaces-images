[![Build](https://app.travis-ci.com/IBM/wazi-devspaces-images.svg?branch=3.0.0.wazi)](https://app.travis-ci.com/IBM/wazi-devspaces-images)
[![Release](https://img.shields.io/github/release/IBM/wazi-devspaces-images.svg)](../../releases/latest)
[![License](https://img.shields.io/github/license/IBM/wazi-devspaces-images)](./LICENSE)
[![Documentation](https://img.shields.io/badge/Documentation-blue?color=1f618d)](https://ibm.biz/wazi-ds-doc)

## What's inside?

This repository contains the registry, operator, bundle, and index images for IBM Wazi for Dev Spaces.

* [IBM Wazi for Dev Spaces - Operator](/devspaces-operator)
  * This repository is based off of the upstream [Red Hat Dev Spaces Images](https://github.com/redhat-developer/devspaces-images) repository, where the generated output is located in the respective branch for each version.
* [IBM Wazi for Dev Spaces - Dashboard](/devspaces-dashboard)
  * The dashboard page is your landing page that has been customized to provide a more inclusive user experience.
* [IBM Wazi for Dev Spaces - Devfile](/devspaces-devfileregistry) / [IBM Wazi for Dev Spaces - Plug-in](/devspaces-pluginregistry)
  * You can use the custom devfile and plug-in registry images to tailor your development environment.
  * _Stacks_ or workspace definitions are pre-configured Workspaces with a dedicated set of tools that cover different developer personas. For example, you can pre-configure a workbench for testers with only the tools needed for their purposes.


## IBM Wazi for Dev Spaces

IBM Wazi for Dev Spaces provides a modern experience for mainframe software developers working with z/OS applications in the Hybrid Cloud. Powered by the open source projects Zowe and Red Hat OpenShift Dev Spaces, IBM Wazi for Dev Spaces offers an easy, streamlined onboarding process to provide mainframe developers with the tools they need. Using container technology, IBM Wazi for Dev Spaces brings the necessary tools to the task at hand.

## Details

IBM Wazi for Dev Spaces provides a custom stack with the all-in-one mainframe development package that enables mainframe developers to:

- Use a modern mainframe editor with rich language support for COBOL, PL/I, Assembler (HLASM), REXX, and JCL, which provides language-specific features such as syntax highlighting, outline view, declaration hovering, code completion, snippets, a preview of copybooks, copybook navigation, and basic refactoring using IBM Z Open Editor, a component of IBM Wazi for VS Code.
- Integrate with any flavor of Git source code management (SCM).
- Perform a user build on a remote z/OS system with IBM Dependency Based Build.
- Work with z/OS resources such as MVS data sets, UNIX files, and JES jobs.
- Connect to the z/OS host with z/OSMF or IBM Remote System Explorer (RSE) API, using Zowe Explorer for a graphical user interface and Zowe CLI utilizing the RSE API plug-in for Zowe CLI for command line access.
- Debug COBOL, PL/I, and HLASM applications using IBM Z Open Debug.
- Leverage the full language support for Ansible and Red Hat Ansible Certified Content for IBM Z to author and execute playbooks to configure z/OS for development.
- Use IBM Wazi Analyze to discover the relationships among z/OS application artifacts.

## Documentation

For details of the features for IBM Wazi for Dev Spaces, see its [official documentation](https://ibm.biz/wazi-ds-doc).

### Image Information

* `catalog` Image - This is an index image used for the operator hub. 
* `operator` Image - This is the operator image inherited by Red Hat. 
* `dashboard` Image - This is a customized dashboard image. 
* `devfile` Image - This is a customized devfile registry image. 
* `plugin` Image - This is a customized plugin registry image. 
* `sidecar` Image - This is a sidecar image used to provide supporting resources. 

| <sub>Registry and Image</sub> | <sub>Version</sub> |
| :--- | --- |
||<sub><b>3.0.0</b></sub>|
|<sub>icr.io/cpopen/ibm-wazi-for-devspaces-catalog@sha256:13ba2dc729012286d5d72b4874cc77d7afb58205ac7b92ccc20dfc385224af0a</sub>||
|<sub>icr.io/cpopen/ibm-wazi-for-devspaces-operator@sha256:025f8829fe1b4e401ca2d121785db08f55ea6157c0826d51430714101cc1d9c9</sub>||
|<sub>icr.io/wazi-code/ibm-wazi-for-devspaces-dashboard@sha256:410ee4b654acf8fa943da520abe11d8a5138987d1edf58f0e89baba8f9a5867c</sub>||
|<sub>icr.io/wazi-code/ibm-wazi-for-devspaces-devfile@sha256:4a636e886241f802c04c26161708a0e2a6950d744c3e7dfdad81e12b87e2fd1f</sub>||
|<sub>icr.io/wazi-code/ibm-wazi-for-devspaces-plugin@sha256:0397cf795f10f2ecd5f11cc1b9b851e50361c932fe9b620e9f83f08594b65023</sub>||
|<sub>icr.io/wazi-code/ibm-wazi-for-devspaces-sidecar@sha256:6832c15f01f91af86cc5824ce2cc6bf3530b1408c5e53da86ef344ae34932d55</sub>||

## Feedback
  
We would love to hear feedback from you about IBM Wazi for Dev Spaces.  
File an issue or provide feedback here: [IBM Wazi for Dev Spaces Issues](https://github.com/IBM/wazi-devspaces-images/issues)

_Red Hat Content_

Links marked with this icon :door: are _internal to Red Hat_. This includes Jenkins servers, job configs in gitlab, and container sources in dist-git. 

Because these services are internal, in the interest of making all things open, we've copied as much as possible into this repo. Details below.

## Midstream code
This repo is used to house identical copies of the code used to build the **Red Hat OpenShift Dev Spaces (formerly CodeReady Workspaces) images** in Brew/OSBS, but made public to enable pull requests and easier contribution.

* Downstream code can be found in repos _internal to Red Hat_ at http://pkgs.devel.redhat.com/cgit/?q=devspaces :door:
    - select the `devspaces-3-rhel-8` branch for the latest `2.x` synced from upstream main branches, or 
    - select a branch like `devspaces-3.1-rhel-8` for a specific release, synced to a stable branch like `7.46.x`.

## Generated code

In some cases, where we need to house code generated by downstream processes in a public location, this repo will contain folders that end in `-generated` to differentiate from code that's synced from upstream and copied to downstream.

* [devspaces-operator-bundle-generated/manifests](https://github.com/redhat-developer/devspaces-images/tree/devspaces-3-rhel-8/devspaces-operator-bundle-generated/manifests/) contains CSVs with [pinned digests](http://pkgs.devel.redhat.com/cgit/containers/devspaces-operator-bundle/tree/container.yaml?h=devspaces-3-rhel-8#n24) :door:, generated from downstream Brew processes.

## Jenkins jobs

This repo also contains an identical copy of the [Jenkinsfiles and groovy](https://gitlab.cee.redhat.com/codeready-workspaces/crw-jenkins/-/tree/master/jobs/DS_CI) :door: sources used to configure the [jenkins-csb](https://gitlab.cee.redhat.com/ccit/jenkins-csb) :door: Configuration-as-Code (casc) Jenkins instance used to build the artifacts needed for Brew/OSBS builds. Since the server and config sources are _internal to Red Hat_, [this copy](https://github.com/redhat-developer/devspaces-images/blob/devspaces-3-rhel-8/crw-jenkins/jobs/DS_CI/) is provided to make it easier to see how Red Hat OpenShift Dev Spaces is built. Hooray for open source!

* To run a local Jenkins, see [README](https://gitlab.cee.redhat.com/codeready-workspaces/crw-jenkins/-/blob/master/README.md#first-time-user-setup) :door:
* [Job](https://main-jenkins-csb-crwqe.apps.ocp-c1.prod.psi.redhat.com/job/DS_CI/job/Releng/job/sync-jenkins-gitlab-to-github_2.x/) :door: that performs the sync from [gitlab](https://gitlab.cee.redhat.com/codeready-workspaces/crw-jenkins/-/blob/master/jobs/DS_CI/Releng/sync-jenkins-gitlab-to-github.groovy) :door: to [github](https://github.com/redhat-developer/devspaces-images/blob/devspaces-3-rhel-8/crw-jenkins/jobs/DS_CI/Releng/sync-jenkins-gitlab-to-github.groovy) at intervals
* Other jobs are used to:
    * [build a series of artifact and container builds](https://github.com/redhat-developer/devspaces-images/blob/devspaces-3-rhel-8/crw-jenkins/jobs/DS_CI/Releng/build-all-images.groovy),
    * [sync midstream to downstream](https://github.com/redhat-developer/devspaces-images/blob/devspaces-3-rhel-8/crw-jenkins/jobs/DS_CI/crw-sync-to-downstream.groovy),
    * [build artifacts](https://github.com/redhat-developer/devspaces-images/tree/devspaces-3-rhel-8/crw-jenkins/jobs/DS_CI/),
    * [orchestrate Brew builds](https://github.com/redhat-developer/devspaces-images/blob/devspaces-3-rhel-8/crw-jenkins/jobs/DS_CI/get-sources-rhpkg-container-build.groovy),
    * [copy containers to quay](https://github.com/redhat-developer/devspaces-images/blob/devspaces-3-rhel-8/crw-jenkins/jobs/DS_CI/push-latest-container-to-quay.groovy),
    * [check & update digests in registries/metadata images](https://github.com/redhat-developer/devspaces-images/blob/devspaces-3-rhel-8/crw-jenkins/jobs/DS_CI/update-digests-in-registries-and-metadata.groovy)
* Or, to:
    * [send email notifications](https://github.com/redhat-developer/devspaces-images/blob/devspaces-3-rhel-8/crw-jenkins/jobs/DS_CI/Releng/send-email-qe-build-list.groovy) of ER and RC builds
    * [tag sources & collect manifests](https://github.com/redhat-developer/devspaces-images/blob/devspaces-3-rhel-8/crw-jenkins/jobs/DS_CI/Releng/get-3rd-party-deps-manifests.groovy), [collect sources](https://github.com/redhat-developer/devspaces-images/blob/devspaces-3-rhel-8/crw-jenkins/jobs/DS_CI/Releng/get-3rd-party-sources.groovy) to create a release
    * set up subsequent releases ([branching](https://github.com/redhat-developer/devspaces-images/blob/devspaces-3-rhel-8/crw-jenkins/jobs/DS_CI/Releng/create-branches.groovy), [bumping versions](https://github.com/redhat-developer/devspaces-images/blob/devspaces-3-rhel-8/crw-jenkins/jobs/DS_CI/Releng/update-version-and-registry-tags.groovy))
