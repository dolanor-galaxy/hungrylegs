package importer

import (
	"bytes"
	"database/sql"
	"fmt"
	"io/ioutil"

	"github.com/therohans/HungryLegs/src/tcx"
	"github.com/tormoder/fit"
)

type Importer interface {
	Import(file string, db *sql.DB) error
}

////////////////////////////////////

type FitFile struct{}

func (f *FitFile) Import(file string, db *sql.DB) error {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	// Decode the FIT file data
	fit, err := fit.Decode(bytes.NewReader(data))
	if err != nil {
		return err
	}
	fmt.Println(fit.Type())
	return nil
}

////////////////////////////////////

type TcxFile struct{}

func (f *TcxFile) Import(file string, db *sql.DB) error {
	fmt.Println(file)
	tcxdb, err := tcx.ReadFile(file)
	if err != nil {
		return err
	}
	fmt.Printf("%v", tcxdb)
	return nil
}
