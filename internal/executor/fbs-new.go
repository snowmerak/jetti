package executor

import (
	"fmt"
	"os"
	"path/filepath"
)

var fbsFolder = filepath.Join("template", "fbs")

func FbsNew(path string) {
	// packageName := "./" + filepath.Dir(path)
	path = filepath.Join(fbsFolder, path)
	dir := filepath.Dir(path)
	namespace := filepath.Base(dir)
	// base := filepath.Base(path)
	// ext := filepath.Ext(base)
	// name := base[:len(base)-len(ext)]

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		panic(err)
	}

	if _, err := os.Stat(path); !os.IsNotExist(err) {
		panic("file already exists")
	}

	f, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}()

	const template = `namespace %s;

enum Color:byte { Red = 0, Green, Blue = 2 }
 
union Equipment { Weapon } // Optionally add more tables.
 
struct Vec3 {
  x:float;
  y:float;
  z:float;
}
 
table Monster {
  pos:Vec3; // Struct.
  mana:short = 150;
  hp:short = 100;
  name:string;
  friendly:bool = false (deprecated);
  inventory:[ubyte];  // Vector of scalars.
  color:Color = Blue; // Enum.
  weapons:[Weapon];   // Vector of tables.
  equipped:Equipment; // Union.
  path:[Vec3];        // Vector of structs.
}
 
table Weapon {
  name:string;
  damage:short;
}
 
root_type Monster;
`

	if _, err := f.Write([]byte(fmt.Sprintf(template, namespace))); err != nil {
		panic(err)
	}
}
