**dactorio-dat** is a utility for converting factorio mod-settings.dat to the old json format.

## Installing

### From sources
Go 1.2+ is required.

`go get -u github.com/betrok/factorio-dat`

### Prebuild
Static binaries: [linux](https://ttyh.ru/files/factorio-dat/lin/factorio-dat), [windows](https://ttyh.ru/files/factorio-dat/win/factorio-dat.exe), [mac](https://ttyh.ru/files/factorio-dat/mac/factorio-dat).

## Usage
`factorio-dat [in.dat] [out.json]`

Default input/output files are `mod-settings.dat` and `mod-settings.json` in the current directory.
`-` can be used for any of arguments to use stdin/stdout.

Factorio converts `mod-settings.json` to **.dat** at startup if the **.dat** file does not exist.

## License
This project is licensed under the terms of the MIT license.
