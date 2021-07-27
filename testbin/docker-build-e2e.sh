#!/bin/sh

#  Copyright 2021 The UnDistro Authors.
#
#  Licensed under the Apache License, Version 2.0 (the "License");
#  you may not use this file except in compliance with the License.
#  You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
#  Unless required by applicable law or agreed to in writing, software
#  distributed under the License is distributed on an "AS IS" BASIS,
#  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#  See the License for the specific language governing permissions and
#  limitations under the License.

set -o errexit
set -o pipefail

function exit_and_inform {
	err_n=$1
	case $err_n in
		0)
			echo "Error: No git project." 1>&2
			;;
		1)
			echo "Error: No tag provided." 1>&2
			;;
		2)
			echo "Error: No host provided." 1>&2
			;;
		*)
			echo "Usage: $(basename $0) <docker_tag> <registry_host>" 1>&2
			;;
	esac
	exit 1
}

function make_and_push {
	make manager;
	mv ./bin/manager .;
	docker build -t $host:5000/undistro:$tag .;
	docker push $host:5000/undistro:$tag;
	make aws-init;
	mv ./bin/aws-init .;
	docker build -t $host:5000/aws-init:$tag -f aws-init.docker .;
	docker push $host:5000/aws-init:$tag;
}


if test $# -ne 2; then
	exit_and_inform
fi
tag=$1
host=$2
proj_root=$(git rev-parse --show-toplevel)

if test -n "$proj_root"; then
	if test -z "$tag"; then
		exit_and_inform 1
	elif test -z "$host"; then
		exit_and_inform 2
	fi
	cd "$proj_root"
	make_and_push
else
	exit_and_inform 0
fi
