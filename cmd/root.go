package cmd

import "github.com/spf13/cobra"

var (
	Regin    string
	FileName string
)

var RootCmd = &cobra.Command{
	Use:   "xcos",
	Short: "cos tools",
	Long:  "cos tools",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
		return
	},
}

func init() {
	RootCmd.PersistentFlags().StringVarP(&Regin, "regin", "r", "", "cos regin")
	RootCmd.PersistentFlags().StringVarP(&FileName, "filename", "f", "", "filename")
}
