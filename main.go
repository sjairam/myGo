package main

// Try Go with all databases

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Try add common usage function
func usage() {
	fmt.Println(" This program utilises Go.  It shall also doa a pre-check for Go")
	fmt.Println(" The current usage is incorrect.  Please see example input parameters below")
	fmt.Println(" Please provider the database engine type, options are 'mysql' or 'oracle-ee' or 'postgres'.")
	os.Exit(1)
}

func main() {
	//Check GO installed
	goPath, err := exec.LookPath("go")
	if err != nil {
		fmt.Println("ERROR: The Go binary does not exist or is not executable.")
		fmt.Println("Please install Go and ensure it is in your PATH.")
		os.Exit(1)
	}
	fmt.Printf("Go found at: %s\n", goPath)

	if len(os.Args) < 2 {
		usage()
	}

	engine := os.Args[1]
	if engine != "postgres" && engine != "mysql" && engine != "oracle-ee" {
		usage()
		os.Exit(1)
	}

	// Check if AWS CLI is installed
	awsPath, err := exec.LookPath("aws")
	if err != nil {
		fmt.Println("ERROR: The aws binary does not exist or is not executable.")
		fmt.Println("Please install AWS CLI and ensure it is in your PATH.")
		os.Exit(1)
	}

	// Prepare the command to fetch all RDS instances with specified database engine
	cmd := exec.Command(awsPath, "rds", "describe-db-instances",
		"--query", fmt.Sprintf("DBInstances[?Engine=='%s'].[DBInstanceIdentifier, DBInstanceClass, Engine, EngineVersion, DBInstanceStatus, MultiAZ]", engine),
		"--output", "table")

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr

	// Execute the command and capture its output
	err = cmd.Run()
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		os.Exit(1)
	}

	rdsInstances := strings.TrimSpace(out.String())

	// Check if any RDS instances of specified engine are found
	if rdsInstances == "" {
		fmt.Printf("No %s RDS instances found.\n", engine)
	} else {
		fmt.Printf("%s RDS Instances:\n", strings.Title(engine))
		fmt.Println(rdsInstances)
	}
}
