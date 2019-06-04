package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"bufio"

	"github.com/quickfixgo/quickfix"
	"github.com/quickfixgo/quickfix/datadictionary"
)

func main() {
	dd, err := datadictionary.Parse("FIX41.xml")

	if err != nil {
		log.Fatal(err)
	}
	f, err := os.Open("fix.msg")
	if err != nil {
		log.Fatal(err)
	}
	fb := bufio.NewReader(f)
	for line, err := fb.ReadBytes(byte('\n')); err == nil; line, err = fb.ReadBytes(byte('\n')) {
		b := bytes.NewBuffer(line)
		m := quickfix.NewMessage()
		err = nil
		err = quickfix.ParseMessageWithDataDictionary(m, b, dd, dd)
		if err != nil {
			log.Println(err)
			return
		}
		tw := tabwriter.NewWriter(os.Stdout, 1, 1, 1, '.', 0)
		for _, t := range m.Body.Tags() {
			val, err := m.Body.GetBytes(t)
			if err != nil {
				log.Println(err)
				continue
			}

			d, ok := dd.FieldTypeByTag[int(t)]
			if !ok {
				log.Println(t, "not found in dictioary")
			}
			strval := string(val) + "\t"
			if len(d.Enums) > 0 {
				if e, ok := d.Enums[string(val)]; ok {
					strval = e.Value + "\t(" + e.Description + ")"
				}
			}
			fmt.Fprintf(tw, "%s\t(%d)\t=\t%s\n", d.Name(), t, strval)
		}
		tw.Flush()
		fmt.Println("----------------------------")
	}
}
