// Copyright (c) updateip author(s) 2022. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package cmd

import (
	"fmt"

	"github.com/morikuni/aec"
	"github.com/spf13/cobra"
)

var (
	Version   string
	GitCommit string
)

func PrintupdateipASCIIArt() {
	updateipLogo := aec.BlueF.Apply(updateipFigletStr)
	fmt.Print(updateipLogo)
}

func MakeVersion() *cobra.Command {
	var command = &cobra.Command{
		Use:          "version",
		Short:        "Print the version",
		Example:      `  updateip version`,
		Aliases:      []string{"v"},
		SilenceUsage: false,
	}
	command.Run = func(cmd *cobra.Command, args []string) {
		PrintupdateipASCIIArt()
		if len(Version) == 0 {
			fmt.Println("Version: dev")
		} else {
			fmt.Println("Version:", Version)
		}
		fmt.Println("Git Commit:", GitCommit)

	}
	return command
}

const updateipFigletStr = `                 _       _       _       
 _   _ _ __   __| | __ _| |_ ___(_)_ __  
| | | | '_ \ / _` + "`" + ` |/ _` + "`" + ` | __/ _ \ | '_ \ 
| |_| | |_) | (_| | (_| | ||  __/ | |_) |
 \__,_| .__/ \__,_|\__,_|\__\___|_| .__/ 
      |_|                         |_|    

Update your dynamic DNS records	

`
