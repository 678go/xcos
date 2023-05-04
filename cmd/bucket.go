package cmd

import (
	"context"
	"github.com/678go/xcos/cos"
	"github.com/spf13/cobra"
)

var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "upload",
	Long:  "upload",
	Run: func(cmd *cobra.Command, args []string) {
		bucket, err := cos.NewBucket(Regin)
		if err != nil {
			return
		}
		bucket.InitClient()
		if err := bucket.UploadFolder(context.Background(), FileName); err != nil {
			return
		}
	},
}

var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "download",
	Long:  "download",
	Run: func(cmd *cobra.Command, args []string) {
		bucket, err := cos.NewBucket(Regin)
		if err != nil {
			return
		}
		bucket.InitClient()
		if err := bucket.DownloadFile(context.Background(), FileName); err != nil {
			return
		}
	},
}

func init() {
	RootCmd.AddCommand(uploadCmd)
	RootCmd.AddCommand(downloadCmd)
}
