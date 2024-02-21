package fynetool

import (
	"errors"
	"fmt"
	"io"
	"jiratool/conf"
	"jiratool/gittool"
	"jiratool/lib"
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

const title = "Git工具@ricky in 2024/02"
const admintxt = "temp1.txt"
const sourcetxt = "temp2.txt"

// 初始化視窗
func InitFyneApp() fyne.Window {
	a := app.New()
	// 設置字體
	a.Settings().SetTheme(&lib.MyTheme{})

	// 初始化視窗顯示內容
	w := a.NewWindow(title)

	// 設定寬高
	w.Resize(fyne.NewSize(1024, 768))

	return w
}

// 產生各個元件
func SetttingWidget(w fyne.Window) fyne.Window {

	input1 := widget.NewEntry()
	input2 := widget.NewEntry()

	inputCon1 := container.NewHScroll(input1)
	inputCon1.SetMinSize(fyne.NewSize(800, 20))
	inputCon2 := container.NewHScroll(input2)
	inputCon2.SetMinSize(fyne.NewSize(800, 20))
	loadPath(input1, admintxt)
	loadPath(input2, sourcetxt)

	fileSelect1 := newButtonWidget(w, input1, admintxt)
	fileSelect2 := newButtonWidget(w, input2, sourcetxt)

	clearButton1 := widget.NewButton("Clear", func() {
		input1.SetText("")
		savePath(input1.Text, admintxt)
	})
	clearButton2 := widget.NewButton("Clear", func() {
		input2.SetText("")
		savePath(input2.Text, sourcetxt)
	})

	progressBar := widget.NewProgressBarInfinite()
	progressBar.Hide()
	progressBarCon := container.NewHScroll(progressBar)
	progressBarCon.SetMinSize(fyne.NewSize(1000, 20))

	errLabel := widget.NewLabel("")
	errLabel.Hidden = true

	actionButton1 := widget.NewButton("admin update", func() {
		progressBar.Show()
		defer progressBar.Hide()
		err := mergeCsv(input1.Text, conf.GetConfig().AdminPath, conf.GetConfig().AdminFileName)
		if err == nil {
			errLabel.Hidden = true
		} else {
			errLabel.SetText("執行失敗或檔案沒有更新!!")
			errLabel.Hidden = false
		}
	})
	actionButton2 := widget.NewButton("source update", func() {
		progressBar.Show()
		defer progressBar.Hide()
		err := mergeCsv(input2.Text, conf.GetConfig().WebPath, conf.GetConfig().WebFileName)
		if err == nil {
			errLabel.Hidden = true
		} else {
			errLabel.SetText("執行失敗或檔案沒有更新!!")
			errLabel.Hidden = false
		}
	})

	w.SetContent(container.NewVBox(
		container.NewHBox(widget.NewLabel(" admin"),
			fileSelect1,
			inputCon1,
			clearButton1,
		),
		container.NewHBox(widget.NewLabel("source"),
			fileSelect2,
			inputCon2,
			clearButton2,
		),
		container.NewHBox(
			actionButton1, actionButton2,
		),
		progressBarCon,
		widget.NewLabel("備註: 需裝過git並且指令可正常Git push專案"),
		errLabel,
	))

	return w
}

func newButtonWidget(w fyne.Window, input *widget.Entry, saveTxt string) *widget.Button {
	return widget.NewButton("Select File", func() {
		fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err == nil && reader != nil {
				input.SetText(reader.URI().Path())
				savePath(input.Text, saveTxt)
			}
		}, w)

		filter := &customFileFilter{extensions: []string{".csv"}}
		fd.SetFilter(filter)

		fd.Resize(fyne.NewSize(800, 600))
		fd.Show()
	})
}

type customFileFilter struct {
	extensions []string
}

func (cff *customFileFilter) Matches(uri fyne.URI) bool {
	// Exclude hidden files and filter by extensions
	return !strings.HasPrefix(uri.Name(), ".") && cff.hasValidExtension(uri)
}

func (cff *customFileFilter) hasValidExtension(uri fyne.URI) bool {
	// Check if the file has a valid extension
	for _, ext := range cff.extensions {
		if strings.HasSuffix(uri.Name(), ext) {
			return true
		}
	}
	return false
}

func savePath(path string, filename string) {
	file, err := os.Create(filename)
	if err != nil {
		return
	}
	defer file.Close()

	_, err = file.WriteString(path)
	if err != nil {
		return
	}
}

func loadPath(entry *widget.Entry, filename string) {
	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer file.Close()

	buf := make([]byte, 1024)
	n, _ := file.Read(buf)
	entry.SetText(string(buf[:n]))
}

func mergeCsv(replaceFile, filePath, fileName string) error {
	var repoName string
	var err error
	defer func() {
		if repoName != "" {
			checkAndRemoveDirectory(repoName)
		}
	}()

	if replaceFile == "" || filePath == "" || fileName == "" {
		return errors.New("執行失敗!!")
	}

	if repoName, err = gittool.GitClone(filePath); err != nil {
		return err
	}

	if err = copyFile(replaceFile, fmt.Sprintf("./%s/%s", repoName, fileName)); err != nil {
		return err
	}

	if err = gittool.GitAdd(repoName, fileName); err != nil {
		return err
	}

	if err = gittool.GitCommitAndPush(repoName); err != nil {
		return err
	}

	if err = gittool.AddTag(repoName); err != nil {
		return err
	}

	return nil
}

func showBusyIndicator(progressBar *widget.ProgressBarInfinite) {
	progressBar.Show()
}

func hideBusyIndicator(progressBar *widget.ProgressBarInfinite) {
	progressBar.Hide()
}

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		fmt.Println("Error opening source file:", err)
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		fmt.Println("Error creating destination file:", err)
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		fmt.Println("Error copying file:", err)
		return err
	}
	return nil
}

func checkAndRemoveDirectory(path string) error {
	if directoryExists(path) {
		return removeDirectory(path)
	} else {
		return fmt.Errorf("Directory '%s' does not exist", path)
	}
}

func directoryExists(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func removeDirectory(path string) error {
	return os.RemoveAll(path)
}
