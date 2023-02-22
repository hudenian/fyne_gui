package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	contract "fyne_gui/contract"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	"strconv"
)

func main() {
	myApp := app.New()
	w := myApp.NewWindow("UserErc721A contract query")

	//节点rpc
	rpc := widget.NewEntry()
	rpc.SetText("https://eth-goerli.public.blastapi.io")
	//合约地址
	contractAddress := widget.NewEntry()
	contractAddress.SetText("0x8ece72a85879e2c289b32f38ea35c307b8e50d1c")
	tokenId := widget.NewEntry()
	tokenId.SetText("0")
	owner := widget.NewLabel("")

	from := widget.NewForm(
		widget.NewFormItem("rpc:", rpc),
		widget.NewFormItem("contract address:", contractAddress),
		widget.NewFormItem("tokenId:", tokenId),
		widget.NewFormItem("owner:", owner),
	)
	from.CancelText = "clear"
	from.SubmitText = "query"

	//from.Items[0].HintText = "chain rpc url"
	//from.Items[1].HintText = "chain contract address"
	//from.Items[2].HintText = "contract token id"
	//from.Items[3].HintText = "token id owner"

	from.OnSubmit = func() {
		client, err := ethclient.Dial(rpc.Text)
		if err != nil {
			owner.SetText("Failed to Dial ")
			return
		}
		defer client.Close()
		owner.SetText("success to Dial")

		address := common.HexToAddress(contractAddress.Text)
		instance, err := contract.NewUserERC721A(address, client)
		if err != nil {
			owner.SetText(err.Error())
		}
		valInt, _ := strconv.Atoi(tokenId.Text)
		address, _ = instance.OwnerOf(&bind.CallOpts{Pending: true}, big.NewInt(int64(valInt)))
		owner.SetText(address.String())

	}
	from.OnCancel = func() {
		owner.SetText("")
	}

	c := container.NewVBox(from)
	w.SetContent(c)
	w.Resize(fyne.NewSize(640, 460))
	w.ShowAndRun()
}
