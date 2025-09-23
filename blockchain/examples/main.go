package main

import (
	"encoding/json"
	"fmt"

	"github.com/Conflux-Chain/go-conflux-util/cmd"
	"github.com/ethereum/go-ethereum/common"
	"github.com/openweb3/web3go"
	"github.com/v3-Swampy/points-service/blockchain"
)

var (
	swappiAddresses = blockchain.SwappiAddresses{
		Factory: common.HexToAddress("0xE2a6F7c0ce4d5d300F97aA7E125455f5cd3342F5"),
		USDT:    common.HexToAddress("0xfe97e85d13abd9c1c33384e796f10b73905637ce"),
		WCFX:    common.HexToAddress("0x14b2d3bc65e74dae1030eafd8ac30c533c976a9b"),
	}

	// ABC/WCFX/USDT
	abc   = common.HexToAddress("0x905f2202003453006eaf975699545f2e909079b8")
	abcLP = common.HexToAddress("0x700d841e087f4038639b214e849beab622f178c6")
)

func main() {
	client := web3go.MustNewClient("http://evm.confluxrpc.com")
	defer client.Close()

	caller, _ := client.ToClientForContract()
	erc20 := blockchain.NewERC20(caller)
	swappi := blockchain.NewSwappi(caller, erc20, swappiAddresses)

	// get token info: USDT
	token, err := erc20.GetTokenInfo(swappiAddresses.USDT)
	cmd.FatalIfErr(err, "Failed to get token info")
	fmt.Println("Token info:", mustToJson(token))

	// get LP token info: ABC/WCFX
	pool, err := swappi.GetPairInfo(abcLP)
	cmd.FatalIfErr(err, "Failed to get pool info")
	fmt.Println("Pool info:", mustToJson(pool))

	// get ABC price
	tokenPrice, err := swappi.GetTokenPriceAuto(nil, abc)
	cmd.FatalIfErr(err, "Failed to get token price in Swappi")
	fmt.Println("Token price:", tokenPrice)

	// ABC-WCFX pool TVL
	tvl, err := swappi.GetPairTVL(nil, abcLP)
	cmd.FatalIfErr(err, "Failed to get TVL in Swappi")
	fmt.Println("Pool TVL:", tvl)
}

func mustToJson(v any) string {
	data, err := json.MarshalIndent(v, "", "    ")
	cmd.FatalIfErr(err, "Failed to marshal data")
	return string(data)
}
