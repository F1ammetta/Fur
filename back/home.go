package back

import (
	"fmt"
	"html/template"
	"math"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
)

type FolderData struct {
	Path   string
	Folder string
	Drives []Drive
}

type File struct {
	Url   string
	Type  string
	Size  float64
	IsDir bool
	Name  string
	Noext string
	Ext   string
	Hide  bool
}

type Drive struct {
	Name string
	Path string
}

var parsedFiles []File

// function to read file and print full relative path, preferably recursively

// func Urls(c *gin.Context) []string {
// 	return
// }

func grid(c *gin.Context) {

	if dir == "" {
		tmpl, err := template.ParseFiles(run_dir + "/templates/blank.html")
		if err != nil {
			fmt.Println("Error parsing template: ", err)
			c.String(http.StatusInternalServerError, fmt.Sprintf("Error parsing template: %v", err))
			return
		}
		tmpl.Execute(c.Writer, nil)
		return
	}

	subdir := c.Param("dir")
	more_subdir := c.Param("deer")

	if subdir != "" {
		if !strings.HasPrefix(more_subdir, subdir) {
			subdir = filepath.Join(subdir, more_subdir)
		} else {
			subdir = more_subdir
		}
		subdir = strings.ReplaceAll(subdir, "\\", "/")
		dir = abs_dir + string(os.PathSeparator) + subdir
	} else {
		dir = abs_dir
	}

	// subdir = strings.ReplaceAll(subdir, "/", "\\")

	tmpl, err := template.ParseFiles(run_dir + "/templates/grid.html")

	if err != nil {
		fmt.Println("Error parsing template: ", err)
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error parsing template: %v", err))
		return
	}

	// read files
	ReadFiles(dir, c)
	// trim urls
	parsedFiles = TrimDir(abs_dir)
	parsedFiles = SortType(parsedFiles)
	for i, file := range parsedFiles {
		if i > len(parsedFiles)-1 {
			break
		}
		if strings.Contains(file.Url, "\\") {
			parsedFiles[i].Url = strings.ReplaceAll(file.Url, "\\", "/")
		}
		if file.IsDir {
			if 0 != strings.Compare(abs_dir, "/") {
				parsedFiles[i].Url = strings.ReplaceAll(parsedFiles[i].Url, abs_dir, "")
			}
			parsedFiles[i].Url = strings.TrimPrefix(parsedFiles[i].Url, "/")

			// if file name is previews, remove it
			if strings.Contains(parsedFiles[i].Url, "previews") {
				parsedFiles = append(parsedFiles[:i], parsedFiles[i+1:]...)
			}
		}
		parsedFiles[i].Hide = file.Name[0] == '.'
	}
	data := struct {
		Files []File
	}{
		Files: parsedFiles,
	}
	tmpl.Execute(c.Writer, data)
}

func SortType(files []File) []File {

	// sort by extension first and then by type, so videos go before images, etc., then by name

	sort.Slice(files, func(i, j int) bool {
		if files[i].IsDir && !files[j].IsDir {
			return true
		}
		if !files[i].IsDir && files[j].IsDir {
			return false
		}
		if files[i].IsDir && files[j].IsDir {
			return files[i].Name < files[j].Name
		}
		if files[i].Ext == files[j].Ext {
			if files[i].Type == files[j].Type {
				return files[i].Name < files[j].Name
			}
			return files[i].Type < files[j].Type
		}
		return files[i].Ext < files[j].Ext
	})

	return files
}

func ReadFiles(path string, c *gin.Context) {
	parsedFiles = []File{}
	files, err := os.ReadDir(path)
	if err != nil {
		fmt.Println("Error reading files: ", err)
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error parsing template: %v", err))
		return
	}
	for _, file := range files {
		filename := file.Name()
		file, err := file.Info()

		if err != nil {
			fmt.Println("Error getting file info: ", err)
			c.String(http.StatusInternalServerError, fmt.Sprintf("Error parsing template: %v", err))
			return
		}

		sizeMB := float64(file.Size()) / 1000000
		//round to 2 decimal places
		sizeMB = math.Round(sizeMB*100) / 100

		var noext string
		filename_split := strings.Split(filename, ".")
		if file.IsDir() {
			noext = filename
		} else if len(filename_split) > 1 {
			noext = strings.TrimSuffix(filename, filepath.Ext(filename))
		} else {
			noext = filename
		}
		parsedFiles = append(parsedFiles, File{
			Url:   path + "\\" + filename,
			Type:  file.Mode().String(),
			IsDir: file.IsDir(),
			Size:  sizeMB,
			Name:  filename,
			Noext: noext,
			Ext:   strings.ToUpper(strings.TrimPrefix(filepath.Ext(filename), ".")),
		})
	}
}

func contains(slice []string, str string) bool {
	for _, v := range slice {
		if v == str {
			return true
		}
	}
	return false
}

