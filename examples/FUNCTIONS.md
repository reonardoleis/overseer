# Functions

Function commands allow to create JavaScript functions and call them with args, which will be passed as args parameter into the created function. 

## Commands

### `!fncreate <name>` <br> `<code>`
Creates a function with the specified name and code. The code must be a valid JavaScript function body, without the function keyword and without brackets. The function will be called with the args parameter, which is an array of arguments passed to the function (on !fnrun command).
**Note:**: in order to work properly, the command must be executed exactly as shown above, with the code parameter on the next line (single line break), and the closing backtick on the line after that.

### `!fnrun <name> <args...>`
Runs the function with the specified name and arguments. The arguments will be passed as an array named **args** to the function, so they are accessible as **args[0]**, **args[1]**, etc.

## Example

### fncreate
`!fncreate add` <br> `return args[0] + args[1];`

### fnrun
`!fnrun add 1 2` <br> `> 3`

