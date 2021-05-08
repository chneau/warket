# warket

warframe.market client API

## dev

```bash
# use nodemon for hotreload
nodemon -e go -x "go run . s -g 0 || false"
```

## Doc

Please find it [here](https://docs.google.com/document/d/1121cjBNN4BeZdMBGil6Qbuqse-sWpEXPpitQH5fb_Fo/edit##heading=h.irwashnbboeo).

Other links:

- https://github.com/search?o=desc&q=api.warframe.market&s=updated&type=Repositories
- https://github.com/LastExceed/WarframeMarKT
- https://warframe.market/api_docs

```bash
# interesting path
/profile/orders # POST to create a new order
# payload: {"order_type":"sell","item_id":"54a74455e779892d5e5156cc","platinum":210,"quantity":1,"visible":false,"rank":1}
/profile/orders/{orderId} # PUT to update an already existing order
# payload: {order_id: "60965e5dea937404198b037c", platinum: 211, quantity: 1, visible: false, rank: 1}
```

## weird console output

Please run warket like this `winpty warket ...`

## Easy install

TODO: download with powershell like this `(new-object System.Net.WebClient).DownloadFile("https://github.com/chneau/warket/releases/download/pre-release/windows_amd64_warket.exe","C:\tmp\file.txt")`

Install Chocolatey (sort of a package manager for Windows)  
`https://chocolatey.org/install`

Then, install Go.  
`choco install golang`

Then, install this repo.  
`go install -v github.com/chneau/warket@latest`

If on Windows, be sure to have this in your PATH system environment variable:  
`%userprofile%\go\bin`
