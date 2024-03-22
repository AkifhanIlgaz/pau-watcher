package chain

type Chain struct {
	Name string
	Scan string
	Swap Swap
}

type Swap struct {
	Name string
	Url  string
}

// TODO: Add other chains
var (
	base = Chain{
		Name: "base",
		Swap: Swap{
			Name: "Uniswap",
			Url:  "https://app.uniswap.org/swap",
		},
		Scan: "https://basescan.org/tokentxns",
	}
	fantom = Chain{
		Name: "fantom",
		Swap: Swap{
			Name: "SpookySwap",
			Url:  "https://spooky.fi/#/swap",
		},
		Scan: "https://ftmscan.com/tokentxns",
	}
)

var Chains = map[string]Chain{
	"base":   base,
	"fantom": fantom,
}
