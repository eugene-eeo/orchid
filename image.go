package main

const DefaultImage string = "\u001B[38;5;147m" + `        _
    _ (` + " - " + `) _
  /` + "` '.\\ /.' `" + `\
  ` + "``" + `'-.,=,.-'` + "``" + `
    .'//v\\'.
   (_/\ " /\_)
       '-'`

type image interface {
	Render() string
}

type defaultImage struct{}

func (d *defaultImage) Render() string {
	return DefaultImage
}
