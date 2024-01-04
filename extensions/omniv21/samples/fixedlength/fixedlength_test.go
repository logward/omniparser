package fixedlength

import (
	"bytes"
	"io"
	"io/ioutil"
	"testing"

	"github.com/bradleyjkemp/cupaloy"
	"github.com/jf-tech/go-corelib/jsons"
	"github.com/logward/omniparser"

	"github.com/logward/omniparser/extensions/omniv21/samples"
	"github.com/logward/omniparser/transformctx"
)

type testCase struct {
	schemaFile string
	inputFile  string
	schema     omniparser.Schema
	input      []byte
}

const (
	test1_Single_Row = iota
	test2_Multi_Rows
	test3_Header_Footer
)

var tests = []testCase{
	{
		// test1_Single_Row
		schemaFile: "./1_single_row.schema.json",
		inputFile:  "./1_single_row.input.txt",
	},
	{
		// test2_Multi_Rows
		schemaFile: "./2_multi_rows.schema.json",
		inputFile:  "./2_multi_rows.input.txt",
	},
	{
		// test3_Header_Footer
		schemaFile: "./3_header_footer.schema.json",
		inputFile:  "./3_header_footer.input.txt",
	},
}

func init() {
	for i := range tests {
		schema, err := ioutil.ReadFile(tests[i].schemaFile)
		if err != nil {
			panic(err)
		}
		tests[i].schema, err = omniparser.NewSchema("bench", bytes.NewReader(schema))
		if err != nil {
			panic(err)
		}
		tests[i].input, err = ioutil.ReadFile(tests[i].inputFile)
		if err != nil {
			panic(err)
		}
	}
}

func (tst testCase) doTest(t *testing.T) {
	cupaloy.SnapshotT(t, jsons.BPJ(samples.SampleTestCommon(t, tst.schemaFile, tst.inputFile)))
}

func (tst testCase) doBenchmark(b *testing.B) {
	for i := 0; i < b.N; i++ {
		transform, err := tst.schema.NewTransform(
			"bench", bytes.NewReader(tst.input), &transformctx.Ctx{})
		if err != nil {
			b.FailNow()
		}
		for {
			_, err = transform.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				b.FailNow()
			}
		}
	}
}

func Test1_Single_Row(t *testing.T) {
	tests[test1_Single_Row].doTest(t)
}

func Test2_Multi_Rows(t *testing.T) {
	tests[test2_Multi_Rows].doTest(t)
}

func Test3_Header_Footer(t *testing.T) {
	tests[test3_Header_Footer].doTest(t)
}

// Benchmark1_Single_Row-8      	   25869	     45576 ns/op	   27721 B/op	     644 allocs/op
func Benchmark1_Single_Row(b *testing.B) {
	tests[test1_Single_Row].doBenchmark(b)
}

// Benchmark2_Multi_Rows-8      	   18813	     63901 ns/op	   29167 B/op	     635 allocs/op
func Benchmark2_Multi_Rows(b *testing.B) {
	tests[test2_Multi_Rows].doBenchmark(b)
}

// Benchmark3_Header_Footer-8   	    5857	    197326 ns/op	   82234 B/op	    2009 allocs/op
func Benchmark3_Header_Footer(b *testing.B) {
	tests[test3_Header_Footer].doBenchmark(b)
}
