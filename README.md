# checkfatal
Static analyzer to search for dead code after `log.Fatal` and other functions that call `os.Exit`

Go is fantastic at fidning dead code and won't build if there are unused imports or variables, but if a function in one of your dependencies ultimately calls `os.Exit,` any code you have following that will not be executed ever but the compiler does not pick up that dead code. `fatalcheck` seeks to address that so that you can find potential dead code.

Code such as 

```
package main

import (
  "fmt"
  "log"
)

func main() {
  log.Fatal("This will write to the log and immediately exit the application")
  fmt.Println("You will never see this line printed.") // This won't ever execute but the compiler doesn't pick it up
}
```
