package commands

import (
	"fmt"
	"log"

	"github.com/tkopczynski/dockerclient"
)

// RemoveAllContainers removes all containers in Docker Engine.
// It can be used with additional force flag to remove running containers as well.
func RemoveAllContainers(docker *dockerclient.DockerClient, force bool) {
	containers, err := docker.ListContainers(true, force, "")

	if err != nil {
		log.Fatalf("Listing containers: %s\n", err)
	}

	for _, container := range containers {
		fmt.Println(container.Id)
		docker.RemoveContainer(container.Id, force, false)
	}
}
