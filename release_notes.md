# Release Notes 

V1.01 - September 14, 2021
* [#17](https://github.com/dstockto/csv-chef/issues/17) - Recipes that define the same column, header or variable more than once will now result in a parse error. Previous behavior was that the latter would override and previous definition silently which could be confusing if you accidentally left in identity column recipes.
* [#16](https://github.com/dstockto/csv-chef/issues/16) - Running identity on UTF-8 with BOM (byte order marker) files was resulting in the header value for the first column containing those BOMs which would show up as unprintable characters in some editors. The BOM is removed from the column header names in the comments for the column and column header recipes now.
