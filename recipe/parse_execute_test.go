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
		{
			name:          "column recipe, more complex",
			recipe:        "1 <- 3 + 2\n2 <- 1 + 3\n3 <- 2 + 1\n",
			input:         "a,b,c\nd,e,f\ng,h,i",
			processHeader: false,
			want:          "cb,ac,ba\nfe,df,ed\nih,gi,hg\n",
		},
		{
			name:          "column recipe, same as above, but variables first",
			recipe:        "$a <- 3+2\n$b<-1+3\n$c<-2+1\n1<-$a\n2<-$b\n3<-$c\n",
			input:         "a,b,c\nd,e,f\ng,h,i",
			processHeader: false,
			want:          "cb,ac,ba\nfe,df,ed\nih,gi,hg\n",
		},
		{
			name:          "upper 1, lower 2 - function test #1",
			recipe:        "!1 <- \"FRUIT\"\n1 <- 1 -> uppercase\n!2 <- \"veggies\"\n2 <- 2 -> lowercase",
			input:         "thing1,thing2\napple,artichoke\nBANANA,BEET\nCucumber,Carrot\n",
			processHeader: true,
			want:          "FRUIT,veggies\nAPPLE,artichoke\nBANANA,beet\nCUCUMBER,carrot\n",
		},
		{
			name:          "same as above but not using placeholder",
			recipe:        "!1 <- \"FRUIT\"\n1 <- uppercase(1)\n!2 <- \"veggies\"\n2 <- lowercase(2)",
			input:         "thing1,thing2\napple,artichoke\nBANANA,BEET\nCucumber,Carrot\n",
			processHeader: true,
			want:          "FRUIT,veggies\nAPPLE,artichoke\nBANANA,beet\nCUCUMBER,carrot\n",
		},
		{
			name:          "using join as a pipe function",
			recipe:        "1 <- 1 -> join -> 1",
			input:         "a\nb\n",
			processHeader: false,
			want:          "aa\nbb\n",
		},
		{
			name:          "using join as a function",
			recipe:        "1 <- 1 -> join(1)",
			input:         "a\nb\n",
			processHeader: false,
			want:          "aa\nbb\n",
		},
		{
			name:          "using join as a function and joining to it",
			recipe:        "1 <- 1 + join(1)",
			input:         "a\nb\n",
			processHeader: false,
			want:          "aa\nbb\n",
		},
		{
			name:          "use add to sum two integer columns",
			recipe:        "!1 <- \"fruits\"\n!2 <- \"veggies\"\n!3 <- \"total\"\n1 <- 1\n2 <- 2\n3 <- add(1,2)",
			input:         "a,b\n1,2\n555,444\n13,31\n",
			processHeader: true,
			want:          "fruits,veggies,total\n1,2,3\n555,444,999\n13,31,44\n",
		},
		{
			name:          "use addFloat to sum two float/int columns",
			recipe:        "!1 <- \"fruits\"\n!2 <- \"veggies\"\n!3 <- \"total\"\n1 <- 1\n2 <- 2\n3 <- addFloat(1,2,\"2\")\n",
			input:         "a,b\n1,2\n555.55,444.44\n13.55,31.44\n",
			processHeader: true,
			want:          "fruits,veggies,total\n1,2,3.00\n555.55,444.44,999.99\n13.55,31.44,44.99\n",
		},
		{
			name:          "use addFloat to sum two float/int into rounded ints",
			recipe:        "!1 <- \"fruits\"\n!2 <- \"veggies\"\n!3 <- \"total\"\n1 <- 1\n2 <- 2\n3 <- addFloat(1,2,\"0\")\n",
			input:         "a,b\n1,2\n555.55,444.44\n13.55,31.44\n",
			processHeader: true,
			want:          "fruits,veggies,total\n1,2,3\n555.55,444.44,1000\n13.55,31.44,45\n",
		},
		{
			name:          "use addFloat to sum two float/int with no rounding",
			recipe:        "!1 <- \"fruits\"\n!2 <- \"veggies\"\n!3 <- \"total\"\n1 <- 1\n2 <- 2\n3 <- addFloat(1,2,\"-1\")\n",
			input:         "a,b\n1,2\n555.55,444.44\n13.55,31.44\n",
			processHeader: true,
			want:          "fruits,veggies,total\n1,2,3.000000\n555.55,444.44,999.990000\n13.55,31.44,44.990000\n",
		},
		{
			name:          "add with non-int arg1 is an error",
			recipe:        "1 <- add(1, 2)\n",
			input:         "a,2\n",
			processHeader: false,
			wantErr:       true,
			wantErrText:   "first arg to Add was not an integer: a",
		},
		{
			name:          "add with non-int arg2 is an error",
			recipe:        "1 <- add(2,1)\n",
			input:         "a,2\n",
			processHeader: false,
			wantErr:       true,
			wantErrText:   "second arg to Add was not an integer: a",
		},
		{
			name:          "addFloat with non-int arg1 is an error",
			recipe:        "1 <- addfloat(1, 2)\n",
			input:         "a,2\n",
			processHeader: false,
			wantErr:       true,
			wantErrText:   "first arg to AddFloat was not numeric: a",
		},
		{
			name:          "addFloat with non-int arg2 is an error",
			recipe:        "1 <- addfloat(2, 1)\n",
			input:         "a,2\n",
			processHeader: false,
			wantErr:       true,
			wantErrText:   "second arg to AddFloat was not numeric: a",
		},
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
			if tt.wantErr && (err != nil) && err.Error() != tt.wantErrText {
				t.Errorf("got execute error text = %v, want error text = %v", err.Error(), tt.wantErrText)
			}

			got := b.String()
			if got != tt.want {
				t.Errorf("Dump() = %v, want %v", got, tt.want)
			}
		})
	}
}
