package executor

import "path/filepath"

const generated = "gen"
const pkg = "lib"
const internal = "internal"

var clientFolder = filepath.Join("lib", "client")
var workerFolder = filepath.Join("lib", "worker")
var serverFolder = filepath.Join("lib", "server")
var serviceFolder = filepath.Join("lib", "service")

var templateConfig = filepath.Join("template", "configs")
