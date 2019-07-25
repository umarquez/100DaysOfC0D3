/*
   _ _
 _| | |_  ___    ___  ___  ____                  _____  ___  _____  ___  ____   ___
|_     _||_  |  |   ||   ||    \  ___  _ _  ___ |     ||  _||     ||   ||    \ |_  |
|_     _| _| |_ | | || | ||  |  || .'|| | ||_ -||  |  ||  _||   --|| | ||  |  ||_  |
  |_|_|  |_____||___||___||____/ |__,||_  ||___||_____||_|  |_____||___||____/ |___|
                                      |___|
- [17/100] Endinanness | Big-endian & Little-endian
*/

package main

import (
	"encoding/binary"
	"fmt"
	"strings"
)

func main() {
	bContent := []byte("abcdwxyz")
	for _, hx := range bContent {
		fmt.Printf("%X, ", hx)
	}
	fmt.Printf("\n%v\n", strings.Join(strings.Split(string(bContent), ""), "\t"))

	be := binary.BigEndian.Uint64(bContent)
	le := binary.LittleEndian.Uint64(bContent)
	fmt.Printf("%X\n", be)
	fmt.Printf("%X\n", le)
}
