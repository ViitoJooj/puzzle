package cmd

import (
	"clown-core/helpers"
	_ "embed"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

//go:embed files/auth-service.txt
var auth_service []byte

//go:embed files/auth-middlewares.txt
var auth_middleware []byte

//go:embed files/auth-pkg-jwt.txt
var auth_jwt []byte

//go:embed files/auth-pkg-response.txt
var auth_responses []byte

//go:embed files/helpers-validator.txt
var helpers_validator []byte

//go:embed files/auth-repository.txt
var auth_repo []byte

func create_add(service string) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	projectName := filepath.Base(wd)

	if service == "auth" {
		authDir := filepath.Join(wd, "internal", "auth")
		middlewareDir := filepath.Join(wd, "internal", "middlewares")

		if err := os.MkdirAll(authDir, 0755); err != nil {
			return err
		}

		if err := os.MkdirAll(middlewareDir, 0755); err != nil {
			return err
		}

		files := map[string][]byte{
			filepath.Join(authDir, "service.go"):     helpers.ReplaceProjectName(auth_service, projectName),
			filepath.Join(authDir, "repository.go"):  helpers.ReplaceProjectName(auth_repo, projectName),
			filepath.Join(middlewareDir, "auth.go"):  helpers.ReplaceProjectName(auth_middleware, projectName),
			filepath.Join(wd, "pkg", "jwt.go"):       helpers.ReplaceProjectName(auth_jwt, projectName),
			filepath.Join(wd, "pkg", "responses.go"): helpers.ReplaceProjectName(auth_responses, projectName),
		}

		for path, content := range files {
			if err := os.WriteFile(path, content, 0644); err != nil {
				return err
			}
		}

		cmdTidy := exec.Command("go", "mod", "tidy")
		cmdTidy.Stdout = os.Stdout
		cmdTidy.Stderr = os.Stderr
		if err := cmdTidy.Run(); err != nil {
			log.Fatal(err)
		}

		fmt.Println("Sucess!")
		return nil
	}

	if service == "hash" {
		authDir := filepath.Join(wd, "internal", "helpers")

		if err := os.MkdirAll(authDir, 0755); err != nil {
			return err
		}

		files := map[string][]byte{
			filepath.Join(authDir, "hashPassword.go"): helpers.ReplaceProjectName(helpers_validator, projectName),
		}

		for path, content := range files {
			if err := os.WriteFile(path, content, 0644); err != nil {
				return err
			}
		}

		cmdTidy := exec.Command("go", "mod", "tidy")
		cmdTidy.Stdout = os.Stdout
		cmdTidy.Stderr = os.Stderr
		if err := cmdTidy.Run(); err != nil {
			log.Fatal(err)
		}

		fmt.Println("Sucess!")
		return nil
	}

	fmt.Println("This service not exits.")
	return nil
}

var add = &cobra.Command{
	Use:   "add <service>",
	Short: "create a service",
	Args:  cobra.ExactArgs(1),

	RunE: func(cmd *cobra.Command, args []string) error {
		service := args[0]
		return create_add(service)
	},
}

func init() {
	rootCmd.AddCommand(add)
}
