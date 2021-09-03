package recipe

import (
	"bytes"
	"encoding/csv"
	"strings"
	"testing"
)

func BenchmarkExecute(b *testing.B) {
	// 5 lines of input
	input := "voter_id,first,last,address,city,state,zipcode,birthdate,party,sent\n68438357,Hazel,Dooley,290 Brekke Center,New Ford,Alaska,38953,1950-02-16,IND,\n75375390,Melyna,Yost,9768 Tina Terrace,Wunschview,North Carolina,10773,1954-03-23,,\n44195534,Uriah,Padberg,22332 Princess Point,West Sadieborough,Nebraska,77076,1992-09-17,REP,\n44371895,Helene,Kiehn,45605 Virgil Stravenue,Port Morrisfurt,Maryland,67655,1922-11-06,,\n47327331,Janet,Gaylord,33631 Winifred Estate,Port Wilmaville,Texas,41727,1970-07-06,REP,\n"
	reader := csv.NewReader(strings.NewReader(input))
	recipe := "!1 <- 1 # voter_id header\n1 <- 1 # voter_id\n!2 <- 2 # first header\n2 <- 2 # first\n!3 <- 3 # last header\n3 <- 3 # last\n!4 <- 4 # address header\n4 <- 4 # address\n!5 <- 5 # city header\n5 <- 5 # city\n!6 <- 6 # state header\n6 <- 6 # state\n!7 <- 7 # zipcode header\n7 <- 7 # zipcode\n!8 <- 8 # birthdate header\n8 <- 8 # birthdate\n!9 <- 9 # party header\n9 <- 9 # party\n!10 <- 10 # sent header\n10 <- 10 # sent\n$username <- firstchars(\"1\", 2) + 3 -> replace(\"'\", \"\") -> lowercase\n!11 <- \"username\"\n11 <- $username\n!12 <- \"email\"\n12 <- $username + \"@gmail.com\"\n"
	transformation, _ := Parse(strings.NewReader(recipe))
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)
	buf.Reset()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		transformation.Execute(reader, writer, true, -1)
	}
}
