package Tests

import (
	"fmt"
	"os/exec"
)

func RunTester(userFolder, fileName, url, user string) string {

	//avaiable_port := "12344"
	port1 := "12345"
	port2 := "12346"
	port3 := "12347"
	bash := "./Tests/create_sandbox.sh"

	_, err := exec.Command(bash, userFolder, fileName, url, user, port1, port2, port3).Output()
	if err != nil {
		fmt.Printf("error %s", err)
	}

	start_docker := "./Tests/start_docker.sh"
	cmd, err := exec.Command(start_docker, userFolder, user, fileName).Output()
	if err != nil {
		fmt.Printf("error %s", err)
	}
	output := string(cmd)
	return output

}
