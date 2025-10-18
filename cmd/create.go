package cmd

import (
	"fmt"
	"torrent-cli/pkg/torrent"

	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create [file/directory]",
	Short: "Create a torrent file for the specified file or directory",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		options := &torrent.CreateOptions{
			InputPath:   args[0],
			OutputPath:  cmd.Flag("output").Value.String(),
			AnnounceURL: cmd.Flag("announce").Value.String(),
			Comment:     cmd.Flag("comment").Value.String(),
			TorrentName: cmd.Flag("name").Value.String(),
			Source:      cmd.Flag("source").Value.String(),
			CreatedBy:   cmd.Flag("created-by").Value.String(),
			Publisher:   cmd.Flag("publisher").Value.String(),
		}

		privateFlag, _ := cmd.Flags().GetBool("private")
		options.PrivateFlag = privateFlag

		noCreatedBy, _ := cmd.Flags().GetBool("no-created-by")
		options.NoCreatedBy = noCreatedBy

		noCreationDate, _ := cmd.Flags().GetBool("no-creation-date")
		options.NoCreationDate = noCreationDate

		noPublisher, _ := cmd.Flags().GetBool("no-publisher")
		options.NoPublisher = noPublisher

		noSource, _ := cmd.Flags().GetBool("no-source")
		options.NoSource = noSource

		creator := torrent.NewCreator(options)
		if err := creator.Create(); err != nil {
			return fmt.Errorf("failed to create torrent: %w", err)
		}

		fmt.Println("✅ Torrent created successfully!")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	createCmd.Flags().StringP("output", "o", "", "Output file path")
	createCmd.Flags().StringP("announce", "a", "https://example.com", "Announce URL")
	createCmd.Flags().StringP("comment", "c", "torrent-cli", "Comment")
	createCmd.Flags().StringP("name", "n", "", "Torrent name")
	createCmd.Flags().StringP("source", "s", "anicmv :)", "Source tag")
	createCmd.Flags().String("created-by", "torrent-cli", "Created by")
	createCmd.Flags().String("publisher", "anicmv :)", "Publisher")
	createCmd.Flags().BoolP("private", "p", false, "Private flag")
	createCmd.Flags().Bool("no-created-by", false, "Do not include created by")
	createCmd.Flags().Bool("no-creation-date", false, "Do not include creation date")
	createCmd.Flags().Bool("no-publisher", false, "Do not include publisher")
	createCmd.Flags().Bool("no-source", false, "Do not include source")
}
