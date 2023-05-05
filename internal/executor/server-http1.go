package executor

import (
	"os"
	"path/filepath"
	"strings"
)

func ServerHttp(path string) {
	dep := getDependency(path)
	if dep == nil {
		return
	}

	lowerName := strings.ToLower(dep.Type)
	folder := makeSubPath(serverFolder, lowerName)

	if err := os.MkdirAll(folder, os.ModePerm); err != nil {
		panic(err)
	}

	buffer := strings.Builder{}
	buffer.WriteString("package " + lowerName + "\n\n")
	buffer.WriteString("import (\n")
	buffer.WriteString("\t\"github.com/julienschmidt/httprouter\"\n")
	buffer.WriteString("\t\"golang.org/x/net/http2\"\n")
	buffer.WriteString("\t\"github.com/quic-go/quic-go/http3\"\n")
	buffer.WriteString("\t\"net/http\"\n")
	buffer.WriteString("\t\"time\"\n")
	buffer.WriteString("\t\"log\"\n")
	buffer.WriteString("\t\"context\"\n")
	buffer.WriteString(")\n\n")
	buffer.WriteString("type " + dep.Type + " struct {\n")
	buffer.WriteString("\trouter *httprouter.Router\n")
	buffer.WriteString("\tserver *http.Server\n")
	buffer.WriteString("}\n\n")
	buffer.WriteString("func New" + dep.Type + "(addr string, readTimeOut time.Duration, writeTimeOut time.Duration) *" + dep.Type + " {\n")
	buffer.WriteString("\trouter := httprouter.New()\n")
	buffer.WriteString("\treturn &" + dep.Type + "{\n")
	buffer.WriteString("\t\trouter: router,\n")
	buffer.WriteString("\t\tserver: &http.Server{\n")
	buffer.WriteString("\t\t\tAddr:         addr,\n")
	buffer.WriteString("\t\t\tHandler:      router,\n")
	buffer.WriteString("\t\t\tReadTimeout:  readTimeOut,\n")
	buffer.WriteString("\t\t\tWriteTimeout: writeTimeOut,\n")
	buffer.WriteString("\t\t},\n")
	buffer.WriteString("\t}\n")
	buffer.WriteString("}\n\n")
	buffer.WriteString("func (s *" + dep.Type + ") Serve() {\n")
	buffer.WriteString("\tgo func() {\n")
	buffer.WriteString("\t\tif err := s.server.ListenAndServe(); err != nil {\n")
	buffer.WriteString("\t\t\tlog.Fatal(err)\n")
	buffer.WriteString("\t\t}\n")
	buffer.WriteString("\t}()\n")
	buffer.WriteString("}\n\n")
	buffer.WriteString("func (s *" + dep.Type + ") ServeTLS(certFile string, keyFile string) {\n")
	buffer.WriteString("\tgo func() {\n")
	buffer.WriteString("\t\tif err := s.server.ListenAndServeTLS(certFile, keyFile); err != nil {\n")
	buffer.WriteString("\t\t\tlog.Fatal(err)\n")
	buffer.WriteString("\t\t}\n")
	buffer.WriteString("\t}()\n")
	buffer.WriteString("}\n\n")
	buffer.WriteString("func (s *" + dep.Type + ") Stop() {\n")
	buffer.WriteString("\tif err := s.server.Shutdown(context.Background()); err != nil {\n")
	buffer.WriteString("\t\tlog.Fatal(err)\n")
	buffer.WriteString("\t}\n")
	buffer.WriteString("}\n\n")
	buffer.WriteString("func (s *" + dep.Type + ") Register(method string, path string, handler httprouter.Handle) {\n")
	buffer.WriteString("\ts.router.Handle(method, path, handler)\n")
	buffer.WriteString("}\n\n")
	buffer.WriteString("func (s *" + dep.Type + ") GET(path string, handler httprouter.Handle) {\n")
	buffer.WriteString("\ts.router.GET(path, handler)\n")
	buffer.WriteString("}\n\n")
	buffer.WriteString("func (s *" + dep.Type + ") POST(path string, handler httprouter.Handle) {\n")
	buffer.WriteString("\ts.router.POST(path, handler)\n")
	buffer.WriteString("}\n\n")
	buffer.WriteString("func (s *" + dep.Type + ") PUT(path string, handler httprouter.Handle) {\n")
	buffer.WriteString("\ts.router.PUT(path, handler)\n")
	buffer.WriteString("}\n\n")
	buffer.WriteString("func (s *" + dep.Type + ") DELETE(path string, handler httprouter.Handle) {\n")
	buffer.WriteString("\ts.router.DELETE(path, handler)\n")
	buffer.WriteString("}\n\n")
	buffer.WriteString("func (s *" + dep.Type + ") PATCH(path string, handler httprouter.Handle) {\n")
	buffer.WriteString("\ts.router.PATCH(path, handler)\n")
	buffer.WriteString("}\n\n")
	buffer.WriteString("func (s *" + dep.Type + ") OPTIONS(path string, handler httprouter.Handle) {\n")
	buffer.WriteString("\ts.router.OPTIONS(path, handler)\n")
	buffer.WriteString("}\n\n")
	buffer.WriteString("func (s *" + dep.Type + ") HEAD(path string, handler httprouter.Handle) {\n")
	buffer.WriteString("\ts.router.HEAD(path, handler)\n")
	buffer.WriteString("}\n\n")
	buffer.WriteString("func (s *" + dep.Type + ") ServeFiles(path string, root http.FileSystem) {\n")
	buffer.WriteString("\ts.router.ServeFiles(path, root)\n")
	buffer.WriteString("}\n\n")
	buffer.WriteString("func (s *" + dep.Type + ") ServeHTTP(w http.ResponseWriter, r *http.Request) {\n")
	buffer.WriteString("\ts.router.ServeHTTP(w, r)\n")
	buffer.WriteString("}\n\n")
	buffer.WriteString("func (s *" + dep.Type + ") UpgradeHTTP2() error {\n")
	buffer.WriteString("\treturn http2.ConfigureServer(s.server, nil)\n")
	buffer.WriteString("}\n\n")
	buffer.WriteString("func (s *" + dep.Type + ") ListenAndServe3(certFile string, keyFile string) error {\n")
	buffer.WriteString("\treturn http3.ListenAndServe(s.server.Addr, certFile, keyFile, s.router)\n")
	buffer.WriteString("}\n\n")
	buffer.WriteString("func (s *" + dep.Type + ") ListenAndServe3QUIC(certFile string, keyFile string) error {\n")
	buffer.WriteString("\treturn http3.ListenAndServeQUIC(s.server.Addr, certFile, keyFile, s.router)\n")
	buffer.WriteString("}\n\n")

	f, err := os.Create(filepath.Join(folder, lowerName+".go"))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}()

	if _, err := f.WriteString(buffer.String()); err != nil {
		panic(err)
	}

	if err := goGet("github.com/julienschmidt/httprouter"); err != nil {
		panic(err)
	}

	if err := goGet("golang.org/x/net/http2"); err != nil {
		panic(err)
	}

	if err := goGet("github.com/quic-go/quic-go"); err != nil {
		panic(err)
	}
}
