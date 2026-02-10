package main

import (
	"github.com/spf13/cobra"

	"github.com/smtdfc/contractor/cmd"
)

func main() {

	var rootCmd = &cobra.Command{
		Use:     "contractor",
		Short:   "Type-Safe IDL & Code Generation Toolchain",
		Long:    "Contractor is a specialized Interface Definition Language (IDL) designed to enforce data integrity across distributed systems. It provides a robust mechanism to define cross-platform data contracts, generating validated and idiomatic code for TypeScript and Go, eliminating the risks of manual synchronization.",
		Version: "1.0.0",
	}

	cmd.InitAllCommand(rootCmd)
	rootCmd.Execute()
}
