package pdiutil

import (
	"testing"
)

var sampleXML = `
	<content>
		<p>this is content area</p>
		<animal>
			<p>This id dog</p>
			<dog>
			   <p>tommy</p>
			</dog>
		</animal>
		<birds>
			<p>this is birds</p>
			<p>this is birds</p>
		</birds>
		<animal>
			<p>this is animals</p>
		</animal>
	</content>`

func TestWalkXMLNode(t *testing.T) {

	type args struct {
		nodes []*XMLNode
		f     func(*XMLNode) bool
	}
	tests := []struct {
		name string
		args args
	}{
		{
			"Basic Parse Test",
			args{
				[]*XMLNode{ParseXML([]byte(sampleXML))},
				func(n *XMLNode) bool {
					return true
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			WalkXMLNode(tt.args.nodes, tt.args.f)
		})
	}
}

func TestCountXMLComplexity(t *testing.T) {
	type args struct {
		xmlBytes []byte
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			"Simple XML Count",
			args{[]byte(sampleXML)},
			11,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CountXMLComplexity(tt.args.xmlBytes); got != tt.want {
				t.Errorf("CountXMLComplexity() = %v, want %v", got, tt.want)
			}
		})
	}
}
