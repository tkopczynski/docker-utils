package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	// TODO: using fork repo for now, pull request pending
	"github.com/tkopczynski/dockerclient"
)

const dockerSocketPath string = "/var/run/docker.sock"
const dockerSocketURL string = "unix:///var/run/docker.sock"

func connectToDocker() (*dockerclient.DockerClient, error) {
	dockerHost := os.Getenv("DOCKER_HOST")
	var tlsConfig tls.Config

	if dockerHost == "" {
		_, err := os.Stat(dockerSocketPath)

		if err != nil {
			log.Fatalf("Could not stat docker unix socket %s: %s\n", dockerSocketPath, err)
		}

		dockerHost = dockerSocketURL
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
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage:")
		fmt.Fprintf(os.Stderr, "%s [OPTIONS] [COMMAND]\n\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "Commands:")
		fmt.Fprintln(os.Stderr, "  rm-all")
		fmt.Fprintln(os.Stderr, "\tremoves all containers, unless running. If used with -force option, removes all containers (including running)")
		fmt.Fprintln(os.Stderr, "  rmi-dangling")
		fmt.Fprintln(os.Stderr, "\tremoves all dangling images (untagged)")
		fmt.Fprintln(os.Stderr, "\nOptions:")

		flag.PrintDefaults()
	}

	forcePtr := flag.Bool("force", false, "Adds force option to command.")

	flag.Parse()

	if flag.NArg() > 1 {
		fmt.Fprintf(os.Stderr, "Too many commands specified, should be only one: %s\n\n", strings.Join(flag.Args(), " "))
		flag.Usage()
		os.Exit(1)
	} else if flag.NArg() == 0 {
		fmt.Fprintln(os.Stderr, "No command specified, there should be at least one")
		fmt.Fprintln(os.Stderr, "")
		flag.Usage()
		os.Exit(1)
	}

	docker, err := connectToDocker()

	if err != nil {
		fmt.Fprintf(os.Stderr, "connect to Docker Engine: %s\n", err)
		os.Exit(1)
	}

	switch flag.Arg(0) {
	case "rm-all":
		removeAllContainers(docker, *forcePtr)
	case "rmi-dangling":
		removeDanglingImages(docker, *forcePtr)
	default:
		fmt.Println("Wrong command specified")
		fmt.Println()
		flag.Usage()
	}

}
