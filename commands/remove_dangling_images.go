package commands

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/tkopczynski/dockerclient"
)

func createSingleValueFilter(name string, value string) (string, error) {
	filters := make(map[string][]string)

	val := make([]string, 1)
	val[0] = value

	filters[name] = val

	jsonFilters, err := json.Marshal(filters)

	return string(jsonFilters), err
}

// RemoveDanglingImages removes dangling (untagged) images from Docker.
// It can be used with additional force flag to remove those currently used as well.
func RemoveDanglingImages(docker *dockerclient.DockerClient, force bool) {
	jsonFilters, err := createSingleValueFilter("dangling", "true")

	if err != nil {
		log.Fatalf("Encoding JSON: %s\n", err)
	}

	images, err := docker.ListImages(true, jsonFilters)

	if err != nil {
		log.Fatalf("Listing images: %s\n", err)
	}

	for _, image := range images {
		fmt.Println(image.Id)
		docker.RemoveImage(image.Id, force)
	}
}
