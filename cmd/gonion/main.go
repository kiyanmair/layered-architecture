package main

import (
	"fmt"
	"layer/config"
	"layer/db"
	"layer/sandbox"
	v0 "layer/versions/v0"
	v1 "layer/versions/v1"
	v2 "layer/versions/v2"
	v3 "layer/versions/v3"
	v4 "layer/versions/v4"
	"log"
	"os"

	"github.com/spf13/cobra"
)

// Map of version numbers to their runner functions
var versionRunners = map[int]func(){
	0: v0.Run,
	1: v1.Run,
	2: v2.Run,
	3: v3.Run,
	4: v4.Run,
}

func main() {
	var rootCmd = &cobra.Command{Use: "gonion"}

	var isSandbox bool
	var version int
	var runCmd = &cobra.Command{
		Use:   "run",
		Short: "Execute a specific version of the program",
		Run: func(cmd *cobra.Command, args []string) {
			// If the sandbox flag is set, run the sandbox
			if isSandbox {
				fmt.Println("Starting sandbox...")
				sandbox.Run()
				return
			}

			// If the version is valid, run the version
			runner, valid := versionRunners[version]
			if valid {
				fmt.Printf("Starting version %d...\n", version)
				runner()
				return
			}

			// If neither is valid or set, print a help message
			fmt.Println("Please provide a valid version or use the sandbox.")
		},
	}

	runCmd.Flags().BoolVarP(&isSandbox, "sandbox", "s", false, "Run the sandbox version")
	runCmd.Flags().IntVarP(&version, "version", "v", -1, "Version to run")

	var dbCmd = &cobra.Command{
		Use:   "db",
		Short: "Database-related operations",
	}

	var resetCmd = &cobra.Command{
		Use:   "reset",
		Short: "Reset the database to its initial state",
		Run: func(cmd *cobra.Command, args []string) {
			db.ResetDB(config.DBPath)
		},
	}

	var showCmd = &cobra.Command{
		Use:   "show",
		Short: "Display the contents of the database",
		Run: func(cmd *cobra.Command, args []string) {
			db.ShowDB(config.DBPath)
		},
	}
	var email string
	var registerCmd = &cobra.Command{
		Use:   "register",
		Short: "Register a new user",
		Run: func(cmd *cobra.Command, args []string) {
			if !ServerRunning() {
				return
			}

			password, err := PromptPassword()
			if err != nil {
				log.Fatalf("Error prompting for password: %v", err)
			}

			err = RegisterUser(email, password)
			if err != nil {
				log.Fatalf("Error registering user: %v", err)
			}

			fmt.Println("User successfully registered!")
		},
	}

	registerCmd.Flags().StringVarP(&email, "email", "e", "", "Email address for registration")
	registerCmd.MarkFlagRequired("email")

	dbCmd.AddCommand(resetCmd, showCmd)
	rootCmd.AddCommand(runCmd, dbCmd, registerCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
