package cmd

import (
	"clown-core/helpers"
	_ "embed"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

//go:embed files/main.txt
var main []byte

//go:embed files/postgres.txt
var postgres []byte

//go:embed files/builder-create.txt
var builder_create []byte

//go:embed files/builder-update.txt
var builder_update []byte

//go:embed files/builder-delete.txt
var builder_delete []byte

//go:embed files/builder-base.txt
var builder_base []byte

//go:embed files/internal-dotenv.txt
var internal_dotenv []byte

//go:embed files/dotenv.txt
var dotenv []byte

//go:embed files/gitignore.txt
var gitignore []byte

func initProject(name string, mode string) error {
	if mode == "" {
		fmt.Println("no mode selected")
	}

	if mode == "api" {
		err := os.Mkdir(name, 0755)
		if err != nil {
			fmt.Println("Error creating directory:", err)
			return err
		}

		cmd := exec.Command("go", "mod", "init", name)
		cmd.Dir = "./" + name
		err = cmd.Run()
		if err != nil {
			panic(err)
		}

		err = os.Mkdir(name+"/cmd", 0755)
		if err != nil {
			fmt.Println("Error creating directory:", err)
			return err
		}

		err = os.Mkdir(name+"/cmd/api", 0755)
		if err != nil {
			fmt.Println("Error creating directory:", err)
			return err
		}

		err = os.Mkdir(name+"/internal", 0755)
		if err != nil {
			fmt.Println("Error creating directory:", err)
			return err
		}

		err = os.Mkdir(name+"/pkg", 0755)
		if err != nil {
			fmt.Println("Error creating directory:", err)
			return err
		}

		err = os.Mkdir(name+"/internal/repository", 0755)
		if err != nil {
			fmt.Println("Error creating directory:", err)
			return err
		}

		err = os.Mkdir(name+"/internal/repository/builders", 0755)
		if err != nil {
			fmt.Println("Error creating directory:", err)
			return err
		}

		err = os.Mkdir(name+"/internal/config", 0755)
		if err != nil {
			fmt.Println("Error creating directory:", err)
			return err
		}

		err = os.WriteFile(name+"/cmd/api/main.go", helpers.ReplaceProjectName(main, name), 0644)
		if err != nil {
			return err
		}

		err = os.WriteFile(name+"/internal/repository/postgres.go", helpers.ReplaceProjectName(postgres, name), 0644)
		if err != nil {
			return err
		}

		err = os.WriteFile(name+"/internal/repository/builders/base.go", helpers.ReplaceProjectName(builder_base, name), 0644)
		if err != nil {
			return err
		}

		err = os.WriteFile(name+"/internal/repository/builders/create.go", helpers.ReplaceProjectName(builder_create, name), 0644)
		if err != nil {
			return err
		}

		err = os.WriteFile(name+"/internal/repository/builders/update.go", helpers.ReplaceProjectName(builder_update, name), 0644)
		if err != nil {
			return err
		}

		err = os.WriteFile(name+"/internal/repository/builders/delete.go", helpers.ReplaceProjectName(builder_delete, name), 0644)
		if err != nil {
			return err
		}

		err = os.WriteFile(name+"/internal/config/dotenv.go", helpers.ReplaceProjectName(internal_dotenv, name), 0644)
		if err != nil {
			return err
		}

		err = os.WriteFile(name+"/.env", dotenv, 0644)
		if err != nil {
			return err
		}

		err = os.WriteFile(name+"/.gitignore", gitignore, 0644)
		if err != nil {
			return err
		}

		cmdTidy := exec.Command("go", "mod", "tidy")
		cmdTidy.Dir = "./" + name
		cmdTidy.Stdout = os.Stdout
		cmdTidy.Stderr = os.Stderr

		if err := cmdTidy.Run(); err != nil {
			log.Fatal(err)
		}

		fmt.Println("Sucess!")
		return nil
	}

	fmt.Println("This mode not exits.")
	return nil
}

var initCmd = &cobra.Command{
	Use:   "init <name> <mode>",
	Short: "Initializer new project",
	Args:  cobra.ExactArgs(2),

	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		mode := args[1]

		return initProject(name, mode)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
