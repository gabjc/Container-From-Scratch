package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

// namespaces

// these will run the same
// docker run <image> <cmd> <params>
// go run main.go run <cmd> <params>

func main() {
	switch os.Args[1] {
	case "run":
		run()
	case "child":
		child()
	default:
		panic("uh oh")
	}
}

func run() {

	fmt.Printf("Running %v as %d\n", os.Args[2:], os.Getpid())

	// Running the executable or command
	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Setting up flags to allow new contianer to be made, and to prevent sharing the filesystem to the host
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags:   syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
		Unshareflags: syscall.CLONE_NEWNS,
	}

	err := cmd.Run()
	if err != nil {
		fmt.Printf("Err: %v", err)
		os.Exit(1)
	}

	// command := os.Args[3]
	// args := os.Args[4:len(os.Args)]
	// cmd := exec.Command(command, args...)

	// cmd.Stdin = os.Stdin
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr

	// err := cmd.Run()
	// if err != nil {
	// 	fmt.Printf("Err: %v", err)
	// 	os.Exit(1)
	// }
}

func child() {

	fmt.Printf("Running %v as %d\n", os.Args[2:], os.Getpid())
	//For setting the container name
	syscall.Sethostname([]byte("container"))

	// Changing the root directory of the contianer
	syscall.Chroot("/MOCK_ROOT")
	syscall.Chdir("/")

	// Mounting the container to a specific folder ("proc") to ultimately isolate the processes to just be those in this filesystem
	syscall.Mount("proc", "proc", "proc", 0, "")

	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Printf("Err: %v", err)
		os.Exit(1)
	}

	syscall.Unmount("/proc", 0)
}
