package cmd

import (
	"fmt"
	"os"
	"regexp"

	"github.com/garry-sharp/Sharder/pkg/crypt"
	"github.com/garry-sharp/Sharder/pkg/settings"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var asciiArt string

func init() {
	wd, err := os.Getwd()
	if err != nil {
		settings.FatalLog(err)
	}
	f, err := os.ReadFile(wd + "/ascii/ascii.txt")
	if err != nil {
		settings.FatalLog(err)
	}
	asciiArt = string(f)

}

// Add global flags
var verbose bool
var debug bool
var lang string
var wordListDir string

func assembleCmd() *cobra.Command {

	var dir string

	cmd := &cobra.Command{
		Use:   "assemble",
		Short: "Assemble mnemonic from shards",
		Long:  "Assemble mnemonic from shards",
		Run: func(cmd *cobra.Command, args []string) {
			files, err := os.ReadDir(dir)
			if err != nil {
				settings.FatalLog(err)
			}

			regexp, _ := regexp.Compile("shard_([a-zA-Z0-9-]*).json")
			shards := []crypt.ShardT{}
			for _, file := range files {
				wordListLoadAndVerify()
				settings.VerboseLog("Checking file", file.Name())
				if matches := regexp.FindStringSubmatch(file.Name()); len(matches) > 1 {
					settings.VerboseLog("Found file", matches[1])
					d, err := os.ReadFile(dir + "/" + file.Name())
					if err != nil {
						settings.FatalLog(err)
					}
					shard, err := crypt.JSONToShard(d)
					if err != nil {
						settings.FatalLog(err)
					}
					//TODO check alias vs filename
					shards = append(shards, shard)
					settings.StdLog("Loaded shard", shard.Alias)
				}
			}

			mnemonic, err := crypt.Assemble(shards, settings.GetSettings().Lang)
			if err != nil {
				settings.FatalLog(err)
			}
			settings.StdLog("Mnemonic assembled:", mnemonic)
		},
	}

	cmd.Flags().StringVarP(&dir, "dir", "d", ".", "The directory containing the shards")
	return cmd
}

func shardCmd() *cobra.Command {
	var mnemonic string
	var k int
	var n int
	var save bool

	cmd := &cobra.Command{
		Use:   "shard",
		Short: "Generate Shards from mnemonic",
		Long:  "Generate Shards from mnemonic",
		Run: func(cmd *cobra.Command, args []string) {
			wordListLoadAndVerify()
			settings.DebugLog(mnemonic, k, n)
			shards, err := crypt.Shard(mnemonic, k, n, settings.GetSettings().Lang)
			if err != nil {
				settings.FatalLog(err)
			}
			settings.StdLog("Shards Created:")

			tb := tablewriter.NewWriter(os.Stdout)

			tb.SetHeader([]string{"#", "Alias", "Shard ID", "Shard Value"})

			for i, shard := range shards {
				tb.Append([]string{fmt.Sprintf("%d", i+1), shard.Alias, fmt.Sprintf("0x%02x", shard.Id), fmt.Sprintf("0x%0x", shard.Data)})
				//settings.StdLog(fmt.Sprintf("Shard #%d (%s): id - %x, value - %x", i, shard.Alias, shard.Id, shard.Data))
			}

			tb.Render()

			if save {
				settings.StdLog("Saving shards to file")
				for _, shard := range shards {
					shardJson, err := crypt.ShardToJson(shard)
					if err != nil {
						settings.FatalLog(err)
					}
					err = os.WriteFile(fmt.Sprintf("shard_%s.json", shard.Alias), shardJson, 0644)
					if err != nil {
						settings.FatalLog(err)
					}
				}
			}
		},
	}

	cmd.Flags().BoolVar(&save, "save", false, "Save the shards to a file")
	cmd.Flags().StringVarP(&mnemonic, "mnemonic", "m", "", "The mnemonic to be sharded")
	cmd.Flags().IntVarP(&k, "parts", "k", 0, "The minimum number of shards required to reconstruct the mnemonic")
	cmd.Flags().IntVarP(&n, "threshold", "n", 0, "The total number of shards to generate")

	cmd.MarkFlagRequired("mnemonic")
	cmd.MarkFlagRequired("k")
	cmd.MarkFlagRequired("n")

	return cmd
}

func downloadCmd() *cobra.Command {
	var dir string
	cmd := &cobra.Command{
		Use:   "download",
		Short: "Download wordlists",
		Long:  "Download wordlists",
		Run: func(cmd *cobra.Command, args []string) {
			settings.VerboseLog("Downloading wordlists")
			err := crypt.Download(dir)
			if err != nil {
				settings.FatalLog(err)
			}
			settings.StdLog("Wordlists downloaded")
		},
	}

	cmd.Flags().StringVar(&dir, "dir", "$HOME/bip39wordlists", "The location to download the wordlists")
	return cmd
}

func SetupCLI() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "cryptosharder",
		Short: "A crypto mnemonic sharder",
		Long:  asciiArt,
	}

	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "Enable verbose output")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Enable debug output")
	rootCmd.PersistentFlags().StringVar(&lang, "lang", "en", "The language to use for the mnemonic")
	rootCmd.PersistentFlags().StringVar(&wordListDir, "wordlists", "$HOME/bip39wordlists", "The directory containing the wordlists")

	rootCmd.AddCommand(shardCmd())
	rootCmd.AddCommand(assembleCmd())
	rootCmd.AddCommand(downloadCmd())

	rootCmd.PersistentFlags().ParseErrorsWhitelist.UnknownFlags = true
	rootCmd.PersistentFlags().Parse(os.Args)

	settings.GetSettings().Verbose = verbose
	settings.GetSettings().Lang = lang
	settings.GetSettings().Debug = debug
	settings.GetSettings().WordListDir = wordListDir

	settings.VerboseLog("Verbose mode enabled")
	settings.VerboseLog("Language set to", lang)

	settings.DebugLog("Debug mode enabled")

	return rootCmd
}

func setup() {

}
