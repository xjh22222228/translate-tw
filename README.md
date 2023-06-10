
将文件内的中文简体转换成繁体。





## 安装
Shell (Mac):
```bash
curl -fsSL https://raw.githubusercontent.com/xjh22222228/translate-tw/main/install.sh | bash

# or
curl -fsSL https://cdn.jsdelivr.net/gh/xjh22222228/translate-tw@main/install.sh | bash
```

Windows:

[Download](https://github.com/xjh22222228/translate-tw/releases/latest/download/tw_windows_amd64.zip)



## 使用
默认会跳过 `.jpg` / `.png` 等非文本文件。

| 参数     | 描述              |
| ------- |------------------ |
| --path  | 指定文件或目录  |
| --ext   | 过滤文件后缀，默认不过滤  |
| --exclude   | 跳过指定目录  |
| --start   | 匹配開始行數和位置    |
| --end     | 匹配結束行數和位置    |
| --version   | 打印版本号  |



```bash
$ tw --path="./src"

# 指定文件后缀
$ tw --path="./src" --ext=".js|.jsx|go"

# 跳过指定目录
$ tw --path="./src" --exclude="node_modules|dist"

# 匹配位置, 必須傳入 --start 和 --end 參數
$ tw --path="./src" --start='{"line":4066,"character":0}' --end='{"line":4066,"character":0}'
```


## LICENSE
[MIT](./LICENSE)
