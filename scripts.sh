#!/usr/bin/env bash

git_tag() {
    echo $(git rev-list ${1:-master} --abbrev-commit --max-count=1)
}

set_image_tag() {
    # create .env file if it doesn't exist already
    ! test -f .env && touch .env
    grep -q '^IMAGE_TAG=' .env \
	&& sed  -i .bak 's@\(IMG_TAG=\)@\1'"$(git_tag)"'@g' \
		.env || echo "IMG_TAG=$(git_tag)" >> .env
}

FUNC=${1:-git_tag}
shift
eval "$FUNC $@"
