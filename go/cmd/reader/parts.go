package reader

import "fmt"

func ReadK() (int, error) {
	k := 0
	fmt.Print("Please enter the number of shards required to reconstruct the secret: ")
	_, err := fmt.Scan(&k)
	if err != nil {
		return 0, err
	}
	if k < 2 {
		fmt.Println("The number of shards required to reconstruct the secret must be at least 2")
		return ReadK()
	}
	return k, nil
}

func ReadN(k int) (int, error) {
	n := 0
	fmt.Print("Please enter the total number of shards you'd like to make: ")
	_, err := fmt.Scan(&n)
	if err != nil {
		return 0, err
	}
	if n < k {
		fmt.Println("The total number of shards must be greater than or equal to the number of shards required to reconstruct the mnemonic")
		return ReadN(k)
	}
	return n, nil
}
