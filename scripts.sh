#!/usr/bin/env bash

git_tag() {
    echo $(git rev-list ${1:-master} --abbrev-commit --max-count=1)
}

set_image_tag() {
    # create .env file if it doesn't exist already
    ENV_FILE=.env
    ENV_NAME=IMG_TAG
    ! test -f $ENV_FILE && touch $ENV_FILE
    grep -q "^$ENV_NAME=" .env \
	&& sed  -i .bak 's@\(IMG_TAG=\)\(.*\)@\1'"$(git_tag)"'@g' \
		.env || echo "IMG_TAG=$(git_tag)" >> .env
}

POSITIONAL=()
while [[ $# -gt 0 ]]
do
key="$1"

case $key in
    -s|--submodule)
    SUBMODULE=YES
    shift # past argument
    ;;
    *)    # unknown option
    POSITIONAL+=("$1") # save it in an array for later
    shift # past argument
    ;;
esac
done
set -- "${POSITIONAL[@]}"

if test $SUBMODULE ; then
    cd todo-react
fi

FUNC=${1:-set_image_tag}
shift
eval "$FUNC $@"
