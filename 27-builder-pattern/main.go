/*
 *    _ _
 *  _| | |_  ___    ___  ___  ____                  _____  ___  _____  ___  ____   ___
 * |_     _||_  |  |   ||   ||    \  ___  _ _  ___ |     ||  _||     ||   ||    \ |_  |
 * |_     _| _| |_ | | || | ||  |  || .'|| | ||_ -||  |  ||  _||   --|| | ||  |  ||_  |
 *   |_|_|  |_____||___||___||____/ |__,||_  ||___||_____||_|  |_____||___||____/ |___|
 *                                       |___|
 *
 * - [27/100] Builder Pattern
 */
package main

import "fmt"

var builders map[CookieType]CookieBuilder
var myMachine CookiesMachine

func supplyOrder(order map[CookieType]int) []AnyCookie {
	var result []AnyCookie

	for cookieType, qtt := range order {
		myMachine.SetBuilder(builders[cookieType])
		for i := 0; i < qtt; i++ {
			c := myMachine.MakeCookie()
			if c == nil {
				continue
			}
			result = append(result, c)
		}
	}

	return result
}

func init() {
	builders = make(map[CookieType]CookieBuilder)
	builders[ChocolateChipsCookieType] = new(ChocolateChipsCookiesBuilder)
	builders[PastisetaCookieType] = new(PastisetaCookieBuilder)
}

func main() {
	myMachine = GetMachineInstance()
	order := make(map[CookieType]int)
	order[ChocolateChipsCookieType] = 10
	order[PastisetaCookieType] = 6

	result := supplyOrder(order)

	for _, c := range result {
		fmt.Printf("- %v\n", c.GetDescription())
		fmt.Printf("\tingredientes: %s\n", c.GetIngredients())
		fmt.Printf("\tpeso: %vgr.\n", c.GetWeightGr())
		fmt.Printf("\tcaducidad: %s\n", c.GetExpirationDate().Format("2 Jan 2006"))
	}
}
