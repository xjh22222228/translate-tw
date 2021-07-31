package main

import (
    _ "embed"
    "encoding/json"
    "flag"
    "fmt"
    "github.com/fatih/color"
    "io/fs"
    "io/ioutil"
    "log"
    "os"
    "path/filepath"
    "regexp"
    "strings"
    "sync"
    "time"
)

//go:embed tw.json
var s []byte

//go:embed ascii.txt
var ascii string

const (
    version = "v1.0.0"
)

// 跳过这些后缀文件
var ignoreExtMap = map[string]bool{
    ".jpg": true,
    ".jpeg": true,
    ".png": true,
    ".gif": true,
    ".svg": true,
    ".exe": true,
    ".zip": true,
    ".tar": true,
    ".tar.gz": true,
    ".rar": true,
    ".mp3": true,
    ".mp4": true,
    ".avi": true,
    ".pdf": true,
    ".doc": true,
    ".docx": true,
    ".xls": true,
    ".xlsx": true,
    ".ppt": true,
    ".pptx": true,
    ".rm": true,
    ".mid": true,
    ".iso": true,
}

var (
    pathFlag string
    extFlag string
    excludeFlag string
    verFlag bool

    wg sync.WaitGroup
)

var chineseReg = regexp.MustCompile("[\u4e00-\u9fa5]")

func ReadPath(p string) []string {
    var paths []string

    fileInfo, err := os.Stat(p)
    if err != nil {
        log.Panicln("Not path：", err)
    }

    // ext
    extList := strings.Split(extFlag, "|")
    extMap := make(map[string]string, len(extList))
    for _, v := range extList {
        v = strings.TrimSpace(v)
        extMap[v] = v
    }

    // exclude
    excludeList := strings.Split(excludeFlag, "|")
    excludeMap := make(map[string]string, len(excludeList))
    for _, v := range excludeList {
        v = strings.TrimSpace(v)
        excludeMap[v] = v
    }

    if fileInfo.IsDir() {
        err := filepath.Walk(p, func(path string, info fs.FileInfo, err error) error {
            if !info.IsDir() {
                ext := filepath.Ext(path)
                _, isIgnore := ignoreExtMap[ext]
                if _, ok := extMap[ext]; (ok || extFlag == "") && !isIgnore {
                    paths = append(paths, path)
                }
            }

            return err
        })

        if err != nil {
            log.Panicln("Walk：", err)
        }
    } else {
        paths = append(paths, p)
    }

    return paths
}

func writeFile(p string, twMap map[string]string)  {
    defer func() {
        wg.Done()
    }()

    b, err := ioutil.ReadFile(p)
    if err != nil {
        fmt.Printf(
            "%v %v, %v\n",
            color.RedString("ERR"),
            color.YellowString(p),
            err,
        )
        return
    }

    content := string(b)
    f := chineseReg.ReplaceAllStringFunc(content, func(s2 string) string {
        if _, ok := twMap[s2]; ok {
            s2 = twMap[s2]
        }
        return s2
    })
    err = os.WriteFile(p, []byte(f), 0666)
    if err != nil {
        fmt.Printf(
            "%v %v, %v\n",
            color.RedString("ERR"),
            color.YellowString(p),
            err,
        )
    } else {
        fmt.Printf(
            "%v %v\n",
            color.GreenString(" OK"),
            color.YellowString(p),
        )
    }
}

func Translate()  {
    var twMap map[string]string
    err := json.Unmarshal(s, &twMap)
    if err != nil {
        panic(err)
    }

    fmt.Printf(" %v\n\n", color.CyanString("Processing..."))

    paths := ReadPath(filepath.Join(pathFlag))
    wg.Add(len(paths))

    for _, p := range paths {
        go func(p string) {
            writeFile(p, twMap)
        }(p)
    }
    wg.Wait()

    fmt.Printf("\n %v files, ", len(paths))
}

func main()  {
    fmt.Println(ascii)

    flag.StringVar(&pathFlag, "path", ".", "--path")
    flag.StringVar(&extFlag, "ext", "", "--ext")
    flag.StringVar(&excludeFlag, "exclude", "", "--exclude")
    flag.BoolVar(&verFlag, "version", false, "--version")
    flag.Parse()

    if !verFlag {
        start := time.Now()
        Translate()

        n := time.Since(start).Seconds()
        fmt.Printf("Time: %vs", n)
    } else {
        fmt.Printf("Version %v \n", version)
    }
}