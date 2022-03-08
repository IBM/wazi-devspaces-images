[![Build Status](https://travis-ci.com/IBM/wazi-codeready-workspaces-images.svg?branch=main)](https://travis-ci.com/IBM/wazi-codeready-workspaces-images)
[![Release](https://img.shields.io/github/release/IBM/wazi-codeready-workspaces-images.svg)](../../releases/latest)
[![License](https://img.shields.io/github/license/IBM/wazi-codeready-workspaces-images)](LICENSE)
[![Documentation](https://img.shields.io/badge/Documentation-blue?color=1f618d)](https://ibm.biz/wazi-crw-doc)

## What's inside?

This repository contains the registry, operator, bundle, and index images for IBM Wazi Developer for Red Hat CodeReady Workspaces (IBM Wazi Developer for Workspaces).

* [IBM Wazi Developer for Workspaces - Operator](/codeready-workspaces-operator)
  * This repository is based off of the upstream [Red Hat CodeReady for Workspaces Operator](https://github.com/redhat-developer/codeready-workspaces-operator), whose code is in another upstream [Eclipse Che Operator](https://github.com/eclipse/che-operator/) repository. Repositories gathered together downstream are all brought together into a [Red Hat CodeReady Workspaces Images](https://github.com/redhat-developer/codeready-workspaces-images) repository, where the generated output is located in the respective branch for each version.
* IBM Wazi Developer for Workspaces - [Devfile](/codeready-workspaces-devfileregistry) / [Plug-in](/codeready-workspaces-pluginregistry)
  * You can use the custom devfile and plug-in registry images to tailor your development environment. After you rebuild the registry images with modifications, and push the images into a container registry, you can deploy these images when creating an instance of Wazi Developer for Workspaces or apply the images to an existing instance.
  * _Stacks_ or workspace definitions are pre-configured CodeReady Workspaces with a dedicated set of tools that cover different developer personas. For example, you can pre-configure a workbench for testers with only the tools needed for their purposes. You can add new stacks by customizing the Wazi Developer for Workspaces devfile and plug-in registry.


## IBM Wazi Developer for Workspaces

IBM Wazi Developer for Workspaces is a component shipped with either IBM Wazi Developer for Red Hat CodeReady Workspaces (Wazi Developer for Workspaces) or IBM Developer for z/OS Enterprise Edition (IDzEE).

IBM Wazi Developer for Workspaces provides a modern experience for mainframe software developers working with z/OS applications in the cloud. Powered by the open source projects Zowe and Red Hat CodeReady Workspaces, IBM Wazi Developer for Workspaces offers an easy, streamlined onboarding process to provide mainframe developers with the tools they need. Using container technology, IBM Wazi Developer for Workspaces brings the necessary tools to the task at hand.

For more benefits of IBM Wazi Developer for Workspaces, see the [IBM Wazi Developer product page](https://www.ibm.com/products/wazi-developer) or [IDzEE product page](https://www.ibm.com/products/developer-for-zos).

## Details

IBM Wazi Developer for Workspaces provides a custom stack with the all-in-one mainframe development package that enables mainframe developers to:

- Use a modern mainframe editor with rich language support for COBOL, JCL, Assembler (HLASM), REXX, and PL/I, which provides language-specific features such as syntax highlighting, outline view, declaration hovering, code completion, snippets, a preview of copybooks, copybook navigation, and basic refactoring using IBM Z Open Editor, a component of IBM Wazi Developer for VS Code
- Integrate with any flavor of Git source code management (SCM)
- Perform a user build with IBM Dependency Based Build for any flavor of Git
- Work with z/OS resources such as MVS, UNIX files, and JES jobs
- Connect to the Z host with z/OSMF or IBM Remote System Explorer (RSE) API, using Zowe Explorer plus IBM Z Open Editor for a graphical user interface and Zowe CLI plus the RSE API plug-in for Zowe CLI for command line access
- Debug COBOL and PL/I applications using IBM Z Open Debug
- Use a mainframe development package with a custom plug-in and devfile registry support from the [IBM Wazi Developer stack](https://github.com/IBM/wazi-codeready-workspaces-images)
- Leverage the full language support for Ansible and Red Hat Ansible Certified Content for IBM Z to author and execute playbooks to configure z/OS for development. Automate the deployment and configuration of tools and dependencies such as build, test and debug, to quickly build, deploy and run your applications.

## Documentation

For details of the features for IBM Wazi Developer for Workspaces, see its [official documentation](https://ibm.biz/wazi-crw-doc).

| Repository | Description |
| --- | --- |
| [IBM Wazi Developer for Workspaces](https://github.com/ibm/wazi-codeready-workspaces-images) |  The devfile and plug-in registries and operator |
| [IBM Wazi Developer for Workspaces Sidecars](https://github.com/ibm/wazi-codeready-workspaces-samples) | Supporting resources for the Wazi Developer |

## Feedback
  
We would love to hear feedback from you about IBM Wazi Developer for Workspaces.  
File an issue or provide feedback here: [IBM Wazi Developer for Workspaces Issues](https://github.com/IBM/wazi-codeready-workspaces-images/issues)

---
  
## Reference: Red Hat content

Links marked with this icon :door: are _internal to Red Hat_. This includes Jenkins servers, job configs in gitlab, and container sources in dist-git. 

Because these services are internal, in the interest of making all things open, we've copied as much as possible into this repo. Details below.

## Midstream code
This repo is used to house identical copies of the code used to build the **CodeReady Workspaces images** in Brew/OSBS, but made public to enable pull requests and easier contribution.

* Downstream code can be found in repos _internal to Red Hat_ at http://pkgs.devel.redhat.com/cgit/?q=codeready-workspaces :door:
    - select the `crw-2-rhel-8` branch for the latest `2.x` synced from upstream main branches, or 
    - select a branch like `crw-2.8-rhel-8` for a specific release, synced to a stable branch like `7.28.x`.

## Generated code

In some cases, where we need to house code generated by downstream processes in a public location, this repo will contain folders that end in `-generated` to differentiate from code that's synced from upstream and copied to downstream.

* [codeready-workspaces-operator-metadata-generated/manifests](https://github.com/redhat-developer/codeready-workspaces-images/tree/crw-2-rhel-8/codeready-workspaces-operator-metadata-generated/manifests/) contains CSVs with [pinned digests](http://pkgs.devel.redhat.com/cgit/containers/codeready-workspaces-operator-metadata/tree/container.yaml?h=crw-2-rhel-8#n24) :door:, generated from downstream Brew processes.

## Jenkins jobs

This repo also contains an identical copy of the [Jenkinsfiles and groovy](https://gitlab.cee.redhat.com/codeready-workspaces/crw-jenkins/-/tree/master/jobs/CRW_CI) :door: sources used to configure the [jenkins-csb](https://gitlab.cee.redhat.com/ccit/jenkins-csb) :door: Configuration-as-Code (casc) Jenkins instance used to build the artifacts needed for Brew/OSBS builds. Since the server and config sources are _internal to Red Hat_, [this copy](https://github.com/redhat-developer/codeready-workspaces-images/blob/crw-2-rhel-8/crw-jenkins/jobs/CRW_CI/) is provided to make it easier to see how CodeReady Workspaces is built. Hooray for open source!

* To run a local Jenkins, see [README](https://gitlab.cee.redhat.com/codeready-workspaces/crw-jenkins/-/blob/master/README.md#first-time-user-setup) :door:
* [Job](https://main-jenkins-csb-crwqe.apps.ocp4.prod.psi.redhat.com/job/CRW_CI/job/Releng/job/sync-jenkins-gitlab-to-github_2.x/) :door: that performs the sync from [gitlab](https://gitlab.cee.redhat.com/codeready-workspaces/crw-jenkins/-/blob/master/jobs/CRW_CI/Releng/sync-jenkins-gitlab-to-github.groovy) :door: to [github](https://github.com/redhat-developer/codeready-workspaces-images/blob/crw-2-rhel-8/crw-jenkins/jobs/CRW_CI/Releng/sync-jenkins-gitlab-to-github.groovy) at intervals
* Other jobs are used to:
    * [build a series of artifact and container builds](https://github.com/redhat-developer/codeready-workspaces-images/blob/crw-2-rhel-8/crw-jenkins/jobs/CRW_CI/Releng/build-all-images.groovy),
    * [sync midstream to downstream](https://github.com/redhat-developer/codeready-workspaces-images/blob/crw-2-rhel-8/crw-jenkins/jobs/CRW_CI/crw-sync-to-downstream.groovy),
    * [build artifacts](https://github.com/redhat-developer/codeready-workspaces-images/tree/crw-2-rhel-8/crw-jenkins/jobs/CRW_CI/),
    * [orchestrate Brew builds](https://github.com/redhat-developer/codeready-workspaces-images/blob/crw-2-rhel-8/crw-jenkins/jobs/CRW_CI/get-sources-rhpkg-container-build.groovy),
    * [copy containers to quay](https://github.com/redhat-developer/codeready-workspaces-images/blob/crw-2-rhel-8/crw-jenkins/jobs/CRW_CI/push-latest-container-to-quay.groovy),
    * [check & update digests in registries/metadata images](https://github.com/redhat-developer/codeready-workspaces-images/blob/crw-2-rhel-8/crw-jenkins/jobs/CRW_CI/update-digests-in-registries-and-metadata.groovy)
* Or, to:
    * [send email notifications](https://github.com/redhat-developer/codeready-workspaces-images/blob/crw-2-rhel-8/crw-jenkins/jobs/CRW_CI/Releng/send-email-qe-build-list.groovy) of ER and RC builds
    * [tag sources & collect manifests](https://github.com/redhat-developer/codeready-workspaces-images/blob/crw-2-rhel-8/crw-jenkins/jobs/CRW_CI/Releng/get-3rd-party-deps-manifests.groovy), [collect sources](https://github.com/redhat-developer/codeready-workspaces-images/blob/crw-2-rhel-8/crw-jenkins/jobs/CRW_CI/Releng/get-3rd-party-sources.groovy) to create a release
    * set up subsequent releases ([branching](https://github.com/redhat-developer/codeready-workspaces-images/blob/crw-2-rhel-8/crw-jenkins/jobs/CRW_CI/Releng/create-branches.groovy), [bumping versions](https://github.com/redhat-developer/codeready-workspaces-images/blob/crw-2-rhel-8/crw-jenkins/jobs/CRW_CI/Releng/update-version-and-registry-tags.groovy))
