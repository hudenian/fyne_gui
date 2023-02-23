package main

import (
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
	w := myApp.NewWindow("UserErc721A合约查询")

	rpc.Set("https://eth-goerli.public.blastapi.io")
	contractAddr.Set("0x8ece72a85879e2c289b32f38ea35c307b8e50d1c")

	rpcEntry, contractAddressEntry := conditionFrom()

	from := widget.NewForm(
		widget.NewFormItem("rpc:", rpcEntry),
		widget.NewFormItem("contract address:", contractAddressEntry),
	)

	tokenIdEntry, ownerLabel, queryBtn := queryOwnerOf()
	from.Append("tokenId", tokenIdEntry)
	from.Append("owner", ownerLabel)
	from.Append("query owner", queryBtn)

	c := container.NewVBox(from)
	w.SetContent(c)
	w.Resize(fyne.NewSize(640, 460))
	w.ShowAndRun()
}

/**
查询token的owner表单
*/
func queryOwnerOf() (*widget.Entry, *widget.Label, *widget.Button) {
	tokenId := widget.NewEntry()
	tokenId.SetText("0")
	owner := widget.NewLabel("")
	ownerQBtn := widget.NewButton("query owner", func() {
		rpcUrl, _ := rpc.Get()
		client, err := ethclient.Dial(rpcUrl)
		if err != nil {
			owner.SetText("Failed to Dial ")
			return
		}
		defer client.Close()
		owner.SetText("success to Dial")
		contractAddress, _ := contractAddr.Get()
		address := common.HexToAddress(contractAddress)
		instance, err := contract.NewUserERC721A(address, client)
		if err != nil {
			owner.SetText(err.Error())
		}
		valInt, _ := strconv.Atoi(tokenId.Text)
		address, _ = instance.OwnerOf(&bind.CallOpts{Pending: true}, big.NewInt(int64(valInt)))
		owner.SetText(address.String())
	})
	ownerQBtn.Importance = widget.HighImportance

	return tokenId, owner, ownerQBtn
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
