package cmd

import (
	"context"
	"github.com/678go/xcos/config"
	"github.com/678go/xcos/cos"
	"github.com/spf13/cobra"
)

var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "upload",
	Long:  "upload",
	Run: func(cmd *cobra.Command, args []string) {
		bucket, err := config.InitBucket(Regin)
		if err != nil {
			return
		}
		if err := cos.NewTenCentBucket(bucket).Upload(context.Background(), FileName); err != nil {
			return
		}
	},
}

func init() {
	RootCmd.AddCommand(uploadCmd)
}