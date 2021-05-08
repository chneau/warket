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

## weird console output

Please run warket like this `winpty warket ...`

## Easy install

TODO: download with powershell like this `(new-object System.Net.WebClient).DownloadFile("https://github.com/chneau/warket/releases/download/pre-release/windows_amd64_warket.exe","C:\tmp\file.txt")`

Install Chocolatey (sort of a package manager for Windows)  
`https://chocolatey.org/install`

Then, install Go.  
`choco install golang`

Then, install this repo.  
`go install -v github.com/chneau/warket`

If on Windows, be sure to have this in your PATH system environment variable:  
`%userprofile%\go\bin`
