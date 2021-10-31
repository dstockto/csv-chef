CSV-Chef
===

This project intends to make CSV->CSV transformations easy to define and execute. The program uses a "recipe" which
is a program that is designed and intended to be simple and easy to understand, even if you're not a developer.

There are two main modes of operation: bake and generate. Generate will be coming in a future version.

Bake
--
This is the main mode.

`csv-chef bake -i /path/to/input.csv -o /path/to/output.csv -r /path/to/recipefile`

The program requires an input CSV file, an output CSV file, and a recipe file. CSVs can be whatever you want as long as they are legitimate, correctly formatted CSVs. You can provide `-n` or `--lines` to specify how many lines you which to process. By default, it considers the first line a header and will follow the header recipe rules if provided. If not provided, the headers will remain unchanged from the input file, according to the column they originally were in with any extra columns written as `col #` where # is the column number that didn't have a header specified. To disable header processing please specify `--no-header` or `-d`.

Please see the recipes section for information about how to build recipes for the program.

Identity
==

The `identity` command provides a starter recipe for you based on the provided input file. You may specify that you want to include header recipe lines as well with the `-w` or `--with-headers` option. The command will output recipe that would give back the input file unchanged. If you specify `-o` or `--output` it can output to a file. If not, it would output to the console (stdout). 

If the output file already exists, `csv-chef identity` will stop and not write. If you want to write the file anyway, please provide the `-f` or `--force` option.

Example:

```
$ csv-chef identity -w input.csv

!1 <- 1 # voter_id header
1 <- 1 # voter_id
!2 <- 2 # first header
2 <- 2 # first
!3 <- 3 # last header
3 <- 3 # last
!4 <- 4 # address header
4 <- 4 # address
!5 <- 5 # city header
5 <- 5 # city
!6 <- 6 # state header
6 <- 6 # state
!7 <- 7 # zipcode header
7 <- 7 # zipcode
!8 <- 8 # birthdate header
8 <- 8 # birthdate
!9 <- 9 # party header
9 <- 9 # party
!10 <- 10 # sent header
10 <- 10 # sent
```

Recipes
==

There are three types of value can have data assigned: headers, columns and variables. Headers process only
the first line of output and only if they are turned on with the `--withHeaders` option.

Variables allow you to store information that can be reused within a column, and can be used to allow for more complex
operations than would otherwise be allowed with what you can do with a single column's recipe.

Within a recipe file, it's easy to identify columns, headers, variables, and functions because they are very limited in
how they can be defined.

Columns consist of only digits. If you see a number by itself, it's a column reference.

Headers are an exclamation point followed by a column number with no spaces, like `!2`. If you want to add a column
header for an inserted column, these can be useful. You could also use them to change existing headers. You can use all
the features of a recipe when defining a header, but remember, for transformations, it will run against existing
header values. For generate recipes, there are no incoming columns, so it doesn't make any sense to try to use column
references.

Variables can be identified because they start with a `$` and consist of letters, for example `$firstname`.

Functions consist of only letters. They can either be just letters, or they can potentially require arguments which
should be provided inside parentheses. If there are more than one, they should be separated by commas. Arguments to a
function can be columns, variables or literals.

Literals are plain text. They are wrapped in double quotes. For example, "Header 3". You can, if needed, include quotes
inside a literal by escaping them with a backslash, like so `"this \" <- is a quote"`. If you want to include a
backslash, you must also escape it as well, which means type two backslashes - `"backslash \\"`.

If you want to put comments in your recipe, you can! Use `#` to indicate a comment. This can be either at the start of a
line, which indicates the entire line is a comment, or at the end of a recipe line. If you put a '#' in literal, it will
not be considered a comment character. However, if you try to add a comment before a recipe line is potentially
complete, you will see an error when trying to run the recipe.

The next important item is the assignment operator. It is a left arrow created with a less-than and a dash - `<-`. Each
recipe line will start with either a variable, column or header followed by the assignment operator, followed by your
recipe to build the value that should go in the output.

The rest of the line consists of combinations of columns, variables, literals and functions. They can be combined using
a right arrow `->` which indicates whatever has been built should be passed along into the next thing (typically this is
going to be a function), or you can use `+` to combine the values on the left and the right of the plus. Recipes
execute strictly left-to-right. This means if you have more complex operations that would require out-of-order
processing, you should consider variables in order to simplify the operation into something that can be executed
linearly.

Let's take a look at some simple recipes. We'll be looking at single lines at a time, but a full recipe file will likely
consist of many of these:

`# full line comment`

This line will not do anything for the output, but it can be helpful to leave notes for yourself or others who might be
helping with maintaining the recipe.

`1 <- 1`

