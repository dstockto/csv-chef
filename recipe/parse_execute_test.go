package recipe

import (
	"bytes"
	"encoding/csv"
	"strings"
	"testing"
)

func TestTransformation_ParseExecute(t *testing.T) {
	tests := []struct {
		name             string
		recipe           string
		input            string
		processHeader    bool
		want             string
		wantParseErr     bool
		wantParseErrText string
		wantErr          bool
		wantErrText      string
	}{
		{
			name:          "simple 1 <- 1",
			recipe:        "!1 <- 1\n1 <- 1\n",
			processHeader: true,
			input:         "a,b\n",
			want:          "a\n",
		},
		{
			name:          "empty recipe is a parse error",
			processHeader: false,
			want:          "",
			wantErr:       true,
			wantErrText:   "no column recipes provided",
		},
		{
			name:          "referencing a header for a column that has no reference is an error",
			recipe:        "1 <- \"hi\"\n!3 <- \"lala\"",
			processHeader: true,
			wantErr:       true,
			wantErrText:   "found header for column 3, but no recipe for column 3",
		},
		{
			name:          "process headers with no header recipe",
			recipe:        "1<-2\n2<-1\n",
			input:         "a,b\n",
			processHeader: true,
			want:          "a,b\n",
		},
		{
			name:          "header recipe with literals",
			recipe:        "1<-1\n2<-2\n!2<-\"apple\"\n",
			input:         "a,b\n",
			processHeader: true,
			want:          "a,apple\n",
		},
		{
			name:          "header recipe with joining literals",
			recipe:        "!1<- \"alpha\"+\" beta\"\n1<-1\n2<-2\n",
			input:         "a,b\n",
			processHeader: true,
			want:          "alpha beta,b\n",
		},
		{
			name:          "double join flip flop headers",
			recipe:        "!1<-2+1\n!2<-1+2\n1<-1\n2<-2\n",
			input:         "alpha,beta\n",
			processHeader: true,
			want:          "betaalpha,alphabeta\n",
		},
		{
			name:          "header referencing variable that does not exist is an error",
			recipe:        "!1<-$bar\n1<-1\n",
			input:         "a,b\n",
			processHeader: true,
			wantErr:       true,
			wantErrText:   "variable '$bar' referenced, but it is not defined",
		},
		{
			name:          "headers via variables",
			recipe:        "$foo<-2\n1<-$foo\n!1<-$foo\n",
			input:         "apple,banana\n",
			processHeader: true,
			want:          "banana\n",
		},
		{
			name:          "referencing header column that does not exist is an error",
			recipe:        "1 <- 1\n!1 <- 3\n",
			input:         "a,b\n",
			processHeader: true,
			wantErr:       true,
			wantErrText:   "column 3 referenced but it does not exist in input file",
		},
		{
			name:          "referencing variable that is not defined is an error",
			recipe:        "1<-1\n!1<-$foo\n",
			input:         "a,b",
			processHeader: true,
			wantErr:       true,
			wantErrText:   "variable '$foo' referenced, but it is not defined",
		},
		{
			name:          "double header using placeholder concatenation",
			recipe:        "!1 <- 1 + ?\n1<-1\n",
			input:         "ab,c\n",
			processHeader: true,
			want:          "abab\n",
		},
		{
			name:          "quad header using placeholder concatenation",
			recipe:        "!1 <- 1 + ? + ?\n1<-1\n",
			input:         "ab,c\n",
			processHeader: true,
			want:          "abababab\n",
		},
		{
			name:          "headers and column recipe, swap columns",
			recipe:        "!1 <- \"col1\"\n!2<-\"col2\"\n1<-2\n2<-1",
			input:         "first,last\na,b\nc,d\ne,f",
			processHeader: true,
			want:          "col1,col2\nb,a\nd,c\nf,e\n",
		},
		// TODO need to start executing functions
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transformation, err := Parse(strings.NewReader(tt.recipe))

			if (err != nil) != tt.wantParseErr {
				t.Errorf("parse error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantParseErr && err.Error() != tt.wantParseErrText {
				t.Errorf("got parse error text = %v, want error text = %v", err.Error(), tt.wantErrText)
			}

			var b bytes.Buffer
			writer := csv.NewWriter(&b)

			err = transformation.Execute(csv.NewReader(strings.NewReader(tt.input)), writer, tt.processHeader)
			if (err != nil) != tt.wantErr {
				t.Errorf("execute error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && err.Error() != tt.wantErrText {
				t.Errorf("got execute error text = %v, want error text = %v", err.Error(), tt.wantErrText)
			}

			got := b.String()
			if got != tt.want {
				t.Errorf("Dump() = %v, want %v", got, tt.want)
			}
		})
	}
}
