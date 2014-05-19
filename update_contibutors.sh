#!/usr/bin/env bash

git log --all --format='%cN <%cE>' | sort -u | grep -v karan.misra@gmail.com > CONTRIBUTORS
