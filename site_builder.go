package main

import (
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"nikastudio/files"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var currentFolder string
var version = "2019-04-30-A"
var isSilent bool

var baseTmpl = new(baseTemplate)

type baseTemplate struct {
	baseTemplate *template.Template
	templates    map[string]string
}

func fromTemplate(w io.Writer, path string) {

	contentData, ok := baseTmpl.templates[path]
	if !ok {
		log.Panicf("can't find content for path %s", path)
	}

	tmpl, err := baseTmpl.baseTemplate.Clone()
	if err != nil {
		log.Panicf("can't create clone from base template: %v", err)
	}

	template.Must(tmpl.Parse(contentData))

	executionErr := tmpl.ExecuteTemplate(w, "base", path)
	if executionErr != nil {
		panic(executionErr)
	}
}

func images() []string {
	result := make([]string, 0)

	images, err := ioutil.ReadDir(filepath.Join(currentFolder, "static", "images", "portfolio", "thumbnails"))
	if err != nil {
		panic(err)
	}

	for _, img := range images {
		result = append(result, img.Name())
	}
	return result
}

func projects() []string {
	result := make([]string, 0)

	projects, err := ioutil.ReadDir(filepath.Join(currentFolder, "static", "images", "projects"))
	if err != nil {
		panic(err)
	}

	for _, p := range projects {
		result = append(result, p.Name())
	}
	return result
}

func imagesForProject(prj string) []string {
	result := make([]string, 0)

	prjPath := filepath.Join(currentFolder, "static", "images", "projects", truncProj(prj), "thumbnails")
	prjFiles, err := ioutil.ReadDir(prjPath)
	if err != nil {
		log.Panicf("can't read folder %s for project %s, error: %v", prjPath, prj, err)
	}

	for _, f := range prjFiles {
		result = append(result, f.Name())
	}

	return result
}

func truncProj(path string) string {
	noExt := strings.Replace(path, ".html", "", -1)
	return strings.Replace(noExt, "projects-", "", -1)
}

func isGallery(page string) bool {
	return page == "index.html" || strings.HasPrefix(page, "projects")
}

func isProject(page string) bool {
	return strings.HasPrefix(page, "projects")
}

func prepareTemplates() {
	layouts, err := files.ReadFiles(filepath.Join(currentFolder, "templates", "layout"))
	if err != nil {
		panic(err)
	}

	baseTmpl.baseTemplate = template.New("main")

	baseTmpl.baseTemplate.Funcs(
		template.FuncMap{
			"images":           images,
			"projects":         projects,
			"imagesForProject": imagesForProject,
			"truncProj":        truncProj,
			"isGallery":        isGallery,
			"isProject":        isProject,
		})

	for l := range layouts {
		if l == "projects_content.tmpl.html" {
			continue
		}
		template.Must(baseTmpl.baseTemplate.Parse(layouts[l]))
		logItf("layout %v parsed", l)
	}

	buildProjectsTmpl()

	baseTmpl.templates, err = files.ReadFiles(filepath.Join(currentFolder, "templates", "content"))
	if err != nil {
		panic(err)
	}
	resizePortfolio()
	resizeProjects()

	logIt("templates map prepared")
}

func resizePortfolio() {
	originPath := filepath.Join(currentFolder, "static", "images", "portfolio", "originals")
	thmbPath := filepath.Join(currentFolder, "static", "images", "portfolio", "thumbnails")
	origs, err := ioutil.ReadDir(originPath)
	if err != nil {
		log.Panicf("can't read files from portfolio originals %s,  error: %v", originPath, err)
	}
	for _, origImg := range origs {
		err := files.Thumbnail(
			filepath.Join(originPath, origImg.Name()),
			filepath.Join(thmbPath, origImg.Name()))
		if err != nil {
			log.Panicf("can't resize portfolio file %s  error: %v", origImg.Name(), err)
		}
		log.Printf("portfolio resize image ready-> %s", filepath.Join(thmbPath, origImg.Name()))
	}
}

func resizeProjects() {
	projectsPath := filepath.Join(currentFolder, "static", "images", "projects")
	projects, err := ioutil.ReadDir(projectsPath)
	if err != nil {
		log.Panicf("cant' read projects folder %s  error: %v", projectsPath, err)
	}

	for _, prj := range projects {
		originPath := filepath.Join(projectsPath, prj.Name(), "originals")
		thmbPath := filepath.Join(projectsPath, prj.Name(), "thumbnails")

		origs, err := ioutil.ReadDir(originPath)
		if err != nil {
			log.Panicf("can't read files from project originals %s,  error: %v", originPath, err)
		}
		for _, origImg := range origs {
			err := files.Thumbnail(
				filepath.Join(originPath, origImg.Name()),
				filepath.Join(thmbPath, origImg.Name()))
			if err != nil {
				log.Panicf("can't resize projects file %s  error: %v", origImg.Name(), err)
			}
			log.Printf("resize project image ready-> %s", filepath.Join(thmbPath, origImg.Name()))
		}
	}
}

func buildProjectsTmpl() {
	for _, p := range projects() {
		err := files.CreateProjectFile(currentFolder, p)
		if err != nil {
			log.Panicf("can't create project file %s, error: %v", p, err)
		}
		logItf("%v project file created", p)
	}

}

func validateCurrentFolder(folder string) {
	info, err := os.Stat(folder)
	if err != nil {
		log.Panic(err)
	}

	if !info.IsDir() {
		log.Panic(folder + " is not directory")
	}

	info, err = os.Stat(folder + "/static")
	if err != nil {
		log.Panicf("can't find static folder at %s, error: %v", folder+"/static", err)
	}

	if !info.IsDir() {
		log.Panic(folder + " is not directory")
	}
}

func serve(port string) {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		path := request.URL.Path

		logItf("requested %v ", path)

		if path == "/" {
			path = "index.html"
		}
		fromTemplate(writer, strings.Replace(path, "/", "", 1))
	})

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(currentFolder+"/static/"))))

	logItf("web server is ready on port %v", port)
	panic(http.ListenAndServe(port, nil))
}

