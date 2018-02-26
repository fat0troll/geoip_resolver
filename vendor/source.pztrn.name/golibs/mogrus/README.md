# Mogrus

Logger thing built on top of github.com/sirupsen/logrus with ability to
create multiple loggers (e.g. console and file loggers) for one logger
instance.

The reason to create this handler was a need of logging things to both
console and, e.g., file, which is unsupported by logrus itself (you have
to create several loggers for each output).

## Example

```
package main

import (
    // stdlib
    "os"

    // tools
    "bitbucket.org/pztrn/mogrus"
)

func main() {
    l := mogrus.New()
    l.Initialize()
    log := l.CreateLogger("helloworld")
    log.CreateOutput("stdout", os.Stdout, true)

    // File output.
    file_output, err := os.Create("/tmp/hellorowld.log")
    if err != nil {
        log.Errorln("Failed to create file output:", err.Error())
    }
    log.CreateOutput("file /tmp/hellorowld.log", file_output, false)

    log.Println("Starting log experiment tool...")
    log.Debugln("Debug here!")
    log.Infoln("This is INFO level")
    log.Println("This is also INFO level.")
    log.Warnln("This is WARN.")
    log.Errorln("This is ERROR level.")
    log.Fatalln("We will exit here.")
}

```
