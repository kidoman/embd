#!/usr/bin/env bash

git log --all --format='%cN <%cE>' | sort -u | grep -v karan.misra@gmail.com |\
	grep -v noreply@ | cat CONTRIBUTORS - | sort | uniq > CONTRIBUTORS.new
mv CONTRIBUTORS.new CONTRIBUTORS
