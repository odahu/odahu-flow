#!/usr/bin/env bash
#
# This script will remove an image, but first it will remove all child images of that image.

# ANSI codes for printing out different colors.
RED="\033[0;31m"
GREEN="\033[0;32m"
NC="\033[0m"

## parent_short_id=$1
## parent_id=`docker inspect --format '{{.Id}}' $1`
##
## get_kids() {
##     local parent_id=$1
##     docker inspect --format='ID {{.Id}} PAR {{.Parent}}' $(docker images -a -q) | grep "PAR ${parent_id}" | sed -E "s/ID ([^ ]*) PAR ([^ ]*)/\1/g"
## }
##
## print_kids() {
##     local parent_id=$1
##     local prefix=$2
##     local tags=`docker inspect --format='{{.RepoTags}}' ${parent_id}`
##     echo "${prefix}${parent_id} ${tags}"
##
##     local children=`get_kids "${parent_id}"`
##
##     for c in $children;
##     do
##         print_kids "$c" "$prefix  "
##     done
## }
##
## print_kids "$parent_id" ""

# Check to see if the specified image exists.
function sanity_check() {
	local IMAGE=$1
	echo "# Verifying that image ${IMAGE} exists..."
	RESULT=$(docker images -a | grep $IMAGE || true)
	echo $RESULT
	if test ! "$RESULT"; then
		echo "! "
		echo "! Sorry, image ${IMAGE} not found!"
		echo "! "
		return 1
	fi
}

# Get child images for the supplied argument. Nothing is returned if none found.
function get_child_images() {
  local PARENT_IMAGE=$1
  for i in $(docker images -q); do
      docker history $i | grep -q ${PARENT_IMAGE} && echo $i
  done | sort -u
}

# Check for and remove stopped containers.
function remove_containers() {
  local IMAGE=$1

  echo "# Checking for containers..."
  local CONTAINERS=$(docker ps -aq ancestor=${IMAGE})
  if test "$CONTAINERS"; then
    echo "# Removing the following containers: ${CONTAINERS}"
    docker rm -f ${CONTAINERS}
  fi
}

# Recursively removes all child images.
function remove_child_images() {

  local MAIN_IMAGE=$1

  echo "# Checking for child images of ${MAIN_IMAGE}...."
  local IMAGES=$(get_child_images ${MAIN_IMAGE})

  if test ! "$IMAGES"; then
    echo "# No child images of ${MAIN_IMAGE}!"
    return 0
  fi

  printf "# ${RED}-> Found the following child images of ${MAIN_IMAGE}: ${NC}\n"
  echo "# ${IMAGES}"

  for IMAGE in $IMAGES; do
    remove_image $IMAGE
  done

  printf "# ${GREEN}<- Done removing child images of ${MAIN_IMAGE}${NC}\n"
}

# Remove an image.
function remove_image() {
  local IMAGE=$1
  remove_containers $IMAGE

  printf "#${RED} Processing image ${IMAGE}...${NC}\n"
  remove_child_images $IMAGE

  docker rmi $IMAGE
}

# if test ! "$1"; then
#   echo "! "
#   echo "! Syntax: $0 docker_image_id"
#   echo "! "
#   exit 1
# fi
#
# MAIN_IMAGE=$1
#
# sanity_check $MAIN_IMAGE
# remove_image $MAIN_IMAGE
#
# printf "# ${GREEN}Done removing image ${MAIN_IMAGE} and its children!${NC}\n"