func generate() {

	targetFolder := "docs"
	targetFolderPath := filepath.Join(currentFolder, targetFolder)

	err := os.RemoveAll(targetFolderPath)

	err = os.Mkdir(targetFolderPath, 0755)
	if err != nil {
		logItf("can't create dir %v, error: %v \n", targetFolderPath, err)
	}

	err = files.CopyDir(
		filepath.Join(currentFolder, "static"),
		filepath.Join(currentFolder, targetFolder, "static"))

	if err != nil {
		logItf("can't copy folders: %v \n", err)
	}

	tmplFiles, filesErr := ioutil.ReadDir(filepath.Join(currentFolder, "templates", "content"))
	if filesErr != nil {
		logItf("can't read tmplFiles in directory: %v \n", filesErr)
	}

	for _, file := range tmplFiles {
		pageBuf := bytes.NewBufferString("")
		fromTemplate(pageBuf, file.Name())

		fileErr := ioutil.WriteFile(filepath.Join(currentFolder, targetFolder, file.Name()), pageBuf.Bytes(), 0755)
		if fileErr != nil {
			logItf("can't create file: %v", fileErr)
		}

		logItf("%v -> ready \n", file.Name())
	}

	logIt("Completed.")
}

func logItf(f string, v ...interface{}) {
	logIt(fmt.Sprintf(f, v))
}

func logIt(msg string) {
	if isSilent {
		return
	}
	log.Print(msg)
}

func main() {
	isGenerate := flag.Bool("generate", false, "generate file, otherwise run local web server")
	port := flag.Int("port", 8080, "web port")
	folder := flag.String("folder", "", "source folder")
	flag.BoolVar(&isSilent, "s", false, "enable silent mode")

	flag.Parse()

	logItf("version: %s", version)

	if *folder == "" {
		var err error
		currentFolder, err = files.LocalFolder()
		if err != nil {
			log.Panicf("can't find local folder: %v", err)
		}
	} else {
		currentFolder = *folder
	}
	validateCurrentFolder(currentFolder)
	logItf("using folder: %v", currentFolder)

	prepareTemplates()

	if *isGenerate {
		generate()
		return
	}

	serve(":" + strconv.Itoa(*port))
}
