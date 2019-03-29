// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
	absp "github.com/rhysd/abspath"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "pete",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		input, err := absp.ExpandFrom(viper.GetString("input"))
		if err != nil {
			panic(fmt.Sprintf("error expanding input: %+v", err))
		}
		output, err := absp.ExpandFrom(viper.GetString("output"))
		if err != nil {
			panic(fmt.Sprintf("error expanding output: %+v", err))
		}
		deli := strings.Replace(viper.GetString("deli"), "\\n", "\n", -1)
		linepad := viper.GetString("linepad")
		prefix := viper.GetString("prefix")
		tabsize := "  "

		fmt.Println("input: ", input)
		fmt.Println("output: ", output)
		fmt.Println("deli: ", deli)
		fmt.Println("prefix: ", prefix)

		protofile, queryStart, queryEnd, err := protoFileQueriesPos(output.String())
		if err != nil {
			panic(err)
		}
		// get unformatted queries from pete file
		queries, err := peteQueriesFromFile(input.String(), deli)
		if err != nil {
			panic(err)
		}
		// now format our queries
		decoratePeteQueries(queries, linepad, prefix, tabsize)

		joinedQueries := strings.Join(queries, "")

		data := protofile[0:queryStart] +
			header(linepad) +
			joinedQueries +
			footer(linepad) +
			protofile[queryEnd:]

		if err = ioutil.WriteFile(output.String(), []byte(data), 0644); err != nil {
			panic(err)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// cobra.OnInitialize(newViper)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.pete.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.PersistentFlags().StringP("deli", "d", "\n\n", "the delimiter to use")
	viper.BindPFlag("deli", rootCmd.PersistentFlags().Lookup("deli"))

	rootCmd.Flags().StringP("input", "i", "persist.pete", "file to parse (default is \"persist.pete\"")
	rootCmd.Flags().StringP("output", "o", "", "file to write to")
	rootCmd.Flags().StringP("linepad", "l", "    ", "the padding string for each line defaults to 4 spaces")
	rootCmd.PersistentFlags().StringP("prefix", "p", "", "the package prefix for your in and out types")
	viper.BindPFlag("input", rootCmd.Flags().Lookup("input"))
	viper.BindPFlag("output", rootCmd.Flags().Lookup("output"))
	viper.BindPFlag("linepad", rootCmd.Flags().Lookup("linepad"))
	viper.BindPFlag("prefix", rootCmd.PersistentFlags().Lookup("prefix"))
}

// initConfig reads in config file and ENV variables if set.
func newViper() *viper.Viper {
	snake := viper.New()
	if cfgFile != "" {
		fmt.Println("parsing config flag: ", cfgFile)
		// Use config file from the flag.
		snake.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".pete" (without extension).
		snake.AddConfigPath(home)
		snake.SetConfigName(".pete")
	}

	snake.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := snake.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", snake.ConfigFileUsed())
	}
	return snake
}
