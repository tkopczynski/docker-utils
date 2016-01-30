package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	// TODO: using fork repo for now, pull request pending
	"github.com/tkopczynski/dockerclient"
)

const dockerSocketPath string = "/var/run/docker.sock"
const dockerSocketUrl string = "unix:///var/run/docker.sock"

func connectToDocker() (*dockerclient.DockerClient, error) {
	dockerHost := os.Getenv("DOCKER_HOST")
	var tlsConfig tls.Config

	if dockerHost == "" {
		_, err := os.Stat(dockerSocketPath)

		if err != nil {
			log.Fatalf("Could not stat docker unix socket %s: %s\n", dockerSocketPath, err)
		}

		dockerHost = dockerSocketUrl
	} else {
		certPath := os.Getenv("DOCKER_CERT_PATH")

		if certPath == "" {
			log.Fatalf("DOCKER_CERT_PATH should be set when DOCKER_HOST is used. DOCKER_HOST = %s\n", dockerHost)
		}

		cert, err := tls.LoadX509KeyPair(certPath+"/cert.pem", certPath+"/key.pem")

		if err != nil {
			log.Fatalf("loadkeys: %s\n", err)
		}

		tlsConfig = tls.Config{Certificates: []tls.Certificate{cert}, InsecureSkipVerify: true}
	}

	return dockerclient.NewDockerClient(dockerHost, &tlsConfig)

}

func createSingleValueFilter(name string, value string) (string, error) {
	filters := make(map[string][]string)

	val := make([]string, 1)
	val[0] = value

	filters[name] = val

	jsonFilters, err := json.Marshal(filters)

	return string(jsonFilters), err
}

func removeDanglingImages(docker *dockerclient.DockerClient, force bool) {
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

func removeAllContainers(docker *dockerclient.DockerClient, force bool) {
	containers, err := docker.ListContainers(true, force, "")

	if err != nil {
		log.Fatalf("Listing containers: %s\n", err)
	}

	for _, container := range containers {
		fmt.Println(container.Id)
		docker.RemoveContainer(container.Id, force, false)
	}
}

func main() {
	rmAllPtr := flag.Bool("rm-all", false, "Command: removes all containers, unless it's running. If used with -force option, removes all containers (including running)")
	rmiDangling := flag.Bool("rmi-dangling", false, "Command: removes all dangling images (untagged)")
	forcePtr := flag.Bool("force", false, "Adds force option to command.")

	flag.Parse()

	docker, err := connectToDocker()

	if err != nil {
		log.Fatalf("connect to Docker Engine: %s\n", err)
	}

	switch {
	case *rmAllPtr == true:
		removeAllContainers(docker, *forcePtr)
	case *rmiDangling == true:
		removeDanglingImages(docker, *forcePtr)
	default:
		fmt.Println("Nothing to do, please specify command.")
		fmt.Println()
		flag.Usage()
	}

}