package main

import (
	"fmt"
	"image"
	"image/png"
	"log"
	"os"

	"github.com/surma/mnist"
	"github.com/voxelbrain/goptions"
)

var (
	options = struct {
		ImageFile *os.File      `goptions:"-i, --image, description='Image file', obligatory, rdonly"`
		LabelFile *os.File      `goptions:"-l, --label, description='Label file', rdonly"`
		ImageList []int         `goptions:"-e, --extract, description='Index to extract (default: all)'"`
		Help      goptions.Help `goptions:"-h, --help, description='Show this help'"`
	}{}
)

type LabelReader interface {
	ReadLabel(i int) (mnist.Label, error)
}

type NullLabelReader string

func (NullLabelReader) ReadLabel(i int) (mnist.Label, error) {
	return mnist.Label(0), nil
}

type Item struct {
	Index int
	Label int
}

func main() {
	goptions.ParseAndFail(&options)
	defer options.ImageFile.Close()

	ir := &mnist.ImageReader{
		ReadSeeker: options.ImageFile,
	}
	if err := ir.ValidateHeader(); err != nil {
		log.Fatalf("Image file has invalid formalt: %s", err)
	}

	var lr LabelReader
	if options.LabelFile != nil {
		defer options.LabelFile.Close()
		mlr := &mnist.LabelReader{
			ReadSeeker: options.LabelFile,
		}
		if err := mlr.ValidateHeader(); err != nil {
			log.Fatalf("Label file has invalid formalt: %s", err)
		}
		if mlr.NumLabels() != ir.NumImages() {
			log.Fatalf("Image file and label file don't contain the same number of items")
		}
		lr = mlr
	} else {
		lr = NullLabelReader("")
	}

	c := make(chan int)
	go func() {
		if options.ImageList != nil {
			for _, i := range options.ImageList {
				if i < 0 || i > ir.NumImages() {
					log.Printf("Index out of range: %d (must be 1<x<%d)", i, ir.NumImages())
					continue
				}
				c <- i
			}
		} else {
			for i := 0; i < ir.NumImages(); i++ {
				c <- i
			}
		}
		close(c)
	}()

	items := labelize(c, lr)

	for item := range items {
		img, err := ir.ReadImage(item.Index)
		if err != nil {
			log.Printf("Could not read image #%d: %s. Skipping...", item.Label, err)
			continue
		}
		if err := saveImage(fmt.Sprintf("%06d_%d.png", item.Index, item.Label), img); err != nil {
			log.Printf("Could not save image #%d: %s", item.Index, err)
		}
	}
}

func labelize(cin <-chan int, lr LabelReader) <-chan Item {
	cout := make(chan Item)
	go func() {
		for i := range cin {
			l, err := lr.ReadLabel(i)
			if err != nil {
				log.Fatalf("Could not read label #%d: %s", i, err)
			}
			cout <- Item{
				Index: i,
				Label: int(l),
			}
		}
		close(cout)
	}()
	return cout
}

func saveImage(name string, img image.Image) error {
	f, err := os.Create(name)
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}
