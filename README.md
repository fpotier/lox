My implementation of the [Lox](https://craftinginterpreters.com/the-lox-language.html) scripting language written in Go.

Lox is a programming language designed by [Robert Nystrom](http://stuffwithstuff.com/) in the book [Crafting Interpreters](https://craftinginterpreters.com/).

# Tests

The test suite in `test/official_tests` comes from the [original implementation](https://github.com/munificent/craftinginterpreters/tree/master/test)

## Current state

### Official test suite

<details>

| Feature              | Implementation     |
| -------------------- | ------------------ |
| assignment           | :white_check_mark: |
| benchmark            | :white_check_mark: |
| block                | :white_check_mark: |
| bool                 | :white_check_mark: |
| call                 | :white_check_mark: |
| class                | :white_check_mark: |
| closure              | :white_check_mark: |
| comments             | :white_check_mark: |
| constructor          | :white_check_mark: |
| field                | :white_check_mark: |
| for                  | :white_check_mark: |
| function             | :white_check_mark: |
| if                   | :white_check_mark: |
| inheritance          | :white_check_mark: |
| limit                | :x:                |
| logical operator     | :white_check_mark: |
| method               | :white_check_mark: |
| nil                  | :white_check_mark: |
| number               | :x:                |
| operator             | :white_check_mark: |
| print                | :white_check_mark: |
| regression           | :white_check_mark: |
| return               | :white_check_mark: |
| scanning             | :white_check_mark: |
| string               | :white_check_mark: |
| super                | :white_check_mark: |
| this                 | :white_check_mark: |
| variable             | :white_check_mark: |
| while                | :white_check_mark: |
| empty file           | :white_check_mark: |
| precedence           | :white_check_mark: |
| unexpected character | :white_check_mark: |

</details>