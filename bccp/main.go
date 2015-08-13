package main

import (
    "io"
    "io/ioutil"
    "fmt"
    "os"
    "github.com/codegangsta/cli"
    "strings"
)

func CopyFile(source string, dest string) (err error) {
    sourcefile, err := os.Open(source)
    if err != nil {
        return err
    }

    defer sourcefile.Close()

    destfile, err := os.Create(dest)
    if err != nil {
        return err
    }

    defer destfile.Close()

    _, err = io.Copy(destfile, sourcefile)
    if err == nil {
        sourceinfo, err := os.Stat(source)
        if err != nil {
            err = os.Chmod(dest, sourceinfo.Mode())
        }
    }

    return
}


func CopyDir(source string, dest string) (err error) {

    // get properties of source dir
    sourceinfo, err := os.Stat(source)
    if err != nil {
        return err
    }

    // create dest dir

    err = os.MkdirAll(dest, sourceinfo.Mode())
    if err != nil {
        return err
    }

    directory, _ := os.Open(source)
    objects, err := directory.Readdir(-1)

    for _, obj := range objects {
        sourcefilepointer := source + "/" + obj.Name()
        destinationfilepointer := dest + "/" + obj.Name()

        if obj.IsDir() {
            // create sub-directories - recursively
            err = CopyDir(sourcefilepointer, destinationfilepointer)
            if err != nil {
                fmt.Println(err)
            }
        } else {
            // perform copy
            err = CopyFile(sourcefilepointer, destinationfilepointer)
            if err != nil {
                fmt.Println(err)
            }
        }
    }
    return
}


func main() {
    app := cli.NewApp()
    app.Name = "bccp"
    app.Usage = "Copy file, directory to directories with prefix within directory"

    app.Flags = []cli.Flag {
  	    cli.StringFlag{
            Name: "source, s",
            Usage: "source file or directory",
        },
        cli.StringFlag{
            Name: "destination, d",
            Usage: "destination directory",
        },
        cli.StringFlag{
            Name: "prefix, p",
            Usage: "subdirectory prefix",
        },
        cli.StringFlag{
            Name: "subdirectory, sd",
            Usage: "subdirectory",
        },
    }

    app.Action = func(c *cli.Context) {
        source := c.String("source")
        destination := c.String("destination")
        prefix := c.String("prefix")
        subdirectory := c.String("subdirectory")
        
        if (source == "") {
            panic("%v\n", "source is not provied")
        }

        if (destination == "") {
            panic("%v\n", "destination is not provied")
        }

        if (prefix == "") {
            panic("%v\n", "prefix is not provied")
        }

        if (subdirectory == "") {
            panic("%v\n", "subdirectory is not provied")
        }
        
        src, err := os.Stat(source)
        if err != nil {
            panic(err)
        }
        
        dest, err := os.Stat(destination)
        if err != nil {
            panic(136, err)
        }

        if !dest.IsDir() {
            fmt.Println("Destination is not a directory")
            os.Exit(1)
        }

        list, err := ioutil.ReadDir(destination);
        if err != nil {
            panic(146, err)
        }

        for _, entry := range list {
            if strings.HasPrefix(entry.Name(), prefix) {
                if entry.IsDir() {
                    path := strings.Join([]string{dest.Name(), entry.Name(), subdirectory}, "/")
                    
                    if (!src.IsDir()) {
                        destfile := strings.Join([]string{path, src.Name()}, "/")
                        err = CopyFile(source, destfile)
                        if err != nil {
                            panic(158, err)
                        } else {
                            fmt.Println("File copied ", destfile)
                        }
                    } else {
                        err = CopyDir(source, path)
                        if err != nil {
                            panic(165, err)
                        } else {
                            fmt.Println("Directory copied")
                        }    
                    }
                }
            }
        }
    }

    app.Run(os.Args)
}