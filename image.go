package main

const DefaultImage string = "\u001B[38;5;109m" + `⠀⠀⠀⠀⠀⢀⣔⣂⢀⡀⠀⠀⠀⠀
⠀⠀⠐⣔⠍⠄⠉⠉⠙⠓⣯⡵⡄⠀
⢠⣧⡯⠁⠀⠀⠀⠀⠀⠀⠀⣿⣿⡄
⣸⡷⡧⠀⠀⠀⠀⠀⠀⠀⠀⣸⣿⡅
⢨⣿⡯⠄⠀⠀⠀⠀⠀⠀⢀⣼⣿⠃
⠀⠹⢟⣶⣤⣤⠀⠀⢀⣰⣾⣟⠟⠀
⠀⠀⠈⠙⠟⠻⠁⠀⢀⠢⠑⠁⠀⠀`

type image interface {
	Render() string
}

type defaultImage struct{}

func (d *defaultImage) Render() string {
	return DefaultImage
}
