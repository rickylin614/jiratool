default:
	echo 'please input command'

# 編譯mac執行檔
mac:
	fyne-cross darwin -app-id="com.example.myapp" -image="fyneio/fyne-cross-images:darwin"

# 編譯windows執行檔
# go install github.com/fyne-io/fyne-cross@v1.3.0
windows:
	fyne-cross windows
	copy .\conf.yml .\fyne-cross\bin\windows-amd64

# 將設定檔移至編譯好得windows路徑
.PHONY: cp
cp:
	copy .\conf.yml .\fyne-cross\bin\windows-amd64

# 將設定檔移至編譯好得mac路徑
cp_mac:
	copy ./conf.yml ./fyne-cross/bin/darwin-amd64

# 本地啟動
run:
	go run main.go
