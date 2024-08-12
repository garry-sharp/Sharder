package cmd

import (
	"fmt"
	"os"
	"regexp"

	"github.com/garry-sharp/Sharder/cmd/reader"
	"github.com/garry-sharp/Sharder/pkg/crypt"
	"github.com/garry-sharp/Sharder/pkg/settings"
	"github.com/manifoldco/promptui"
	"github.com/savioxavier/termlink"

	"embed"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var asciiArt string

//go:embed ascii/ascii.txt
var asciiArtDF embed.FS

func init() {
	wd, err := os.Getwd()
	fmt.Println(wd)
	if err != nil {
		settings.FatalLog(err)
	}
	f, err := asciiArtDF.ReadFile("ascii/ascii.txt")
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

			if lang == "" {
				lang = reader.ReadLang()
				fmt.Println(lang)
			}

			wordListLoadAndVerify()
			settings.GetSettings().Lang = lang

			shards := []crypt.ShardT{}

			if dir != "" {
				files, err := os.ReadDir(dir)
				if err != nil {
					settings.FatalLog(err)
				}
				regexp, _ := regexp.Compile("shard_([a-zA-Z0-9-]*).json")
				for _, file := range files {
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
			} else {
				shards = reader.AddShardPrompt(shards)
			}

			mnemonic, err := crypt.Assemble(shards, settings.GetSettings().Lang)
			if err != nil {
				settings.FatalLog(err)
			}
			fmt.Println("Mnemonic assembled:\n", mnemonic)
		},
	}

	cmd.Flags().StringVarP(&dir, "dir", "d", "", "The directory containing the shards")
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

			if lang == "" {
				lang = reader.ReadLang()
				fmt.Println(lang)
			}

			wordListLoadAndVerify()
			settings.GetSettings().Lang = lang

			if k == 0 {
				var err error
				k, err = reader.ReadK()
				if err != nil {
					settings.ErrLog(err)
					os.Exit(1)
				}
			}

			if n == 0 {
				var err error
				n, err = reader.ReadN(k)
				if err != nil {
					settings.ErrLog(err)
					os.Exit(1)
				}
			}

			if mnemonic == "" {
				var err error
				wordList, err := crypt.GetWordList(lang)
				if err != nil {
					settings.ErrLog(err)
					os.Exit(1)
				}
				mnemonic, err = reader.ReadMnemonic(wordList)
				if err != nil {
					settings.ErrLog(err)
					os.Exit(1)
				}
			}

			//settings.DebugLog(mnemonic, k, n)
			shards, err := crypt.Shard(mnemonic, k, n, settings.GetSettings().Lang)
			if err != nil {
				settings.FatalLog(err)
			}
			settings.StdLog("Shards Created:")

			tb := tablewriter.NewWriter(os.Stdout)

			tb.SetHeader([]string{"#", "Alias", "Shard ID", "Shard Value"})

			for i, shard := range shards {
				tb.Append([]string{fmt.Sprintf("%d", i+1), shard.Alias, fmt.Sprintf("0x%02x", shard.Id), fmt.Sprintf("0x%0x", shard.Data)})
			}

			tb.Render()

			if !save {
				p := promptui.Select{Label: "Would you like to save the shards to a file?", Items: []string{"Yes", "No"}}
				o, _, err := p.Run()
				if err != nil {
					settings.ErrLog(err)
					os.Exit(1)
				}
				if o == 0 {
					save = true
				}
			}

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
	cmd.Flags().IntVarP(&k, "threshold", "k", 0, "The minimum number of shards required to reconstruct the mnemonic")
	cmd.Flags().IntVarP(&n, "parts", "n", 0, "The total number of shards to generate")

	// cmd.MarkFlagRequired("mnemonic")
	// cmd.MarkFlagRequired("k")
	// cmd.MarkFlagRequired("n")

	return cmd
}

func SetupCLI() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "cryptosharder",
		Short: "A crypto mnemonic sharder",
		Long:  asciiArt,
		Run: func(cmd *cobra.Command, args []string) {

			fmt.Println(asciiArt)
			fmt.Println("\nWelcome to Sharder, a crypto mnemonic sharder")
			donate := termlink.ColorLink("DONATING HERE", "http://google.com", "blue")
			fmt.Println("This tool has been made for free, please consider supporting this project by", donate)
			prompt := promptui.Select{Label: "What would you like to do?", Items: []string{"Shard", "Assemble", "Exit"}}
			o, _, err := prompt.Run()
			if err != nil {
				settings.ErrLog(err)
				os.Exit(1)
			}

			switch o {
			case 0:
				shardCmd().Run(cmd, args)
			case 1:
				assembleCmd().Run(cmd, args)
			case 2:
				fmt.Println("Thank you for using Sharder, please consider supporting this project by", donate)
				os.Exit(0)
			}
		},
	}

	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "Enable verbose output")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Enable debug output")
	rootCmd.PersistentFlags().StringVar(&lang, "lang", "", "The language to use for the mnemonic")
	rootCmd.PersistentFlags().StringVar(&wordListDir, "wordlists", "$HOME/bip39wordlists", "The directory containing the wordlists")

	rootCmd.AddCommand(shardCmd())
	rootCmd.AddCommand(assembleCmd())

	rootCmd.PersistentFlags().ParseErrorsWhitelist.UnknownFlags = true
	rootCmd.PersistentFlags().Parse(os.Args)

	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.SetHelpCommand(&cobra.Command{
		Hidden: true,
	})

	settings.GetSettings().Verbose = verbose
	settings.GetSettings().Lang = lang
	settings.GetSettings().Debug = debug

	settings.VerboseLog("Verbose mode enabled")
	settings.VerboseLog("Language set to", lang)

	settings.DebugLog("Debug mode enabled")

	return rootCmd
}

func setup() {

}
