package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	// TODO: using fork repo for now, pull request pending
	"github.com/tkopczynski/docker-utils/commands"
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
		commands.RemoveAllContainers(docker, *forcePtr)
	case "rmi-dangling":
		commands.RemoveDanglingImages(docker, *forcePtr)
	default:
		fmt.Println("Wrong command specified")
		fmt.Println()
		flag.Usage()
	}

}
