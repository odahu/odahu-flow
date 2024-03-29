#!/usr/bin/env bash
#
# This script will remove an image, but first it will remove all child images of that image.

# Check to see if the specified image exists.
function sanity_check() {
  local IMAGE=$1
  RESULT=$(docker images -aq | grep $IMAGE || true)
  echo "Result: ${RESULT}"
  if test ! "$RESULT"; then
    echo "! -> Sorry, image ${IMAGE} not found!"
    return 1
  fi
}

# Get child images for the supplied argument. Nothing is returned if none found.
get_child_images() {
  local parent_id=$1
  local kids=$(docker inspect --format='ID {{.Id}} PAR {{.Parent}}' $(docker images -aq) |
    grep "PAR ${parent_id}" | sed -E "s/ID ([^ ]*) PAR ([^ ]*)/\1/g")
  echo $kids
}

print_child_images() {
  local parent_id=$1
  local tags=$(docker inspect --format='{{.Id}}' ${parent_id})
  echo "${tags}"

  local children=$(get_child_images "${parent_id}")

  for c in $children; do
    print_child_images "$c"
  done
}

# Check for and remove stopped containers.
function remove_containers() {
  local IMAGE=$1

  # debugging
  echo "remove_containers: $IMAGE"

  echo "# Checking for containers..."
  local CONTAINERS=$(docker ps -aq --filter ancestor=${IMAGE})
  if test "$CONTAINERS"; then
    echo "# Removing the following containers: ${CONTAINERS}"
    docker rm -f ${CONTAINERS}
  fi
}

# Recursively removes all child images.
function remove_child_images() {

  local parent_id=$(docker inspect --format '{{.Id}}' $1)

  echo "# Checking for child images of ${MAIN_IMAGE}...."
  local IMAGES=$(print_child_images ${parent_id} | tail -n+2 | tac)

  if test ! "$IMAGES"; then
    echo "# No child images of ${MAIN_IMAGE}!"
    return 0
  fi
  echo "# -> Found the following child images of ${MAIN_IMAGE}:"
  echo "# ${IMAGES}"

  for IMAGE in $IMAGES; do
    remove_image ${IMAGE}
  done

  echo "# <- Done removing child images of ${MAIN_IMAGE}"
}

# Remove an image.
function remove_image() {
  local IMAGE=$1
  sanity_check $IMAGE

  # stop the function execution if sanity check failed
  if [ "$?" -eq "1" ]; then
    return 1
  fi

  echo -e "\t\t\tSTART remove_containers"
  remove_containers $IMAGE
  echo -e "\t\t\tEND remove_containers"
  echo "# Processing image ${IMAGE}..."
  echo -e "\t\t\tSTART remove_child_images"
  remove_child_images $IMAGE
  echo -e "\t\t\tEND remove_child_images"
  echo "# Remove image ${IMAGE}..."
  docker rmi -f $IMAGE
}

if test ! "$1"; then
  echo "! Syntax: $0 docker_image_id"
fi

MAIN_IMAGE=$1

remove_image $MAIN_IMAGE

echo "# Done removing image ${MAIN_IMAGE} and its children!"