This is one of the simplest recipes. It means that in the output CSV, the first column will come from the first column
of the input CSV.

`1 <- 2`

In this recipe, we're moving a column over. Column 2 of the input file will be placed into column 1 of the output file.

`$first <- 3`

In this recipe, we're taking column 3 of each row and placing it into a variable that we can reference later for other
recipe lines. This can help with ensuring consistency or efficiency. Variables can also help to simplify more complex
expressions into expressions that can execute linearly, left-to-right.

`5 <- $first`

This recipe has the output column 5 receiving the value of the `$first` variable.

`!4 <- "fruit picked"`

This recipe sets a header for column 4 to the phrase `fruit picked`. The exclamation point in front of 4 indicates this
is a header definition, and the quotes around "fruits picked" mean we're using this value as provided.

While there is quite a lot that can be done with just the examples about, we can get even more power when utilizing
pipe (->) and join (+).

`1 <- 2 + 3 # concatenate 2 and 3 into 1`

This recipe reads the value of column 2, tacks on column 3 to the end of it and puts the result in column 1. I've also
included a comment at the end of the recipe to show how that can work as well.

`1 <- 1 -> uppercase`

This recipe leaves column 1 in the first position, but transforms the value through the uppercase function which
transforms any lowercase characters into uppercase.

`3 <- 3 + uppercase`

This one may be a bit tricky, but it's different from the example above. Instead of using a pipe operator, it's using
join. This means instead of just converting the column to uppercase, it combines the original value with the uppercased value. For
example if the original column value was "foo", then this recipe would result in "fooFOO". Using a pipe would result in
only "FOO".

`2 <- 2 -> lowercase + " particles"`

This recipe will lowercase column 2 and then add the word "particles" with a space. This result ends up in column 2 for
the output.

The next few recipes, while written differently, will all result in the same outcome.

```
1 <- 1 -> uppercase
1 <- uppercase(1)
1 <- uppercase(?)
```

This example introduced one new concept I haven't mentioned before -- the placeholder, represented by ?. This is
implicitly what's used from one step in a recipe to the next and passed into functions for which you do not provide
arguments.

This final example is a bit silly but may help illustrate the placeholder a bit.

`
1 <- 1 + ? + ? + uppercase(?)
`

Suppose column 1 contained "apple". This means after that lookup, the placeholder would contain "apple". Then we combine
the placeholder value with the next operation which is the placeholder value. That means after `1 + ?` the placeholder
would contain "appleapple". This is the new placeholder value. Then the next `+ ?` happens which takes "appleapple" and
combines it with the same, resulting in "appleappleappleapple". After the final operation, which takes the 4x-apple and
combines with an uppercase version of the same, the placeholder and the ultimate result will be "
appleappleappleappleAPPLEAPPLEAPPLEAPPLE" or, "apple" repeated 8 times, with the first 4 lowercase and the last 4
uppercase. I don't know why you'd ever want or need to do this, but... you could I guess.

There's more you can do, but it would be impossible to provide examples for all of them. Please see the functions
section for what the provided functions do to learn more about the possibilities.

Available Functions
==

Functions can be hard to understand at first, but the important thing to know is that they will all return a "string" or
character value. They may require some sort of input or possibly even configuration, and I'll try to document that here.
If the function can accept and return a string, it will also accept a column or variable or literal value. If it doesn't
care about what you provide, I'll indicate that with the `?` or placeholder value, but you can provide one of the other
value options. However, if it does matter what you provide, I'll document that as well, and will probably name the parameter
something to indicate its value. You still could provide a variable or column or literal value, but it would need to
match whatever the function is expecting.

Finally, if the function takes a final parameter of `?` then you can leave it off of your recipe, and the program will
automatically provide it. In fact, it will automatically provide the placeholder value for all arguments that you don't
explicitly provide. This may not be what you want though, and the output may not be correct, or a function that expects
certain input may result in an error. If a function does not need any parameters and won't use them if you provide them,
I'll indicate that with empty parens. You can leave those off too. Functions are case-insensitive when calling them, so you could use `Uppercase` or `uppercase` or even `UPPERCASE` in your recipes. In the docs below sometimes I'll use mixed case just to make the function naming easier to read.

* uppercase(?) - transforms characters in the value to uppercase - ex uppercase("apple") is APPLE.
* lowercase(?) - transforms characters in the value to lowercase - ex lowercase("LOWER") is lower.
* join(?) - This function joins whatever has happened on the left (or in the parameter) with the rest of the recipe on the right. CSV inserts this function automatically whenever you use the `+` operator.
* today() - returns today's date in YYYY-mm-dd format, ex 2021-08-30
* now() - returns the current date and time in RFC-3339 format, ex: `2021-08-30T18:22:13-06:00` 
* add(?, ?) - accepts two values that should be numerical and returns a string representing the sum of those two values.
  Providing non-numerical values will probably not do what you want. Remember, `add(2, 3)` is not 5, it's the sum of the values in columns 2 and 3.
