package cmd

import (
	"fmt"
	"net/url"
	"os"
	"torrent-cli/pkg/bencode"

	"github.com/spf13/cobra"
)

var magnetCmd = &cobra.Command{
	Use:   "magnet [torrent file]",
	Short: "Generate magnet link from torrent file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := os.ReadFile(args[0])
		if err != nil {
			return fmt.Errorf("failed to read file: %w", err)
		}

		encoder := bencode.NewBEncoder()
		infoHash, err := encoder.SHA1(data)
		if err != nil {
			return fmt.Errorf("failed to calculate info hash: %w", err)
		}

		decoded, err := encoder.Decode(data)
		if err != nil {
			return fmt.Errorf("failed to decode torrent: %w", err)
		}

		dict, ok := decoded.(map[string]interface{})
		if !ok {
			return fmt.Errorf("invalid torrent format")
		}

		info, _ := dict["info"].(map[string]interface{})
		name, _ := info["name"].(string)
		announce, _ := dict["announce"].(string)

		magnetURI := fmt.Sprintf("magnet:?xt=urn:btih:%s&dn=%s&tr=%s",
			infoHash,
			url.QueryEscape(name),
			url.QueryEscape(announce),
		)

		fmt.Println(magnetURI)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(magnetCmd)
}
