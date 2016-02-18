package main

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/reiki4040/cstore"
)

func main() {
	pkgs, err := DoDpkg()
	if err != nil {
		fmt.Println(err)
		return
	}

	d := time.Now().Format(time.RFC3339)

	pl := PkgList{
		Pkgs: pkgs,
		Date: d,
	}

	cs, err := cstore.NewCStore("dpkg", "./dpkg.json", cstore.JSON)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = cs.SaveWithoutValidate(pl)
	if err != nil {
		fmt.Println(err)
		return
	}
}

type PkgList struct {
	Pkgs []Pkg  `json:"packages"`
	Date string `json:"date"`
}

type Pkg struct {
	Name        string `json:"name"`
	Arch        string `json:"arch"`
	Version     string `json:"version"`
	Description string `json:"description"`
}

func DoDpkg() ([]Pkg, error) {
	out, err := exec.Command("dpkg", "-l").Output()
	if err != nil {
		return nil, err
	}

	iName := -1
	iVer := -1
	iArch := -1
	iDesc := -1
	lines := strings.Split(string(out), "\n")
	pkgs := make([]Pkg, 0, len(lines))
	for _, l := range lines {
		if iName == -1 || iVer == -1 || iArch == -1 || iDesc == -1 {
			iName = strings.Index(l, "Name")
			iVer = strings.Index(l, "Version")
			iArch = strings.Index(l, "Architecture")
			iDesc = strings.Index(l, "Description")

			continue
		}

		if i := strings.Index(l, "ii "); i == 0 {
			p := Pkg{
				Name:        strings.TrimSpace(l[iName:iVer]),
				Version:     strings.TrimSpace(l[iVer:iArch]),
				Arch:        strings.TrimSpace(l[iArch:iDesc]),
				Description: strings.TrimSpace(l[iDesc:]),
			}

			pkgs = append(pkgs, p)
		}
	}

	return pkgs, nil
}
