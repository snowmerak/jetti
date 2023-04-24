package executor

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
)

func PprofHttp1(address string) {
	folder := filepath.Join(generated, "pprof", "http1")
	if err := os.MkdirAll(folder, os.ModePerm); err != nil {
		panic(err)
	}

	f, err := os.Create(filepath.Join(folder, "http1.go"))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}()

	buffer := bytes.NewBuffer(nil)
	buffer.WriteString("package http1\n\n")
	buffer.WriteString("import (\n")
	buffer.WriteString("\t\"net/http\"\n")
	buffer.WriteString("\t\"net/http/pprof\"\n")
	buffer.WriteString(")\n\n")
	buffer.WriteString("func ListenAndServe() error {\n")
	buffer.WriteString("\tmux := http.NewServeMux()\n")
	buffer.WriteString("\tmux.HandleFunc(\"/debug/pprof/\", pprof.Index)\n")
	buffer.WriteString("\tmux.HandleFunc(\"/debug/pprof/cmdline\", pprof.Cmdline)\n")
	buffer.WriteString("\tmux.HandleFunc(\"/debug/pprof/profile\", pprof.Profile)\n")
	buffer.WriteString("\tmux.HandleFunc(\"/debug/pprof/symbol\", pprof.Symbol)\n")
	buffer.WriteString("\tmux.HandleFunc(\"/debug/pprof/trace\", pprof.Trace)\n")
	buffer.WriteString("\treturn http.ListenAndServe(\"" + address + "\", mux)\n")
	buffer.WriteString("}\n\n")
	buffer.WriteString("func ListenAndServeTLS(certFile, keyFile string) error {\n")
	buffer.WriteString("\tmux := http.NewServeMux()\n")
	buffer.WriteString("\tmux.HandleFunc(\"/debug/pprof/\", pprof.Index)\n")
	buffer.WriteString("\tmux.HandleFunc(\"/debug/pprof/cmdline\", pprof.Cmdline)\n")
	buffer.WriteString("\tmux.HandleFunc(\"/debug/pprof/profile\", pprof.Profile)\n")
	buffer.WriteString("\tmux.HandleFunc(\"/debug/pprof/symbol\", pprof.Symbol)\n")
	buffer.WriteString("\tmux.HandleFunc(\"/debug/pprof/trace\", pprof.Trace)\n")
	buffer.WriteString("\treturn http.ListenAndServeTLS(\"" + address + "\", certFile, keyFile, mux)\n")
	buffer.WriteString("}\n\n")

	if _, err := f.Write(buffer.Bytes()); err != nil {
		panic(err)
	}
}

