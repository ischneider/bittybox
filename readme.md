# bittybox | Simple floating point calculator

`bittybox` is a small, underpowered floating point expression evaluator.

## Motivation

Performing safe, efficient user-supplied math operations.

## Features

Basic operators (+-*/^), grouping, predictable operator precedence, variables, some constants and common unary functions (sin, cos, etc).

## Grammar

    Expr --> Unit (Binary Unit)*
    Unit --> Number | "(" Expr ")" | "-" Unit | Ident
    Binary --> "+" | "-" | "*" | "/" | "^"

## Some Credits

Initial inspiration from the long abandoned https://github.com/marcak/calc/tree/master/

With advice from https://www.engr.mun.ca/~theo/Misc/exp_parsing.htm