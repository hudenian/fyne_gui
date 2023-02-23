package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	contract "fyne_gui/contract"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	"strconv"
)

var (
	rpc          = binding.NewString()
	contractAddr = binding.NewString()
)

func main() {
	myApp := app.New()
	w := myApp.NewWindow("Erc721A contract query")

	_ = rpc.Set("https://eth-goerli.public.blastapi.io")
	_ = contractAddr.Set("0x8ece72a85879e2c289b32f38ea35c307b8e50d1c")

	rpcEntry, contractAddressEntry := conditionFrom()

	from := widget.NewForm(
		widget.NewFormItem("rpc", rpcEntry),
		widget.NewFormItem("contract address", contractAddressEntry),
	)
	//生成查询ownerOf接口表单
	queryOwnerOf(from)

	c := container.NewVBox(from)
	w.SetContent(c)
	w.Resize(fyne.NewSize(720, 460))
	w.ShowAndRun()
}

/**
查询token的owner表单
*/
func queryOwnerOf(from *widget.Form) {
	tokenId := widget.NewEntry()
	tokenId.SetText("0")
	totalSupply := widget.NewEntry()
	owner := widget.NewEntry()
	baseUri := widget.NewEntry()
	ownerQBtn := widget.NewButton("query", func() {
		rpcUrl, _ := rpc.Get()
		client, err := ethclient.Dial(rpcUrl)
		if err != nil {
			owner.SetText("Failed to Dial ")
			return
		}
		defer client.Close()
		totalSupply.SetText("success to Dial")
		contractAddress, _ := contractAddr.Get()
		address := common.HexToAddress(contractAddress)
		instance, err := contract.NewUserERC721A(address, client)
		if err != nil {
			owner.SetText(err.Error())
		}
		supply, _ := instance.TotalSupply(&bind.CallOpts{Pending: true})
		totalSupply.SetText(supply.String())

		valInt, _ := strconv.Atoi(tokenId.Text)
		addr, _ := instance.OwnerOf(&bind.CallOpts{Pending: true}, big.NewInt(int64(valInt)))
		if addr.String() == "0x0000000000000000000000000000000000000000" {
			owner.SetText("token id not exist")
		} else {
			owner.SetText(addr.String())
		}

		uri, _ := instance.TokenURI(&bind.CallOpts{Pending: true}, big.NewInt(int64(valInt)))
		if len(uri) > 10 {
			baseUri.SetText(fmt.Sprintf("%s%s", "https://gateway.ipfs.io/ipfs/", uri[7:]))
		} else {
			baseUri.SetText("token id not exist")
		}
	})
	ownerQBtn.Importance = widget.HighImportance

	from.Append("totalSupply", totalSupply)
	from.Append("tokenId", tokenId)
	from.Append("owner", owner)
	from.Append("token uri", baseUri)
	from.Append("", ownerQBtn)
}

/**
查询条件窗口
*/
func conditionFrom() (*widget.Entry, *widget.Entry) {
	//节点rpc
	rpcEntry := widget.NewEntryWithData(rpc)
	//合约地址
	contractAddressEntry := widget.NewEntryWithData(contractAddr)
	return rpcEntry, contractAddressEntry
}
