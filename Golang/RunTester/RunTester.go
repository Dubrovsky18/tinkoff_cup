package RunTester

import (
	"fmt"
	"os/exec"
)

func RunTester(url, name, pathFolder, port1, port2, port3 string) {

	bash := "RunTester/create_sandbox.sh"
	_, err := exec.Command(bash, name, pathFolder, url, port1).Output()
	if err != nil {
		fmt.Printf("error sandBox: %s \n", err)
	}

	start_docker := "RunTester/start_docker.sh"
	cmd1, err := exec.Command(start_docker, name, pathFolder, port1).Output()
	if err != nil {
		fmt.Printf("error docker: %s \n", err)
	}
	output1 := string(cmd1)
	fmt.Printf(output1)
}
