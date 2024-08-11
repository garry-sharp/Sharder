package reader

import (
	"encoding/hex"
	"fmt"
	"os"

	"github.com/garry-sharp/Sharder/pkg/crypt"
	"github.com/garry-sharp/Sharder/pkg/crypt/alias"
	"github.com/garry-sharp/Sharder/pkg/settings"
	"github.com/manifoldco/promptui"
)

func AddShardPrompt(shards []crypt.ShardT) []crypt.ShardT {
	fmt.Printf("Shards collected %d\n", len(shards))
	prompt := promptui.Select{Label: "What would you like to do: ", Items: []string{"Add Shard", "Assemble Mnemonic", "Exit"}}
	var op int
	var menuErr error

	for {
		op, _, menuErr = prompt.Run()

		if menuErr != nil {
			settings.ErrLog(menuErr)
		}

		switch op {
		case 0:
			idS := ""
			fmt.Printf("Type in your shard id: ")
			_, err := fmt.Scan(&idS)
			if err != nil {
				fmt.Println("Unable to read shard id")
				return AddShardPrompt(shards)
			}

			if idS[0:2] == "0x" {
				idS = idS[2:]
			}

			idA, err := hex.DecodeString(idS)
			if err != nil || len(idA) != 1 {
				fmt.Println("Invalid shard id")
				return AddShardPrompt(shards)
			}

			id := idA[0]

			shardS := ""
			fmt.Printf("Type in your shard: ")
			_, err = fmt.Scan(&shardS)
			if err != nil {
				fmt.Println("Unable to read shard")
				return AddShardPrompt(shards)
			}

			if shardS[0:2] == "0x" {
				shardS = shardS[2:]
			}

			shard, err := hex.DecodeString(shardS)
			if err != nil {
				fmt.Println("Invalid shard string")
				return AddShardPrompt(shards)
			}

			shards = append(shards, crypt.ShardT{Id: id, Data: shard})
			alias := alias.GetAlias(id, shard)

			confirm := promptui.Select{Label: fmt.Sprintf("Alias for shard %s: ", alias), Items: []string{"Yes", "No"}}
			yn, _, _ := confirm.Run()
			if yn == 0 {
				shards = append(shards, crypt.ShardT{Id: id, Data: shard, Alias: alias})
			}

		case 1:
			return shards
		case 2:
			os.Exit(0)
		}
	}

}