func PprofHttp2(addr string) {
	folder := filepath.Join(generated, "pprof", "http2")
	if err := os.MkdirAll(folder, os.ModePerm); err != nil {
		panic(err)
	}

	f, err := os.Create(filepath.Join(folder, "http2.go"))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}()

	buffer := bytes.NewBuffer(nil)
	buffer.WriteString("package http2\n\n")
	buffer.WriteString("import (\n")
	buffer.WriteString("\t\"net/http\"\n")
	buffer.WriteString("\t\"net/http/pprof\"\n")
	buffer.WriteString("\t\"golang.org/x/net/http2\"\n")
	buffer.WriteString(")\n\n")
	buffer.WriteString("func ListenAndServe() error {\n")
	buffer.WriteString("\tmux := http.NewServeMux()\n")
	buffer.WriteString("\tmux.HandleFunc(\"/debug/pprof/\", pprof.Index)\n")
	buffer.WriteString("\tmux.HandleFunc(\"/debug/pprof/cmdline\", pprof.Cmdline)\n")
	buffer.WriteString("\tmux.HandleFunc(\"/debug/pprof/profile\", pprof.Profile)\n")
	buffer.WriteString("\tmux.HandleFunc(\"/debug/pprof/symbol\", pprof.Symbol)\n")
	buffer.WriteString("\tmux.HandleFunc(\"/debug/pprof/trace\", pprof.Trace)\n")
	buffer.WriteString("\tserver := http.Server{\n")
	buffer.WriteString("\t\tAddr: \"" + addr + "\",\n")
	buffer.WriteString("\t\tHandler: mux,\n")
	buffer.WriteString("\t}\n")
	buffer.WriteString("\tif err := http2.ConfigureServer(&server, nil); err != nil {\n")
	buffer.WriteString("\t\treturn err\n")
	buffer.WriteString("\t}\n")
	buffer.WriteString("\treturn server.ListenAndServe()\n")
	buffer.WriteString("}\n\n")
	buffer.WriteString("func ListenAndServeTLS(certFile, keyFile string) error {\n")
	buffer.WriteString("\tmux := http.NewServeMux()\n")
	buffer.WriteString("\tmux.HandleFunc(\"/debug/pprof/\", pprof.Index)\n")
	buffer.WriteString("\tmux.HandleFunc(\"/debug/pprof/cmdline\", pprof.Cmdline)\n")
	buffer.WriteString("\tmux.HandleFunc(\"/debug/pprof/profile\", pprof.Profile)\n")
	buffer.WriteString("\tmux.HandleFunc(\"/debug/pprof/symbol\", pprof.Symbol)\n")
	buffer.WriteString("\tmux.HandleFunc(\"/debug/pprof/trace\", pprof.Trace)\n")
	buffer.WriteString("\tserver := http.Server{\n")
	buffer.WriteString("\t\tAddr: \"" + addr + "\",\n")
	buffer.WriteString("\t\tHandler: mux,\n")
	buffer.WriteString("\t}\n")
	buffer.WriteString("\tif err := http2.ConfigureServer(&server, nil); err != nil {\n")
	buffer.WriteString("\t\treturn err\n")
	buffer.WriteString("\t}\n")
	buffer.WriteString("\treturn server.ListenAndServeTLS(certFile, keyFile)\n")
	buffer.WriteString("}\n\n")

	if _, err := f.Write(buffer.Bytes()); err != nil {
		panic(err)
	}

	cmd := exec.Command("go", "get", "golang.org/x/net/http2")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		panic(err)
	}
}

func PprofHttp3(addr string) {
	folder := filepath.Join(generated, "pprof", "http3")
	if err := os.MkdirAll(folder, os.ModePerm); err != nil {
		panic(err)
	}

	f, err := os.Create(filepath.Join(folder, "http3.go"))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}()

	buffer := bytes.NewBuffer(nil)
	buffer.WriteString("package http3\n\n")
	buffer.WriteString("import (\n")
	buffer.WriteString("\t\"net/http\"\n")
	buffer.WriteString("\t\"net/http/pprof\"\n")
	buffer.WriteString("\t\"github.com/quic-go/quic-go/http3\"\n")
	buffer.WriteString(")\n\n")
	buffer.WriteString("func ListenAndServeTLS(certFile, keyFile string) error {\n")
	buffer.WriteString("\tmux := http.NewServeMux()\n")
	buffer.WriteString("\tmux.HandleFunc(\"/debug/pprof/\", pprof.Index)\n")
	buffer.WriteString("\tmux.HandleFunc(\"/debug/pprof/cmdline\", pprof.Cmdline)\n")
	buffer.WriteString("\tmux.HandleFunc(\"/debug/pprof/profile\", pprof.Profile)\n")
	buffer.WriteString("\tmux.HandleFunc(\"/debug/pprof/symbol\", pprof.Symbol)\n")
	buffer.WriteString("\tmux.HandleFunc(\"/debug/pprof/trace\", pprof.Trace)\n")
	buffer.WriteString("\treturn http3.ListenAndServeQUIC(\"")
	buffer.WriteString(addr)
	buffer.WriteString("\", certFile, keyFile, mux)\n")
	buffer.WriteString("}\n\n")

	if _, err := f.Write(buffer.Bytes()); err != nil {
		panic(err)
	}

	cmd := exec.Command("go", "get", "github.com/quic-go/quic-go")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		panic(err)
	}
}
