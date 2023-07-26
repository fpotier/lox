My implementation of the [Lox](https://craftinginterpreters.com/the-lox-language.html) scripting language written in Go.

Lox is a programming language designed by [Robert Nystrom](http://stuffwithstuff.com/) in the book [Crafting Interpreters](https://craftinginterpreters.com/).

# Tests

The test suite in `test/official_tests` comes from the [original implementation](https://github.com/munificent/craftinginterpreters/tree/master/test)

## Current state

### Official test suite

<details>

| Feature              | Implementation     |
| -------------------- | ------------------ |
| assignment           | :x:                |
| benchmark            | :x:                |
| block                | :x:                |
| bool                 | :x:                |
| call                 | :x:                |
| class                | :x:                |
| closure              | :x:                |
| comments             | :x:                |
| constructor          | :x:                |
| expressions          | :x:                |
| field                | :x:                |
| for                  | :x:                |
| function             | :x:                |
| if                   | :x:                |
| inheritance          | :x:                |
| limit                | :x:                |
| logical operator     | :x:                |
| method               | :x:                |
| nil                  | :x:                |
| number               | :x:                |
| operator             | :x:                |
| print                | :x:                |
| regression           | :x:                |
| return               | :x:                |
| scanning             | :x:                |
| string               | :x:                |
| super                | :x:                |
| this                 | :x:                |
| variable             | :x:                |
| while                | :x:                |
| empty file           | :white_check_mark: |
| precedence           | :white_check_mark: |
| unexpected character | :x:                |

</details>