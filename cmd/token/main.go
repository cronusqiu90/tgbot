package main

import (
	"context"
	"fmt"

	"github.com/gotd/td/session"
	"github.com/gotd/td/session/tdesktop"
)

func main(){

	tdataDirPath := "tdata/"
	accounts, err := tdesktop.Read(tdataDirPath, nil)
	if err != nil {
		panic(err)
	}
	data, err := session.TDesktopSession(accounts[0])
	fmt.Println("DC:", data.DC, "IP:", data.Addr)

	fs := new(session.FileStorage)
	fs.Path = "session.json"
	loader := session.Loader{Storage: fs}
	if err := loader.Save(context.Background(), data); err != nil {
		panic(err)
	}
	fmt.Println("done.")


}
