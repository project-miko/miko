package tools

import "github.com/shopspring/decimal"

const (
	EthUnitWeiStr = "1000000000000000000"
)

func WeiToEther(d decimal.Decimal) decimal.Decimal {
	wei, _ := decimal.NewFromString(EthUnitWeiStr)
	return d.Div(wei)
}

func EtherToWei(d decimal.Decimal) decimal.Decimal {
	wei, _ := decimal.NewFromString(EthUnitWeiStr)
	return d.Mul(wei)
}
