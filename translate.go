package main

import (
	_ "embed"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/fatih/color"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
	"unicode/utf8"
)

//go:embed tw.json
var s []byte

//go:embed ascii.txt
var ascii string

const (
	version = "v1.1.1"
)

// 跳過這些後綴文件
var ignoreExtMap = map[string]bool{
	".jpg":      true,
	".jpeg":     true,
	".png":      true,
	".gif":      true,
	".svg":      true,
	".exe":      true,
	".zip":      true,
	".tar":      true,
	".tar.gz":   true,
	".rar":      true,
	".mp3":      true,
	".mp4":      true,
	".avi":      true,
	".pdf":      true,
	".doc":      true,
	".docx":     true,
	".xls":      true,
	".xlsx":     true,
	".ppt":      true,
	".pptx":     true,
	".rm":       true,
	".mid":      true,
	".iso":      true,
	".DS_Store": true,
	".mod":      true,
	".sum":      true,
	".ttf":      true,
	".woff":     true,
	".woff2":    true,
	".wav":      true,
	".eot":      true,
	".otf":      true,
	".fon":      true,
	".font":     true,
	".ttc":      true,
}

// 跳過目錄
var defExcludeMap = map[string]bool{
	".git":         true,
	".idea":        true,
	".vscode":      true,
	"node_modules": true,
}

var (
	pathFlag    string
	extFlag     string
	excludeFlag string
	verFlag     bool
	startFlag   string // {"line":4066,"character":0}
	endFlag     string // {"line":4066,"character":0}

	wg sync.WaitGroup
)

// 中文正則字符串匹配
var chineseReg = regexp.MustCompile("[\u4e00-\u9fa5]")

type Position struct {
	Line      int `json:"line"`
	Character int `json:"character"`
}

func ReadPath(p string) []string {
	var paths []string

	fileInfo, err := os.Stat(p)
	if err != nil {
		log.Panicln("[tw]: ", err)
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
	excludeMap := make(map[string]string, len(excludeList)+len(defExcludeMap))
	for _, v := range excludeList {
		v = strings.TrimSpace(v)
		excludeMap[v] = v
	}
	for k := range defExcludeMap {
		excludeMap[k] = k
	}

	if fileInfo.IsDir() {
		err := filepath.Walk(p, func(path string, info fs.FileInfo, err error) error {
			if info.IsDir() {
				if _, ok := excludeMap[info.Name()]; ok {
					return filepath.SkipDir
				}
			}

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

func writeFile(p string, twMap map[string]string) {
	defer func() {
		wg.Done()
	}()

	b, err := os.ReadFile(p)
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

func Translate() {
	var twMap map[string]string
	err := json.Unmarshal(s, &twMap)
	if err != nil {
		panic(err)
	}

	fmt.Printf(" %v\n\n", color.CyanString("Processing..."))

	paths := ReadPath(filepath.Join(pathFlag))

	// 位置匹配
	if startFlag != "" && endFlag != "" {
		var start Position
		var end Position

		err := json.Unmarshal([]byte(startFlag), &start)
		if err != nil {
			panic(err)
		}
		err = json.Unmarshal([]byte(endFlag), &end)
		if err != nil {
			panic(err)
		}

		if len(paths) <= 0 {
			panic(err)
		}

		path := paths[0]
		b, err := os.ReadFile(path)
		if err != nil {
			fmt.Printf(
				"%v %v, %v\n",
				color.RedString("ERR"),
				color.YellowString(path),
				err,
			)
			return
		}

		content := string(b)
		// 將字符串按行分割成切片
		lines := strings.Split(content, "\n")
		if len(lines) < end.Line || end.Line < start.Line {
			fmt.Errorf("invalid line number")
			return
		}

		for i := start.Line - 1; i < end.Line; i++ {
			// 獲取當前行的長度，如果不足結束位置，則直接跳過本次循環
			lineLen := utf8.RuneCountInString(lines[i])
			line := []rune(lines[i])

			// 替換需要替換的部分
			// 截取的是同一行
			if i == start.Line-1 && i == end.Line-1 {
				fmt.Println("::1")
				str := line[start.Character:end.Character]
				lines[i] = chineseReg.ReplaceAllStringFunc(string(str), func(s2 string) string {
					if _, ok := twMap[s2]; ok {
						s2 = twMap[s2]
					}
					return s2
				})
				lines[i] = string(line[:start.Character]) + lines[i] + string(line[end.Character:])
			} else if i == start.Line-1 {
				fmt.Println("::2")
				str := line[start.Character:lineLen]
				lines[i] = chineseReg.ReplaceAllStringFunc(string(str), func(s2 string) string {
					if _, ok := twMap[s2]; ok {
						s2 = twMap[s2]
					}
					return s2
				})
				lines[i] = string(line[:start.Character]) + lines[i]
			} else if i == end.Line-1 {
				fmt.Println("::3")
				str := line[:end.Character]
				lines[i] = chineseReg.ReplaceAllStringFunc(string(str), func(s2 string) string {
					if _, ok := twMap[s2]; ok {
						s2 = twMap[s2]
					}
					return s2
				})
				lines[i] += string(line[end.Character:])
			} else {
				fmt.Println("::4")
				lines[i] = chineseReg.ReplaceAllStringFunc(lines[i], func(s2 string) string {
					if _, ok := twMap[s2]; ok {
						s2 = twMap[s2]
					}
					return s2
				})
			}
		}

		err = os.WriteFile(path, []byte(strings.Join(lines, "\n")), 0666)
		if err != nil {
			panic(err)
		}
	} else {
		wg.Add(len(paths))

		for _, p := range paths {
			go func(p string) {
				writeFile(p, twMap)
			}(p)
		}
		wg.Wait()
	}

	fmt.Printf("\n %v files, ", len(paths))
}

func main() {
	fmt.Println(ascii)

	flag.StringVar(&pathFlag, "path", ".", "--path")
	flag.StringVar(&extFlag, "ext", "", "--ext")
	flag.StringVar(&excludeFlag, "exclude", "", "--exclude")
	flag.StringVar(&startFlag, "start", "", "--start")
	flag.StringVar(&endFlag, "end", "", "--end")
	flag.BoolVar(&verFlag, "version", false, "--version")
	flag.Parse()

	if !verFlag {
		start := time.Now()
		Translate()

		n := time.Since(start).Seconds()
		fmt.Printf("Time: %vs\n", n)
	} else {
		fmt.Printf("Version %v \n", version)
	}
}
