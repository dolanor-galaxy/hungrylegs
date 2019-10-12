package importer

type Validator interface {
	Validate([]byte) error
}

type Importer interface {
	Import([]byte) error
}

////////////////////////////////////

type FitFile struct{}

func (f *FitFile) Validate([]byte) error {
	return nil
}

func (f *FitFile) Import([]byte) error {
	return nil
}

////////////////////////////////////

type TcxFile struct{}

func (f *TcxFile) Validate([]byte) error {
	return nil
}

func (f *TcxFile) Import([]byte) error {
	return nil
}
