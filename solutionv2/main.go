package main

import (
	"fmt"
	"os/exec"
)

func main() {

	team_name := "abosdfdffba"
	avaiable_port := "12344"
	url := "https://www.tinkoff.ru" //без разницы как получать эти значения, самое главное, чтобы они были

	bash := "./create_sandbox.sh"
	cmd, err := exec.Command(bash, team_name, avaiable_port, url).Output()
	if err != nil {
		fmt.Printf("error %s", err)
	}
	output := string(cmd)
	fmt.Printf(output)

	start_docker := "./start_docker.sh"
	cmd1, err := exec.Command(start_docker, team_name, avaiable_port).Output()
	if err != nil {
		fmt.Printf("error %s", err)
	}
	output1 := string(cmd1)
	fmt.Printf(output1)

}