* change(from, to, input) - If `from` is the same as `input` then it returns the `to` value. If it is not matching, then the original value returns.
* changei(from, to, input) - This works the same as change, but it is case-insensitive regarding the matching.
* ifEmpty(emptyVal, notEmptyVal, input) - If input is empty then `emptyVal` is returned, otherwise it returns the `notEmpty` value. Since recipes fill in missing values with the placeholder (?) automatically, if you want non-empty values to be retained, you can simply put `notEmpty(emptyVal)` in your recipe, and it will retain non-empty values unchanged.
* subtract(?, ?) - returns the value of the first parameter minus the second. All the caveats that apply to add apply
  here.
* multiply(?, ?) - returns the product of the two provided numerical values. If either are not numerical, an error will occur.
* divide(?, ?) - provides the result of first value divided by the second. They should of course be numbers, and the second value should not be zero unless you want to cause damage to the space-time continuum.
* numberFormat(digits, ?) - run this after add, subtract, multiply or divide to trim decimals. The `digits` parameter is how many digits after the decimal you want to keep.
* lineno() - this function returns the current line number
* mod(x, y) - returns the remainder of dividing x by y. Both arguments need to be integers. If they are not, an error will happen. If y is zero, an error will be returned.
* trim(?) - returns the argument with any leading or trailing white-space removed
* removeDigits(?) - strips all digit characters from the provided value
* firstChars(count, input) - Returns the first `count` characters of the input. If count is larger than the number of characters in input, it returns all the input.
* lastChars(count, input) - Returns the last `count` characters of the input. If the input is smaller than `count` then all of `input` will be returned. If the `count` parameter is not an integer or is negative, an error will occur.
* onlyDigits(?) - strips all characters except digits from the provided value
* normalize_date(format, date) - This function can accept a date in the provided `format` and return a string of that
  date in a format that other functions that need dates can utilize.
* formatDate(format, date) - Use this at the end of a line of date operations to get a date in a format that you want. Formatting is go style based on "Mon Jan 01, 2006 15:04:05-0700". It can recognize Monday or January if you want it spelled out, and 03 for 12-hour-time, as well as PM or pm if you want that included. The timezone is MST on that day, so MST will spell out the timezone, or America/Denver for the fully spelled out timezone. Incoming date should be normalized to RFC 3339 format first.
* formatDateF(format, date) - Similar to formatDate, this will take an incoming RFC3339 formatted date and return it in the go format specified date format. If it does not recognize the incoming value as RFC3339 format, then an error will occur and processing will stop.
* readDate(format, date) - Reads a date in a given format and returns it in RFC3339 format. Uses go format to specify how to read the date. If it does not recognize the incoming format, it will pass the input through unchanged. This allows you to chain more than one readDate if there are several formats you want to recognize.
* readDateF(format, date) - Reads a date in a given format and returns it in RFC3339 format. If the input does not match the given format, it returns an error which will cause processing to stop.
* if_after(after, not_after, date) - This function will return the `after` value if today is after the provided `date`,
  or the `not_after` value if today is before `date`.
* smartDate(date) - Tries to read a date in any reasonable format. If it cannot read as a date it will have an error. In this case, you may want to try specifying a format and using readDate. The return value will be a string of the date in RFC 3339 format if it was recognized as a date.
* isPast(past, future, date) - If the provided date is in the past, then the `past` arg is returned. If it's not, then it returns the `future` argument.
* isFuture(future, past, date) - If the provided date is in the future, then `future` arg is returned. Otherwise, the `past` arg is returned.
* only_digits(?) - returns all digit characters from the provided value
* trim(?) - removes whitespace from the provided value
* first_chars(num, ?) - returns the first `num` characters of a string
* last_chars(num, ?) - returns the last `num` characters of a string
* repeat(count, ?) - returns the input repeated `count` times, ex: `repeat(3, "apple")` is `appleappleapple`
* replace(search, replace, ?) - If it finds the `search` string within the input, it will be replaced with the `replace` string. If it's not found, then it returns the original input unchanged.
* power(num, power) - Returns `num` raised to the `power` parameter.  Both values must be numeric.  It will return a string representation of the number.

Public Recipes
==

If you're looking for public recipes or want to contribute your recipes for others to use, please take a look here: [https://github.com/dstockto/csv-chef-recipes](https://github.com/dstockto/csv-chef-recipes)

I'll be contributing recipes that I think others might use, and I encourage others to do so as well.
