#!/bin/bash
set -e

current_dir="${0%/*}"
name=$1

aws cloudformation deploy \
--stack-name $name-deploy \
--template-file "$current_dir/deploy.yaml" \
--parameter-overrides Prefix=$name GithubOAuth=$RENWICK_GITHUB_TOKEN \
--capabilities CAPABILITY_NAMED_IAM