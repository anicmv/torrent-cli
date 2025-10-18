package cmd

import (
	"fmt"
	"os"
	"torrent-cli/pkg/bencode"

	"github.com/spf13/cobra"
)

var infoCmd = &cobra.Command{
	Use:   "info [torrent file]",
	Short: "Display information about a torrent file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := os.ReadFile(args[0])
		if err != nil {
			return fmt.Errorf("failed to read file: %w", err)
		}

		encoder := bencode.NewBEncoder()
		decoded, err := encoder.Decode(data)
		if err != nil {
			return fmt.Errorf("failed to decode torrent: %w", err)
		}

		dict, ok := decoded.(map[string]interface{})
		if !ok {
			return fmt.Errorf("invalid torrent format")
		}

		infoHash, _ := encoder.SHA1(data)

		fmt.Println("Metafile:        ", args[0])
		fmt.Println("InfoHash:         v1:", infoHash)

		if announce, ok := dict["announce"].(string); ok {
			fmt.Println("Announce:        ", announce)
		}

		if info, ok := dict["info"].(map[string]interface{}); ok {
			if name, ok := info["name"].(string); ok {
				fmt.Println("Name:            ", name)
			}
			if pieceLength, ok := info["piece length"].(int64); ok {
				fmt.Println("Piece size:      ", pieceLength)
			}
			if length, ok := info["length"].(int64); ok {
				fmt.Println("Length:          ", length)
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
}
