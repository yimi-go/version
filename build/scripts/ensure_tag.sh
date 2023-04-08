#!/usr/bin/env bash

# shellcheck disable=SC2153
version="${VERSION}"
if [ "${version}" == "" ];then
  version=v$(gsemver bump)
fi

if [ -z "$(git tag -l "${version}")" ];then
  git tag -a -m "release version ${version}" "${version}"
fi