func folder(c *gin.Context) {
	folder := c.Param("folder") // render templat

	more_folder := c.Param("path")

	og_folder := folder

	// more_folder may contain folder, remove it

	folderComponents := strings.Split(folder, "/")
	moreFolderComponents := strings.Split(more_folder, "/")

	// Create a slice to hold the unique components
	uniqueFolderComponents := make([]string, 0)

	// Add the components of folder to the slice
	for _, component := range folderComponents {
		if component != "" && !contains(uniqueFolderComponents, component) {
			uniqueFolderComponents = append(uniqueFolderComponents, component)
		}
	}

	// Add the components of more_folder to the slice
	for _, component := range moreFolderComponents {
		if component != "" && !contains(uniqueFolderComponents, component) {
			uniqueFolderComponents = append(uniqueFolderComponents, component)
		}
	}

	// Join the unique components back into a string
	folder = strings.Join(uniqueFolderComponents, "/")

	if strings.HasSuffix(folder, og_folder) && folder != og_folder {
		folder = strings.TrimSuffix(folder, og_folder)
	}

	folder = strings.ReplaceAll(folder, "\\", "/")

	drives := []Drive{}
	var command *exec.Cmd

	command = exec.Command("bash", "-c", "lsblk -o label,MOUNTPOINTS -n | grep -v 'boot' | awk '$2!=\"\"'")
	// command = exec.Command("powershell", "-c", "Get-PSDrive -PSProvider FileSystem | Select-Object -ExpandProperty Root")
	// command.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	output, err := command.Output()

	if err != nil {
		fmt.Println("Error getting drives: ", err)
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error parsing template: %v", err))
		return
	}

	lines := strings.Split(string(output), "\n")

	user := os.Getenv("USER")
	if os.PathSeparator == '\\' {
		drives = append(drives, Drive{Name: "Home", Path: "C:\\Users\\" + user})
	} else {
		drives = append(drives, Drive{Name: "Home", Path: "/home/" + user})
		for _, line := range lines {
			if line == "" {
				continue
			}
			// path is last element
			// name is everything else
			split := strings.Split(line, " ")
			path := split[len(split)-1]
			name := strings.Join(split[:len(split)-1], " ")
			drives = append(drives, Drive{
				Name: name,
				Path: path,
			})
		}
	}

	if folder == "favicon.ico" {
		c.File(run_dir + "/static/favicon.ico")
	}

	tmpl, err := template.ParseFiles(run_dir + "/templates/index.html")
	if err != nil {
		fmt.Println("Error parsing template: ", err)
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error parsing template: %v", err))
		return
	}

	path := filepath.Clean(abs_dir + string(os.PathSeparator) + folder)

	if len(path) > 20 {
		path = path[:20] + "..."
	}

	data := FolderData{
		Path:   path,
		Folder: "/" + folder,
		Drives: drives,
	}

	tmpl.Execute(c.Writer, data)
}

func previews(c *gin.Context) {
	path := c.Request.URL.Path
	path = strings.TrimSuffix(path, "_preview.png")
	ext := strings.Split(path, ".")[len(strings.Split(path, "."))-1]

	switch ext {
	case "jpg", "jpeg", "png", "gif", "bmp", "tiff", "tif":
		ext = "image"
	case "mp4", "mov", "avi", "wmv", "flv", "mkv", "webm", "m4v":
		ext = "video"
	case "mp3", "wav", "ogg", "flac", "m4a", "wma", "aac":
		ext = "audio"
	case "doc", "docx", "txt", "pdf", "rtf", "odt", "md":
		ext = "document"
	case "zip", "rar", "7z", "tar", "gz", "bz2", "xz", "z":
		ext = "archive"
	case "":
		ext = "folder"
	default:
		ext = "unknown"
	}

	c.File(run_dir + "/static/previews/" + ext + ".svg")

	return
}

func files(abs_path string, c *gin.Context) {

	abs_path = filepath.Clean(abs_path)

	// remove leading slash
	// convert to absolute path
	// run command to open file with default program platform agnostic
	go func() {
		cmd := exec.Command("xdg-open", abs_path)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			fmt.Println("Error opening file: ", err)
			return
		}
		return
	}()
	// go func() {
	// 	cmd := exec.Command("powershell", "-c", "Start-Process", "\""+absPath+"\"")
	// 	cmd.Stdout = os.Stdout
	// 	cmd.Stderr = os.Stderr
	// 	fmt.Println("Opening file: ", absPath)
	// 	err := cmd.Run()
	// 	if err != nil {
	// 		fmt.Println("Error opening file: ", err)
	// 		return
	// 	}
	// }()
	c.Status(200)
}

// func trim dir to get rid of trailing dir, but keep the subdirs
func TrimDir(path string) []File {
	trimmedFiles := []File{}
	for _, file := range parsedFiles {

		noext := file.Noext
		// trim to 20 characters add ellipsis
		if len(noext) > 20 {
			noext = noext[:20] + "..."
		}
		trimmedFiles = append(trimmedFiles, File{
			Url:   strings.TrimPrefix(file.Url, path),
			Type:  file.Type,
			Size:  file.Size,
			IsDir: file.IsDir,
			Name:  file.Name,
			Noext: noext,
			Ext:   strings.ToUpper(strings.TrimPrefix(filepath.Ext(file.Name), ".")),
		})
		//replace backslash with forward slash
	}
	return trimmedFiles
}

// func home(c *gin.Context) {
//
// 	// render template
// 	tmpl, err := template.ParseFiles(run_dir + "/templates/index.html")
// 	if err != nil {
// 		fmt.Println("Error parsing template: ", err)
// 		c.String(http.StatusInternalServerError, fmt.Sprintf("Error parsing template: %v", err))
//		return
// 	}
//
// 	data := FolderData{
// 		Folder: "",
// 		Home:   true,
// 	}
//
// 	tmpl.Execute(c.Writer, data)
// }
