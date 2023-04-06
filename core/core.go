package core

import (
	"os"
	"os/exec"
)

func Run(arg ...string) {
	head := arg[0]
	parts := arg[1:len(arg)]
	run := exec.Command(head, parts...)
	run.Run()
}

func Variable(val string) string {
	if val == "home" {
		return string(img_home)
	} else if val == "ico" {
		return string(icon)
	} else {
		return string(borde)
	}
}

func AddFiles() {
	//--Eliminar la carpeta template
	Run("rm", "-rf", "./template")
	Run("mkdir", "./template")
	//--Crear todos los archivos
	file_index, _ := os.Create("./template/index.html")
	file_variante, _ := os.Create("./template/intersecciones.html")
	file_interseccion, _ := os.Create("./template/interseccion.html")
	file_edit, _ := os.Create("./template/edit_interseccion.html")
	file_equivalente, _ := os.Create("./template/analisis.html")
	file_real, _ := os.Create("./template/analisis_real.html")

	//--Agregar lineas a los archivos
	file_index.WriteString(index)
	file_variante.WriteString(variantes)
	file_interseccion.WriteString(ver_interseccion)
	file_edit.WriteString(edit_interseccion)
	file_equivalente.WriteString(analisis)
	file_real.WriteString(analisis_real)
}
