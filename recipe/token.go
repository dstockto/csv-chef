package recipe

//go:generate stringer -type=Token

type Token int

const (
	ILLEGAL     Token = iota
	EOF               //1 - end of file
	WS                //2 - space, tab, newline
	NEWLINE           //3 - \n (probably not needed)
	COLUMN_ID         //4 - digits
	ASSIGNMENT        //5 - <-
	PIPE              //6 - ->
	COMMENT           //7 - # ...
	PLACEHOLDER       //8 - ?
	PLUS              //9 - +
	LITERAL           //10 - "quoted"
	VARIABLE          //11 - starts w/ $
	FUNCTION          //12 - letters
	OPEN_PAREN        //13 - (
	CLOSE_PAREN       //14 - )
	COMMA             //15 - ,
	HEADER            //16 - !<digits>
)
