package p

import (
	"fmt"
	"github.com/logrusorgru/aurora/v4"
)

func Red(info string) {
	fmt.Println(aurora.Red(info))
}

func Green(info string) {
	fmt.Println(aurora.Green(info))
}

func Blue(info string) {
	fmt.Println(aurora.Blue(info))
}

func Cyan(info string) {
	fmt.Println(aurora.Cyan(info))
}

func Yellow(info string) {
	fmt.Println(aurora.Yellow(info))
}

func Hi(s string) {
	fmt.Println(aurora.Bold(s))
}
