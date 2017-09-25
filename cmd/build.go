package cmd

import (
	"errors"
	"strings"

	"github.com/buddyspike/mbt/lib"
	"github.com/spf13/cobra"
)

var pipedArgs []string

func init() {
	buildPr.Flags().StringVar(&source, "source", "", "source branch")
	buildPr.Flags().StringVar(&dest, "dest", "", "destination branch")

	buildCommand.PersistentFlags().StringArrayVarP(&pipedArgs, "args", "a", []string{}, "arguments to be passed into build scripts")
	buildCommand.AddCommand(buildBranch)
	buildCommand.AddCommand(buildPr)
	RootCmd.AddCommand(buildCommand)
}

func preparePipedArgs() []string {
	a := []string{}
	for _, i := range pipedArgs {
		if strings.Contains(i, "=") {
			k := strings.Split(i, "=")
			a = append(a, k...)
		} else {
			a = append(a, i)
		}
	}
	return a
}

var buildBranch = &cobra.Command{
	Use:   "branch <branch>",
	Short: "builds the specific branch",
	RunE: func(cmd *cobra.Command, args []string) error {
		branch := "master"
		if len(args) > 0 {
			branch = args[0]
		}

		m, err := lib.ManifestByBranch(in, branch)
		if err != nil {
			return err
		}

		err = lib.Build(m, preparePipedArgs())
		return err
	},
}

var buildPr = &cobra.Command{
	Use:   "pr",
	Short: "builds the pr from a source branch to destination branch",
	RunE: func(cmd *cobra.Command, args []string) error {
		if source == "" {
			return errors.New("requires source")
		}

		if dest == "" {
			return errors.New("requires dest")
		}

		m, err := lib.ManifestByPr(in, source, dest)
		if err != nil {
			return err
		}

		err = lib.Build(m, preparePipedArgs())
		return err
	},
}

var buildCommand = &cobra.Command{
	Use:   "build",
	Short: "Builds the applications in specified path",
}