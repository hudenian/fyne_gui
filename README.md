## 构建不带图形界面窗口
Go tool arguments
```
-ldflags="-s -w -H windowsgui"
```

## 需要设置中文字体
```go
func init() {
	//设置中文字体
	fontPaths := findfont.List()
	for _, path := range fontPaths {
		if strings.Contains(path, "msyh.ttf") || strings.Contains(path, "simhei.ttf") || strings.Contains(path, "simsun.ttc") || strings.Contains(path, "simkai.ttf") {
			os.Setenv("FYNE_FONT", path)
			break
		}
	}
}
```