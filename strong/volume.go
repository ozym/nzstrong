package strong

import (
	"log"

	"os"
	"os/exec"
	"path/filepath"
)

type Volume struct {
	Processed string
	Source    string

	Vol2 string
	Vol3 string

	PgPlot string
	Ps2pdf string
	Plot1  string
	Plot2  string
	Plot3  string
	Plot4  string
}

func (v Volume) Build(trigger Trigger) error {
	if err := v.V1A(trigger); err != nil {
		return err
	}

	if err := v.V2A(trigger); err != nil {
		return err
	}

	if err := v.V3A(trigger); err != nil {
		return err
	}

	return nil
}

func (v Volume) V1A(trigger Trigger) error {
	if err := trigger.Build(v.Processed); err != nil {
		return err
	}

	input := trigger.V1A(v.Processed)
	if err := copyFiles(input, input+"_"); err != nil {
		return err
	}

	ps := PS{
		Dir:  v.PgPlot,
		Exec: v.Plot1,
	}

	psfile := trigger.PSPath(v.Processed, "Vol1")
	log.Println(psfile)
	defer os.Remove(psfile)

	// may have trouble if the ps file already exists
	if _, err := os.Stat(psfile); err == nil {
		if err := os.Remove(psfile); err != nil {
			return err
		}
	}

	// build ps file
	if err := ps.Plot(input + "_"); err != nil {
		return err
	}

	pdf := PDF{
		Exec: v.Ps2pdf,
	}

	pdffile := trigger.PDFPath(v.Processed, "Vol1")
	log.Println(pdffile)
	// convert postscript to pdf
	if err := pdf.Convert(psfile, pdffile); err != nil {
		return err
	}

	return nil
}

func (v Volume) V2A(trigger Trigger) error {

	input := trigger.V1A(v.Processed)
	defer os.Remove(input + "_")

	output := trigger.V2A(v.Processed)
	intermediate := filepath.Join(filepath.Dir(input), filepath.Base(output)+"_")
	defer os.Remove(intermediate)

	if err := copyFiles(input, input+"_"); err != nil {
		return err
	}

	log.Println(v.Vol2, filepath.Base(input+"_"))
	log.Printf("(cd %s;  %s %s)", filepath.Dir(input), v.Vol2, filepath.Base(input+"_"))
	cmd := exec.Command(v.Vol2, filepath.Base(input+"_"))
	cmd.Dir = filepath.Dir(input)
	if v.Source != "" {
		cmd.Env = []string{"SRC=" + v.Source}
	}
	if _, err := cmd.Output(); err != nil {
		return err
	}

	if err := copyFiles(intermediate, output); err != nil {
		return err
	}

	ps := PS{
		Dir:  v.PgPlot,
		Exec: v.Plot2,
	}

	psfile := trigger.PSPath(v.Processed, "Vol2")
	log.Println(psfile)
	defer os.Remove(psfile)

	// may have trouble if the ps file already exists
	if _, err := os.Stat(psfile); err == nil {
		if err := os.Remove(psfile); err != nil {
			return err
		}
	}

	defer os.Remove(output + "_")
	if err := copyFiles(output, output+"_"); err != nil {
		return err
	}

	// build ps file
	if err := ps.Plot(output + "_"); err != nil {
		return err
	}

	pdf := PDF{
		Exec: v.Ps2pdf,
	}

	pdffile := trigger.PDFPath(v.Processed, "Vol2")
	log.Println(pdffile)
	// convert postscript to pdf
	if err := pdf.Convert(psfile, pdffile); err != nil {
		return err
	}

	ps = PS{
		Dir:  v.PgPlot,
		Exec: v.Plot4,
	}

	psfile = trigger.PSPath(v.Processed, "Vol4")
	psfile = filepath.Join(filepath.Dir(output), filepath.Base(psfile))

	log.Println(psfile)
	//defer os.Remove(psfile)

	// may have trouble if the ps file already exists
	if _, err := os.Stat(psfile); err == nil {
		if err := os.Remove(psfile); err != nil {
			return err
		}
	}

	// build ps file
	if err := ps.Plot(output + "_"); err != nil {
		return err
	}

	pdffile = trigger.PDFPath(v.Processed, "Vol4")
	log.Println(pdffile)
	// convert postscript to pdf
	if err := pdf.Convert(psfile, pdffile); err != nil {
		return err
	}

	return nil
}

func (v Volume) V3A(trigger Trigger) error {

	input := trigger.V2A(v.Processed)
	defer os.Remove(input + "_")

	output := trigger.V3A(v.Processed)
	intermediate := filepath.Join(filepath.Dir(input), filepath.Base(output)+"_")
	defer os.Remove(intermediate)

	if err := copyFiles(input, input+"_"); err != nil {
		return err
	}

	log.Printf("(cd %s;  %s %s)", filepath.Dir(input), v.Vol3, filepath.Base(input+"_"))
	cmd := exec.Command(v.Vol3, filepath.Base(input+"_"))
	cmd.Dir = filepath.Dir(input)
	if v.Source != "" {
		cmd.Env = []string{"SRC=" + v.Source}
	}

	if _, err := cmd.Output(); err != nil {
		return err
	}

	if err := copyFiles(intermediate, output); err != nil {
		return err
	}

	ps := PS{
		Dir:  v.PgPlot,
		Exec: v.Plot3,
	}

	psfile := trigger.PSPath(v.Processed, "Vol3")
	log.Println(psfile)
	defer os.Remove(psfile)

	// may have trouble if the ps file already exists
	if _, err := os.Stat(psfile); err == nil {
		if err := os.Remove(psfile); err != nil {
			return err
		}
	}

	//defer os.Remove(output + "_")
	if err := copyFiles(output, output+"_"); err != nil {
		return err
	}

	// build ps file
	if err := ps.Plot(output + "_"); err != nil {
		return err
	}

	pdf := PDF{
		Exec: v.Ps2pdf,
	}

	pdffile := trigger.PDFPath(v.Processed, "Vol3")
	log.Println(pdffile)
	// convert postscript to pdf
	if err := pdf.Convert(psfile, pdffile); err != nil {
		return err
	}

	return nil
}

/*
func (v Volume) Plot(volume Volume, file File) error {

	// output files
	ps, pdf := volume.PSPath(v, file), v.PDFPath(file)
	defer os.Remove(ps)

	// will have trouble if the ps file already exists
	if _, err := os.Stat(ps); err == nil {
		if err := os.Remove(ps); err != nil {
			return err
		}
	}

	// make a copy of the v?a file ...
	if err := copyFiles(volume.Path(file), volume.Path(file)+"_"); err != nil {
		return err
	}
	defer os.Remove(volume.Path(file) + "_")

	// build ps file
	if err := v.PS.Plot(volume.Path(file) + "_"); err != nil {
		return err
	}

	// convert postscript to pdf
	if err := v.PDF.Convert(ps, pdf); err != nil {
		return err
	}

	return nil
}
*/
