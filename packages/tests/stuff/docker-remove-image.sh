#!/bin/bash
#
# This script will remove an image, but first it will remove all child images of that image.

# Errors are fatal
set -e

#
# ANSI codes for printing out different colors.
#
RED="\033[0;31m"
GREEN="\033[0;32m"
NC="\033[0m"


#
# Check to see if the specified image exists.
# 
# If it does, something will be returned, if not, nothing will be returned.
#
function does_image_exist() {
	local IMAGE=$1
	RESULT=$(docker images -a | grep $IMAGE || true)
	echo $RESULT
}

# Get child images for the supplied argument.  Nothing is returned if none are found
function get_child_images() {
	local IMAGE=$1
	RESULT=$(docker images --filter since=${IMAGE} --quiet 2>/dev/null | sort | uniq | tr '\n' ' ')
	echo $RESULT
}

# Run our sanity checks.
function sanity_check() {

	local IMAGE=$1

	echo "# Verifying that image ${IMAGE} exists..."
	if test ! "$(does_image_exist ${IMAGE})"
	then
		echo "! "
		echo "! Sorry, image ${IMAGE} not found!"
		echo "! "
		return 1
	fi

}

# Check for and remove stopped containers.
function remove_stopped_containers() {

	echo "# Checking for stopped containers..."
	local CONTAINERS=$(docker ps -a | grep Exit | awk '{print $1}' | tr '\n' ' ')
	if test "$CONTAINERS"
	then
		echo "# Removing the following stopped containers: ${CONTAINERS}"
		docker rm ${CONTAINERS}
	fi

}

# Remove all tags from the supplied image.
function remove_tags() {

	local IMAGE=$1

	echo "# Checking image ${IMAGE} for tags that use it..."
	local TAGS=$(docker images | grep $IMAGE | awk '{print $1}' | sort | uniq | tr '\n' ' ')
	if test "${TAGS}"
	then
		echo "# Tags found to remove: ${TAGS}"
		docker rmi ${TAGS}
	fi

}

# Remove an image if it exists.
function remove_image() {

	local IMAGE=$1
	
	printf "#${RED} Processing image ${IMAGE}...${NC}\n"
	remove_child_images $IMAGE

	printf "# ${GREEN}Removing image ${IMAGE}...${NC}\n"
	remove_tags $IMAGE

	if test "$(does_image_exist ${IMAGE})"
	then
		docker rmi $IMAGE
	fi

}

# Recursively remove all child images of an image.
function remove_child_images() {

	local MAIN_IMAGE=$1
	local IMAGE

	echo "# Checking for child images of ${MAIN_IMAGE}...."
	local IMAGES=$(get_child_images ${MAIN_IMAGE})

	if test ! "$IMAGES"
	then
		echo "# No child images of ${MAIN_IMAGE}, bailing out!"
		return 0
	fi

	printf "# ${RED}-> Found the following child images of ${MAIN_IMAGE}: ${NC}\n"
	echo "# ${IMAGES}"

	for IMAGE in $IMAGES
	do
		remove_image $IMAGE
	done

	printf "# ${GREEN}<- Done removing child images of ${MAIN_IMAGE}${NC}\n"

}


if test ! "$1"
then
	echo "! "
	echo "! Syntax: $0 docker_image_id"
	echo "! "
	exit 1
fi

MAIN_IMAGE=$1

sanity_check $MAIN_IMAGE
remove_stopped_containers
remove_image $MAIN_IMAGE

printf "# ${GREEN}Done removing image ${MAIN_IMAGE} and its children!${NC}\n"
