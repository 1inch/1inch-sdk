package amounts

import (
	"github.com/1inch/1inch-sdk-go/internal/helpers"
)

const (
	Ten5  = "100000"
	Ten6  = "1000000"
	Ten16 = "10000000000000000"
	Ten17 = "100000000000000000"
	Ten18 = "1000000000000000000"
)

var (
	BigMaxUint256, _ = helpers.BigIntFromString("115792089237316195423570985008687907853269984665640564039457584007913129639935")
)