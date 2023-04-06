package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"net/http"
	"strconv"
	"strings"
	//"runtime"
	"fmt"
	"github.com/akosej/saiu/core"
	"math"
)

var db *gorm.DB
var err error

//-----------------------------------------------------------------
//---------------ESTRUCTURAS DE LA BASE DE DATOS ------------------
//-----------------------------------------------------------------

type Analisis struct {
	ID          uint   `json:"id" form:"id"`
	Nombre      string `json:"nombre" form:"nombre" binding:"required"`
	Descripcion string `json:"descripcion" form:"descripcion" binding:"required"`
	Autor       string `json:"autor" form:"autor" binding:"required"`
}

type Intersecciones struct {
	ID          uint   `json:"id" form:"id"`
	Analisis    uint   `json:"analisis" form:"analisis"`
	Nombre      string `json:"nombre" form:"nombre" binding:"required"`
	Descripcion string `json:"descripcion" form:"descripcion" binding:"required"`
	Tipo        string `json:"tipo" form:"tipo" binding:"required"`
	Area        string `json:"area" form:"area" binding:"required"`
	Fecha       string `json:"fecha" form:"fecha" binding:"required"`
}

type Accesos struct {
	ID                       uint   `json:"id" form:"id"`
	Interseccion             int    `json:"interseccion" form:"interseccion" binding:"required"`
	Nombre                   string `json:"nombre" form:"nombre" binding:"required"`
	Pendiente                string `json:"pendiente" form:"pendiente" binding:"required"`
	Fhmd                     string `json:"fhmd" form:"fhmd" binding:"required"`
	Intensidad               string `json:"intensidad" form:"intensidad" binding:"required"`
	Maniobrasestacionamiento string `json:"maniobrasestacionamiento" form:"maniobrasestacionamiento" binding:"required"`
	Maniobrasparada          string `json:"maniobrasparada" form:"maniobrasparada" binding:"required"`
	Tpad                     string `json:"tpad" form:"tpad" binding:"required"`
	Dossentidos              string `json:"dossentidos" form:"dossentidos" binding:"required"`
	Parada                   string `json:"parada" form:"parada" binding:"required"`
	Isleta                   string `json:"isleta" form:"isleta" binding:"required"`
	Separadorcentral         string `json:"separadorcentral" form:"separadorcentral" binding:"required"`
	Carrilseparador          string `json:"carrilseparador" form:"carrilseparador" binding:"required"`
	Cebra                    string `json:"cebra" form:"cebra" binding:"required"`
}

type Carriles struct {
	ID     uint   `json:"id" form:"id"`
	Acceso int    `json:"acceso" form:"acceso" binding:"required"`
	Ancho  string `json:"ancho" form:"ancho" binding:"required"`
}

type Movimientos struct {
	ID             uint   `json:"id" form:"id"`
	Carril         int    `json:"carril" form:"carril" binding:"required"`
	Tipo           string `json:"tipo" form:"tipo" binding:"required"`
	Autos          int    `json:"autos" form:"autos"`
	Motos          int    `json:"motos" form:"motos"`
	Ciclos         int    `json:"ciclos" form:"ciclos"`
	Coches         int    `json:"coches" form:"coches"`
	CamionesCarga  int    `json:"camionescarga" form:"camionescarga"`
	CamionesPasaje int    `json:"camionespasaje" form:"camionespasaje"`
	Omnibus        int    `json:"omnibus" form:"omnibus"`
}

type Fase_semaforo struct {
	ID           uint   `json:"id" form:"id"`
	Interseccion int    `json:"interseccion" form:"interseccion" binding:"required"`
	Nombre       string `json:"nombre" form:"nombre" binding:"required"`
	Verde        int    `json:"verde" form:"verde" binding:"required"`
	Amarillo     int    `json:"amarillo" form:"amarillo" binding:"required"`
	Rojo         int    `json:"rojo" form:"rojo" binding:"required"`
	Todo_rojo    int    `json:"todo_rojo" form:"todo_rojo" binding:"required"`
	Ciclo        int    `json:"ciclo" form:"ciclo" binding:"required"`
}

type Grupo_carril struct {
	ID       uint   `json:"id" form:"id"`
	Acceso   int    `json:"acceso" form:"acceso" binding:"required"`
	Fase     int    `json:"fase" form:"fase" binding:"required"`
	Carriles string `json:"carriles" form:"carriles" binding:"required"`
}

// -------------------------------------------------------------------
// -----Estructura para exportar intersecciones ----------------------
// -------------------------------------------------------------------
type Grupo_exp struct {
	Acceso   string
	Fase     int
	Carriles string
}

type Fase_exp struct {
	ID           uint
	Interseccion int
	Nombre       string
	Verde        int
	Amarillo     int
	Rojo         int
	Todo_rojo    int
	Ciclo        int
}

type Movimientos_exp struct {
	Carril         uint
	Tipo           string
	Autos          int
	Motos          int
	Ciclos         int
	Coches         int
	CamionesCarga  int
	CamionesPasaje int
	Omnibus        int
}

type Carriles_exp struct {
	ID     uint
	Acceso string
	Ancho  string
}

type Accesos_exp struct {
	Nombre                   string
	Pendiente                string
	Fhmd                     string
	Intensidad               string
	Maniobrasestacionamiento string
	Maniobrasparada          string
	Tpad                     string
	Dossentidos              string
	Parada                   string
	Isleta                   string
	Separadorcentral         string
	Carrilseparador          string
	Cebra                    string
}

type JSON struct {
	Analsis     uint
	Nombre      string
	Descripcion string
	Tipo        string
	Area        string
	Fecha       string
	Accesos     []Accesos_exp
	Carriles    []Carriles_exp
	Movimientos []Movimientos_exp
	Fases       []Fase_exp
	Grupos      []Grupo_exp
}

// -----------------------------------------------------------------
// ---------------FUNCION PRINCIPAL --------------------------------
// -----------------------------------------------------------------
func main() {
	//------------------------------------------
	//-------SERVICION DE SQLITE ---------------
	//------------------------------------------
	db, err = gorm.Open("sqlite3", "./data/data.db")

	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()
	db.AutoMigrate(&Analisis{}, &Intersecciones{}, &Accesos{}, &Carriles{}, Movimientos{}, &Fase_semaforo{}, &Grupo_carril{})
	//------------------------------------------
	//-------Crear ficheros necesarios ---------
	//------------------------------------------
	//os := runtime.GOOS
	//switch  os {
	//case "linux":
	//	core.AddFiles()
	//default:
	//	fmt.Printf(" %s . ", os)
	//}
	//------------------------------------------
	//-------ELEMENTOS DEL SERVIDOR-------------
	//------------------------------------------
	//gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	//-- Cargar plantillas
	r.LoadHTMLGlob("template/*")
	//-- Cargar archivos estaticos
	r.StaticFS("/static/", http.Dir("./static/"))
	//------------------------------------------
	//-------TODAS LAS PETICIONES GET ----------
	//------------------------------------------
	r.GET("/", func(c *gin.Context) {
		var analisis []Analisis
		db.Find(&analisis)
		c.HTML(http.StatusOK, "index.html", gin.H{"values": analisis, "ico": core.Variable("ico"), "img_home": core.Variable("home"), "img_borde": core.Variable("borde")})
	})
	r.GET("/intersecciones/:id", func(c *gin.Context) {
		id := c.Params.ByName("id")
		var analisis Analisis
		db.Where("ID = ?", id).First(&analisis)
		var intersecciones []Intersecciones
		db.Where("Analisis = ?", id).Find(&intersecciones)
		c.HTML(http.StatusOK, "intersecciones.html", gin.H{"analisis": analisis, "intersecciones": intersecciones, "ico": core.Variable("ico"), "img_borde": core.Variable("borde")})
	})
	r.GET("/interseccion/:id", func(c *gin.Context) {
		id := c.Params.ByName("id")
		var interseccion Intersecciones
		db.Where("ID = ?", id).First(&interseccion)
		var acceso []Accesos
		db.Where("Interseccion = ?", id).Find(&acceso)
		var fases []Fase_semaforo
		db.Where("Interseccion = ?", id).Find(&fases)

		c.HTML(http.StatusOK, "interseccion.html", gin.H{"interseccion": interseccion, "accesos": acceso, "fases": fases, "ico": core.Variable("ico")})
	})
	r.GET("/edit_interseccion/:id", func(c *gin.Context) {
		id := c.Params.ByName("id")
		var interseccion Intersecciones
		db.Where("ID = ?", id).First(&interseccion)
		var acceso []Accesos
		db.Where("Interseccion = ?", id).Find(&acceso)
		var fases []Fase_semaforo
		db.Where("Interseccion = ?", id).Find(&fases)

		c.HTML(http.StatusOK, "edit_interseccion.html", gin.H{"interseccion": interseccion, "accesos": acceso, "fases": fases, "ico": core.Variable("ico"), "img_borde": core.Variable("borde")})
	})
	r.GET("/analisis/:id", func(c *gin.Context) {
		id := c.Params.ByName("id")
		var interseccion Intersecciones
		db.Where("ID = ?", id).First(&interseccion)
		var acceso []Accesos
		db.Where("Interseccion = ?", id).Find(&acceso)
		var fases []Fase_semaforo
		db.Where("Interseccion = ?", id).Find(&fases)

		c.HTML(http.StatusOK, "analisis.html", gin.H{"interseccion": interseccion, "accesos": acceso, "fases": fases, "ico": core.Variable("ico")})
	})
	r.GET("/analisis_real/:id", func(c *gin.Context) {
		id := c.Params.ByName("id")
		var interseccion Intersecciones
		db.Where("ID = ?", id).First(&interseccion)
		var acceso []Accesos
		db.Where("Interseccion = ?", id).Find(&acceso)
		var fases []Fase_semaforo
		db.Where("Interseccion = ?", id).Find(&fases)

		c.HTML(http.StatusOK, "analisis_real.html", gin.H{"interseccion": interseccion, "accesos": acceso, "fases": fases, "ico": core.Variable("ico")})
	})
	r.GET("/exportar/:id", func(c *gin.Context) {
		id := c.Params.ByName("id")
		var j JSON
		var interseccion Intersecciones
		db.Where("ID = ?", id).First(&interseccion)

		//-Nombre de la interseccion
		j.Analsis = interseccion.Analisis
		j.Nombre = interseccion.Nombre
		j.Descripcion = interseccion.Descripcion
		j.Tipo = interseccion.Tipo
		j.Area = interseccion.Area
		j.Fecha = interseccion.Fecha
		var fases []Fase_semaforo
		db.Where("Interseccion = ?", interseccion.ID).Find(&fases)
		for _, fase := range fases {
			j.Fases = append(j.Fases, Fase_exp{ID: fase.ID, Nombre: fase.Nombre, Verde: fase.Verde, Amarillo: fase.Amarillo, Rojo: fase.Rojo, Todo_rojo: fase.Todo_rojo, Ciclo: fase.Ciclo})
		}

		var accesos []Accesos
		db.Where("Interseccion = ?", interseccion.ID).Find(&accesos)
		for _, acceso := range accesos {
			var carriles []Carriles
			db.Where("Acceso = ?", acceso.ID).Find(&carriles)
			j.Accesos = append(j.Accesos, Accesos_exp{Nombre: acceso.Nombre, Pendiente: acceso.Pendiente, Fhmd: acceso.Fhmd, Intensidad: acceso.Intensidad, Maniobrasestacionamiento: acceso.Maniobrasestacionamiento, Maniobrasparada: acceso.Maniobrasparada, Tpad: acceso.Tpad, Dossentidos: acceso.Dossentidos, Parada: acceso.Parada, Isleta: acceso.Isleta, Separadorcentral: acceso.Separadorcentral, Carrilseparador: acceso.Carrilseparador, Cebra: acceso.Cebra})
			for _, carril := range carriles {
				j.Carriles = append(j.Carriles, Carriles_exp{ID: carril.ID, Acceso: acceso.Nombre, Ancho: carril.Ancho})
				var movimientos []Movimientos
				db.Where("Carril = ?", carril.ID).Find(&movimientos)
				for _, movimiento := range movimientos {
					j.Movimientos = append(j.Movimientos, Movimientos_exp{Carril: carril.ID, Tipo: movimiento.Tipo, Autos: movimiento.Autos, Motos: movimiento.Motos, Ciclos: movimiento.Ciclos, Coches: movimiento.Coches, CamionesPasaje: movimiento.CamionesPasaje, CamionesCarga: movimiento.CamionesCarga, Omnibus: movimiento.Omnibus})
				}
			}

			var grupos []Grupo_carril
			db.Where("Acceso = ?", acceso.ID).Find(&grupos)
			for _, grupo := range grupos {
				j.Grupos = append(j.Grupos, Grupo_exp{Acceso: acceso.Nombre, Fase: grupo.Fase, Carriles: grupo.Carriles})
			}
		}

		c.JSON(200, j)
	})
	//------------------------------------------
	//-------TODAS LAS PETICIONES POST----------
	//------------------------------------------
	r.POST("/", Addanalisis)
	r.POST("/intersecciones", Addinterseccion)
	r.POST("/add_acceso", Addaccesos)
	r.POST("/add_carril", Addcarriles)
	r.POST("/add_movimiento", Addmovimientos)
	r.POST("/add_fase", Addfases)
	r.POST("/add_grupo", AddGrupo)
	r.POST("/update_acceso", UpdateAccesos)
	r.POST("/update_movimiento", UpdateMovimientos)
	r.POST("/update_fase", UpdateFases)
	r.POST("/update_grupo", UpdateGrupo)
	r.POST("/update_carril", UpdateCarril)
	r.POST("/del_acceso", DeleteAcceso)
	r.POST("/del_carril", DeleteCarril)
	r.POST("/del_movimiento", DeleteMovimiento)
	r.POST("/del_fase", DeleteFase)
	r.POST("/del_grupo", DeleteGrupo)
	r.POST("/importar", Importar)

	//------------------------------------------
	//----------INICIAR EL SERVIDOR-------------
	//------------------------------------------
	r.Run(":8081")
}

// -----------------------------------------------------------------
// ---------------Importar--------------------------------
// -----------------------------------------------------------------
func Importar(c *gin.Context) {
	//var json JSON
	text := c.PostForm("json")
	var j JSON
	json.Unmarshal([]byte(text), &j)

	var interseccion Intersecciones
	interseccion.Analisis = uint(j.Analsis)
	interseccion.Nombre = j.Nombre
	interseccion.Descripcion = j.Descripcion
	interseccion.Tipo = j.Tipo
	interseccion.Area = j.Area
	interseccion.Fecha = j.Fecha
	db.Create(&interseccion)
	//--Importar fases del semaforo
	for _, j_fase := range j.Fases {
		var fase Fase_semaforo
		fase.Interseccion = int(interseccion.ID)
		fase.Nombre = j_fase.Nombre
		fase.Verde = j_fase.Verde
		fase.Amarillo = j_fase.Amarillo
		fase.Rojo = j_fase.Rojo
		fase.Todo_rojo = j_fase.Todo_rojo
		fase.Ciclo = j_fase.Ciclo
		db.Create(&fase)
	}
	//--Importar accesos
	for _, j_acceso := range j.Accesos {
		var acceso Accesos
		acceso.Interseccion = int(interseccion.ID)
		acceso.Nombre = j_acceso.Nombre
		acceso.Pendiente = j_acceso.Pendiente
		acceso.Fhmd = j_acceso.Fhmd
		acceso.Intensidad = j_acceso.Intensidad
		acceso.Maniobrasestacionamiento = j_acceso.Maniobrasestacionamiento
		acceso.Maniobrasparada = j_acceso.Maniobrasparada
		acceso.Tpad = j_acceso.Tpad
		acceso.Dossentidos = j_acceso.Dossentidos
		acceso.Parada = j_acceso.Parada
		acceso.Isleta = j_acceso.Isleta
		acceso.Separadorcentral = j_acceso.Separadorcentral
		acceso.Carrilseparador = j_acceso.Carrilseparador
		acceso.Cebra = j_acceso.Cebra
		db.Create(&acceso)
		//--Importar carriles
		for _, j_carril := range j.Carriles {
			if j_carril.Acceso == j_acceso.Nombre {
				var carril Carriles
				carril.Acceso = int(acceso.ID)
				carril.Ancho = j_carril.Ancho
				db.Create(&carril)
				//--Importar movimientos
				for _, j_mov := range j.Movimientos {
					if j_mov.Carril == j_carril.ID {
						var movimiento Movimientos
						movimiento.Carril = int(carril.ID)
						movimiento.Tipo = j_mov.Tipo
						movimiento.Autos = j_mov.Autos
						movimiento.Motos = j_mov.Motos
						movimiento.Ciclos = j_mov.Ciclos
						movimiento.Coches = j_mov.Coches
						movimiento.CamionesCarga = j_mov.CamionesCarga
						movimiento.CamionesPasaje = j_mov.CamionesPasaje
						movimiento.Omnibus = j_mov.Omnibus
						db.Create(&movimiento)
					}
				}
				//--Importar grupos de carril
				//for _,j_grupo:=range j.Grupos{
				//	var text_carril string
				//	carriles := strings.Split(j_grupo.Carriles, ",")
				//	for i := 0; i < len(carriles)-1; i++ {
				//		if carriles[i] == string(strconv.Itoa(int(j_carril.ID))) {
				//			var grupo Grupo_carril
				//			db.Where("ID = ?", j_grupo.).First(&grupo)
				//text_carril +=
				//}
				//}
				//}
			}
		}
	}

	c.JSON(200, gin.H{"msg": interseccion.ID})

}

// -----------------------------------------------------------------
// ---------------AGREGAR ELEMENTOS --------------------------------
// -----------------------------------------------------------------
// --Funcion para agragar los analisis
func Addanalisis(c *gin.Context) {
	var analisis Analisis
	analisis.Nombre = c.PostForm("nombre")
	analisis.Descripcion = c.PostForm("descripcion")
	analisis.Autor = c.PostForm("autor")
	if err = c.MustBindWith(&analisis, binding.Form); err == nil {
		db.Create(&analisis)
		c.JSON(200, gin.H{"msg": "Analisis creado"})
		return
	}
	c.JSON(400, gin.H{"msg": "datos no validos"})
}

// --Funcion para agragar las intersecciones
func Addinterseccion(c *gin.Context) {
	var interseccion Intersecciones
	interseccion.Nombre = c.PostForm("nombre")
	interseccion.Descripcion = c.PostForm("descripcion")
	interseccion.Tipo = c.PostForm("tipo")
	interseccion.Area = c.PostForm("area")
	interseccion.Fecha = c.PostForm("fecha")
	if err = c.MustBindWith(&interseccion, binding.Form); err == nil {
		db.Create(&interseccion)
		c.JSON(200, gin.H{"msg": "Inteseccion creada"})
		return
	}
	c.JSON(400, gin.H{"msg": "datos no validos"})
}

// --Funcion para agragar las accesos
func Addaccesos(c *gin.Context) {
	var acceso Accesos
	//--Convertir el strin en int
	id_interseccion, err := strconv.Atoi(c.PostForm("interseccion"))
	//---------------------------------------
	acceso.Interseccion = id_interseccion
	acceso.Nombre = c.PostForm("nombre")
	acceso.Pendiente = c.PostForm("pendiente")
	acceso.Fhmd = c.PostForm("fhmd")
	acceso.Intensidad = c.PostForm("intensidad")
	acceso.Maniobrasestacionamiento = c.PostForm("maniobrasestacionamiento")
	acceso.Maniobrasparada = c.PostForm("Maniobrasparada")
	acceso.Tpad = c.PostForm("tpad")
	acceso.Dossentidos = c.PostForm("Dossentidos")
	acceso.Parada = c.PostForm("parada")
	acceso.Isleta = c.PostForm("isleta")
	acceso.Separadorcentral = c.PostForm("separadorcentral")
	acceso.Carrilseparador = c.PostForm("carrilseparador")
	acceso.Cebra = c.PostForm("cebra")
	//---------------------------------------
	if err = c.MustBindWith(&acceso, binding.Form); err == nil {
		db.Create(&acceso)
		c.JSON(200, gin.H{"msg": "Acceso creado"})
		return
	}
	c.JSON(400, gin.H{"msg": "datos no validos"})

}

// --Funcion para agragar las carriles
func Addcarriles(c *gin.Context) {
	var carril Carriles
	//--Convertir el strin en int
	id_acceso, err := strconv.Atoi(c.PostForm("acceso"))
	//---------------------------------------
	carril.Acceso = id_acceso
	carril.Ancho = c.PostForm("ancho")
	//---------------------------------------
	if err = c.MustBindWith(&carril, binding.Form); err == nil {
		db.Create(&carril)
		c.JSON(200, gin.H{"msg": "Carril creado"})
		return
	}
	c.JSON(400, gin.H{"msg": "datos no validos"})

}

// --Funcion para agragar los moviminetos
func Addmovimientos(c *gin.Context) {

	var movimiento Movimientos
	//--Convertir el strin en int
	id_carril, _ := strconv.Atoi(c.PostForm("carril"))
	autos, _ := strconv.Atoi(c.PostForm("autos"))
	motos, _ := strconv.Atoi(c.PostForm("motos"))
	ciclos, _ := strconv.Atoi(c.PostForm("ciclos"))
	coches, _ := strconv.Atoi(c.PostForm("coches"))
	camionescarga, _ := strconv.Atoi(c.PostForm("camionescarga"))
	camionespasaje, _ := strconv.Atoi(c.PostForm("camionespasaje"))
	omnibus, _ := strconv.Atoi(c.PostForm("omnibus"))
	//---------------------------------------
	movimiento.Carril = id_carril
	movimiento.Tipo = c.PostForm("tipo")
	movimiento.Autos = autos
	movimiento.Motos = motos
	movimiento.Coches = coches
	movimiento.CamionesCarga = camionescarga
	movimiento.Ciclos = ciclos
	movimiento.CamionesPasaje = camionespasaje
	movimiento.Omnibus = omnibus

	//---------------------------------------
	if err = c.MustBindWith(&movimiento, binding.Form); err == nil {
		db.Create(&movimiento)
		c.JSON(200, gin.H{"msg": "Movimiento creado"})
		return
	}

	c.JSON(200, gin.H{"msg": movimiento})

}

// --Funcion para agragar las fases del semaforo
func Addfases(c *gin.Context) {

	var fase Fase_semaforo
	//--Convertir el strin en int
	interseccion, _ := strconv.Atoi(c.PostForm("interseccion"))
	verde, _ := strconv.Atoi(c.PostForm("verde"))
	amarillo, _ := strconv.Atoi(c.PostForm("amarillo"))
	rojo, _ := strconv.Atoi(c.PostForm("rojo"))
	todorojo, _ := strconv.Atoi(c.PostForm("todorojo"))
	//---------------------------------------
	fase.Interseccion = interseccion
	fase.Nombre = c.PostForm("nombre")
	fase.Verde = verde
	fase.Amarillo = amarillo
	fase.Rojo = rojo
	fase.Todo_rojo = todorojo
	fase.Ciclo = verde + amarillo + rojo

	//---------------------------------------
	if err = c.MustBindWith(&fase, binding.Form); err == nil {
		db.Create(&fase)
		c.JSON(200, gin.H{"msg": "Fase creada"})
		return
	}

	c.JSON(200, gin.H{"msg": "Error en los datos"})

}

// --Funcion para agragar los grupos
func AddGrupo(c *gin.Context) {

	var grupo Grupo_carril
	//--Convertir el strin en int
	acceso, _ := strconv.Atoi(c.PostForm("acceso"))
	fase, _ := strconv.Atoi(c.PostForm("fase"))

	//---------------------------------------
	grupo.Acceso = acceso
	grupo.Fase = fase
	grupo.Carriles = c.PostForm("carriles")

	//---------------------------------------
	if err = c.MustBindWith(&grupo, binding.Form); err == nil {
		db.Create(&grupo)
		c.JSON(200, gin.H{"msg": "Grupo de carril creado"})
		return
	}

	c.JSON(200, gin.H{"msg": "Error en los datos"})

}

//-----------------------------------------------------------------
//---------------ACTUALIZAR ELEMENTOS -----------------------------
//-----------------------------------------------------------------

// --Funcion para actualizar los accesos
func UpdateAccesos(c *gin.Context) {
	var acceso Accesos
	id := c.PostForm("id")

	if err := db.Where("id = ?", id).First(&acceso).Error; err != nil {
		c.AbortWithStatus(404)
	}
	acceso.Pendiente = c.PostForm("pendiente")
	acceso.Fhmd = c.PostForm("fhmd")
	acceso.Intensidad = c.PostForm("intensidad")
	acceso.Maniobrasestacionamiento = c.PostForm("maniobrasestacionamiento")
	acceso.Maniobrasparada = c.PostForm("maniobrasparada")
	acceso.Tpad = c.PostForm("tpad")
	acceso.Dossentidos = c.PostForm("dossentidos")
	acceso.Parada = c.PostForm("parada")
	acceso.Isleta = c.PostForm("isleta")
	acceso.Separadorcentral = c.PostForm("separadorcentral")
	acceso.Carrilseparador = c.PostForm("carrilseparador")
	acceso.Cebra = c.PostForm("cebra")
	db.Save(&acceso)
	c.JSON(200, gin.H{"msg": "Acceso update"})
	return
}

// --Funcion para actualizar los moviminetos
func UpdateMovimientos(c *gin.Context) {
	var movimiento Movimientos
	id := c.PostForm("id")

	if err := db.Where("id = ?", id).First(&movimiento).Error; err != nil {
		c.AbortWithStatus(404)
	}

	//--Convertir el strin en int
	id_carril, _ := strconv.Atoi(c.PostForm("carril"))
	autos, _ := strconv.Atoi(c.PostForm("autos"))
	motos, _ := strconv.Atoi(c.PostForm("motos"))
	ciclos, _ := strconv.Atoi(c.PostForm("ciclos"))
	coches, _ := strconv.Atoi(c.PostForm("coches"))
	camionescarga, _ := strconv.Atoi(c.PostForm("camionescarga"))
	camionespasaje, _ := strconv.Atoi(c.PostForm("camionespasaje"))
	omnibus, _ := strconv.Atoi(c.PostForm("omnibus"))
	//---------------------------------------
	movimiento.Carril = id_carril
	movimiento.Tipo = c.PostForm("tipo")
	movimiento.Autos = autos
	movimiento.Motos = motos
	movimiento.Coches = coches
	movimiento.CamionesCarga = camionescarga
	movimiento.Ciclos = ciclos
	movimiento.CamionesPasaje = camionespasaje
	movimiento.Omnibus = omnibus

	//---------------------------------------

	db.Save(&movimiento)
	c.JSON(200, gin.H{"msg": "Movimiento creado"})
	return

}

// --Funcion para actualizar las fases del semaforo
func UpdateFases(c *gin.Context) {

	var fase Fase_semaforo
	id := c.PostForm("id")
	if err := db.Where("id = ?", id).First(&fase).Error; err != nil {
		c.AbortWithStatus(404)
	}
	//--Convertir el strin en int
	interseccion, _ := strconv.Atoi(c.PostForm("interseccion"))
	verde, _ := strconv.Atoi(c.PostForm("verde"))
	amarillo, _ := strconv.Atoi(c.PostForm("amarillo"))
	rojo, _ := strconv.Atoi(c.PostForm("rojo"))
	todorojo, _ := strconv.Atoi(c.PostForm("todorojo"))
	//---------------------------------------
	fase.Interseccion = interseccion
	fase.Nombre = c.PostForm("nombre")
	fase.Verde = verde
	fase.Amarillo = amarillo
	fase.Rojo = rojo
	fase.Todo_rojo = todorojo
	fase.Ciclo = verde + amarillo + rojo

	//---------------------------------------
	if err = c.MustBindWith(&fase, binding.Form); err == nil {
		db.Save(&fase)
		c.JSON(200, gin.H{"msg": "Fase editada"})
		return
	}

	c.JSON(200, gin.H{"msg": "Error en los datos"})

}

// --Funcion para actualizar los grupos
func UpdateGrupo(c *gin.Context) {

	var grupo Grupo_carril
	//--Convertir el strin en int
	id := c.PostForm("id")

	if err := db.Where("id = ?", id).First(&grupo).Error; err != nil {
		c.AbortWithStatus(404)
	}

	acceso, _ := strconv.Atoi(c.PostForm("acceso"))
	fase, _ := strconv.Atoi(c.PostForm("fase"))

	//---------------------------------------
	grupo.Acceso = acceso
	grupo.Fase = fase
	grupo.Carriles = c.PostForm("carriles")

	//---------------------------------------
	if err = c.MustBindWith(&grupo, binding.Form); err == nil {
		db.Save(&grupo)
		c.JSON(200, gin.H{"msg": "Grupo de carril actualizado"})
		return
	}

	c.JSON(200, gin.H{"msg": "Error en los datos"})

}

// --Funcion para actualizar los carriles
func UpdateCarril(c *gin.Context) {

	var carril Carriles
	//--Convertir el strin en int
	id := c.PostForm("id")

	if err := db.Where("id = ?", id).First(&carril).Error; err != nil {
		c.AbortWithStatus(404)
	}

	acceso, _ := strconv.Atoi(c.PostForm("acceso"))
	//---------------------------------------
	carril.Acceso = acceso
	carril.Ancho = c.PostForm("ancho")

	//---------------------------------------
	if err = c.MustBindWith(&carril, binding.Form); err == nil {
		db.Save(&carril)
		c.JSON(200, gin.H{"msg": "Carril actualizado"})
		return
	}

	c.JSON(200, gin.H{"msg": "Error en los datos"})

}

//-----------------------------------------------------------------
//---------------ELIMINAR ELEMENTOS -------------------------------
//-----------------------------------------------------------------

// --Funcion para eliminar los accesos
func DeleteAcceso(c *gin.Context) {
	id := c.PostForm("id")

	var carril []Carriles
	db.Where("acceso = ?", id).Find(&carril)
	var carr Carriles

	var movimiento []Movimientos
	var movs Movimientos

	for _, car := range carril {
		db.Where("carril = ?", car.ID).Find(&movimiento)
		for _, mov := range movimiento {
			db.Where("id = ?", mov.ID).Delete(&movs)
		}
		db.Where("id = ?", car.ID).Delete(&carr)
	}

	var acceso Accesos
	db.Where("id = ?", id).Delete(&acceso)
	c.JSON(200, gin.H{"msg": "Acceso eliminado"})
}

// --Funcion para eliminar los carriles
func DeleteCarril(c *gin.Context) {
	id := c.PostForm("id")
	//--Eliminar los movimientos del carril
	var movimiento []Movimientos
	db.Where("carril = ?", id).Find(&movimiento)
	var movs Movimientos
	for _, mov := range movimiento {
		db.Where("id = ?", mov.ID).Delete(&movs)
	}
	//--Eliminar el carril
	var carril Carriles
	db.Where("id = ?", id).Delete(&carril)
	c.JSON(200, gin.H{"msg": "Carril eliminado"})
}

// --Funcion para eliminar los movimientos
func DeleteMovimiento(c *gin.Context) {
	id := c.PostForm("id")
	var movimiento Movimientos
	db.Where("id = ?", id).Delete(&movimiento)
	c.JSON(200, gin.H{"msg": "Movimiento eliminado"})
}

// --Funcion para eliminar las fases
func DeleteFase(c *gin.Context) {
	id := c.PostForm("id")
	var grupos []Grupo_carril
	db.Where("Fase = ?", id).Find(&grupos)
	for _, grupo := range grupos {
		db.Where("id = ?", grupo.ID).Delete(&grupos)
	}
	var fase Fase_semaforo
	db.Where("id = ?", id).Delete(&fase)

	c.JSON(200, gin.H{"msg": "Fase eliminado"})
}

// --Funcion para eliminar los grupos
func DeleteGrupo(c *gin.Context) {
	id := c.PostForm("id")
	var grupo Grupo_carril
	db.Where("id = ?", id).Delete(&grupo)
	c.JSON(200, gin.H{"msg": "Grupo eliminado"})
}

// -----------------------------------------------------------------
// ----------General -----------------------------
// -----------------------------------------------------------------
func GetAccesos(id uint) Accesos {
	var acceso Accesos
	db.Where("ID = ?", id).Find(&acceso)
	return acceso
}

func GetFase(id uint) Fase_semaforo {
	var fase Fase_semaforo
	db.Where("ID = ?", id).Find(&fase)
	return fase
}

func NivelServicio(demora float32) string {
	var result string
	if demora == 0 {
		result = "POR DEFINIR"
	} else if demora <= 10 {
		result = "A"
	} else if (demora > 10) && (demora < 20) {
		result = "B"
	} else if (demora > 20) && (demora < 35) {
		result = "C"
	} else if (demora > 35) && (demora < 55) {
		result = "D"
	} else if (demora > 55) && (demora < 80) {
		result = "E"
	} else {
		result = "F"
	}
	return result
}

func Redondear(result float32, lugares int) float32 {
	var mostrar float32

	if lugares == 1 {
		mostrar = float32(math.Round(float64(result)))
	} else if lugares == 2 {
		redondeo, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", result), 32)
		mostrar = float32(redondeo)
	} else if lugares == 3 {
		redondeo, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", result), 32)
		mostrar = float32(redondeo)
	} else {
		redondeo, _ := strconv.ParseFloat(fmt.Sprintf("%.4f", result), 32)
		mostrar = float32(redondeo)
	}

	return float32(mostrar)
}

// -----------------------------------------------------------------
// ----------Filters para los template -----------------------------
// -----------------------------------------------------------------
// --Cantidad de acceso de la interseccion
func (i Intersecciones) CantidadAccesos() int {
	var acceso []Accesos
	db.Where("Interseccion = ?", i.ID).Find(&acceso)
	return len(acceso)
}

// --Convertir el id del acceso en string
func (a Accesos) IdConvString() string {
	return strconv.Itoa(int(a.ID))
}

// --Adicionar 1 a la llave del range
func (c Carriles) MasUno(key int) int {
	return key + 1
}

// --Cantidad de carriles por accesos
func (a Accesos) CountCarrilesPorAccesos() int {
	var carriles []Carriles
	db.Where("Acceso = ?", a.ID).Find(&carriles)
	return len(carriles)
}

// --Carriles por acceso
func (a Accesos) CarrilesPorAccesos() []Carriles {
	var carriles []Carriles
	db.Where("Acceso = ?", a.ID).Find(&carriles)
	return carriles
}

// --Movimientos por carriles
func (c Carriles) MoviminetosPorCarriles() []Movimientos {
	var movimineto []Movimientos
	db.Where("Carril = ?", c.ID).Find(&movimineto)
	return movimineto
}

// --Grupos de carriles por acceso
func (a Accesos) GruposPorAccesos() []Grupo_carril {
	var grupo []Grupo_carril
	db.Where("Acceso = ?", a.ID).Find(&grupo)
	return grupo
}

// --Cantidad de Grupos de carriles por acceso
func (a Accesos) CantidadDeGruposPorAccesos() int {
	var grupo []Grupo_carril
	db.Where("Acceso = ?", a.ID).Find(&grupo)
	return len(grupo)
}

// --Cantidad de  Movimientos por carriles
func (c Carriles) CantMoviminetosPorCarriles() int {
	var movimineto []Movimientos
	db.Where("Carril = ?", c.ID).Find(&movimineto)
	return len(movimineto)
}

// --Movimientos por carriles
func (m Movimientos) VolumenPorMoviminetos() int {
	suma := m.Autos + m.Motos + m.Ciclos + m.Coches + m.CamionesCarga + m.CamionesPasaje + m.Omnibus
	return suma
}

// -- Grupo al que pertenece el carril
func (g Grupo_carril) GrupoDelCarril(c Carriles) string {
	var result string
	carriles := strings.Split(g.Carriles, ",")
	for i := 0; i < len(carriles)-1; i++ {

		if carriles[i] == string(strconv.Itoa(int(c.ID))) {

			result += g.NombreGrupo()
		}
	}
	if result != "" {
		return result
	} else {
		return "null"
	}

}

// -- Definir el nombre del grupo del carril
func (g Grupo_carril) NombreGrupo() string {
	carriles := strings.Split(g.Carriles, ",")
	var movimientos []Movimientos
	izq := 0
	ret := 0
	der := 0
	var nombre string
	for i := 0; i < len(carriles)-1; i++ {
		var carril Carriles
		db.Where("id = ?", carriles[i]).First(&carril)
		db.Where("Carril = ?", carril.ID).Find(&movimientos)
		for _, b := range movimientos {
			if b.Tipo == "I" {
				izq += 1
			} else if b.Tipo == "R" {
				ret += 1
			} else {
				der += 1
			}
		}
	}
	if (izq > 0) && (ret == 0) && (der == 0) {
		nombre += "I"
	} else if (izq > 0) && (ret > 0) && (der == 0) {
		nombre += "RI"
	} else if (izq == 0) && (ret > 0) && (der == 0) {
		nombre += "R"
	} else if (izq == 0) && (ret > 0) && (der > 0) {
		nombre += "RD"
	} else if (izq == 0) && (ret == 0) && (der > 0) {
		nombre += "D"
	} else if (izq > 0) && (ret > 0) && (der > 0) {
		nombre += "IRD"
	} else {
		nombre += "ID"
	}

	return nombre
}

// --Comprobar si el carril ya esta en un grupo
func (c Carriles) CarrilEnGrupo() bool {
	var acceso Accesos
	db.Where("ID = ?", c.Acceso).Find(&acceso)
	for _, grupo := range acceso.GruposPorAccesos() {
		carriles := strings.Split(grupo.Carriles, ",")
		for i := 0; i < len(carriles)-1; i++ {
			if carriles[i] == string(strconv.Itoa(int(c.ID))) {
				return true
			}
		}
	}
	return false
}

// --Comprobar se el carril esta en este grupo
func (g Grupo_carril) CarrilEnEsteGrupo(c Carriles) bool {
	carriles := strings.Split(g.Carriles, ",")
	for i := 0; i < len(carriles)-1; i++ {
		if carriles[i] == string(strconv.Itoa(int(c.ID))) {
			return true
		}
	}
	return false
}

// -----------------------------------------------------------------
// ----------Tabla 1 Modulo de ajuste de volumenes -----------------
// -----------------------------------------------------------------
// --Volumenes reales o equivalentes por movimientos
func (a Accesos) T1_volumenes(ordenes string) float32 {
	orden := strings.Split(ordenes, "&")
	var autos int
	var motos int
	var ciclos int
	var coches int
	var omnibus int
	var camionescarga int
	var camionespasage int
	var result float32
	var suma float32
	for _, carril := range a.CarrilesPorAccesos() {
		for _, movimiento := range carril.MoviminetosPorCarriles() {
			if orden[0] == movimiento.Tipo {
				autos += movimiento.Autos
				motos += movimiento.Motos
				ciclos += movimiento.Ciclos
				coches += movimiento.Coches
				omnibus += movimiento.Omnibus
				camionescarga += movimiento.CamionesCarga
				camionespasage += movimiento.CamionesPasaje
			}
		}
	}

	if orden[1] == "Real" {
		result += float32(autos) + float32(motos) + float32(ciclos) + float32(coches) + float32(omnibus) + float32(camionescarga) + float32(camionespasage)

	} else {
		suma += float32(autos) + float32(motos)/2 + float32(ciclos)/4 + float32(coches)/0.73 + float32(omnibus) + float32(camionescarga) + float32(camionespasage)
		result = Redondear(suma, 1)
	}

	return result

}

// --Intensidad reales o equivalentes por movimientos
func (a Accesos) T1_intensidad(ordenes string) float32 {
	var result float32
	vhmd := a.T1_volumenes(ordenes)
	fhmd, _ := strconv.ParseFloat(a.Fhmd, 32)
	result = float32(vhmd) / float32(fhmd)

	return float32(Redondear(result, 1))
}

// --Volumen del grupo de carril
func (g Grupo_carril) T1_volumen_grupo(ordenes string) float32 {
	carriles := strings.Split(g.Carriles, ",")
	orden := strings.Split(ordenes, "&")
	var mayor float32
	var volumen float32
	for i := 0; i < len(carriles)-1; i++ {
		var carrils Carriles
		db.Where("ID = ?", carriles[i]).Find(&carrils)
		for _, movimiento := range carrils.MoviminetosPorCarriles() {
			if orden[0] == "I" {
				if movimiento.Tipo == orden[0] {
					if orden[1] == "Real" {
						volumen += float32(movimiento.Autos) + float32(movimiento.Motos) + float32(movimiento.Ciclos) + float32(movimiento.Coches) + float32(movimiento.Omnibus) + float32(movimiento.CamionesCarga) + float32(movimiento.CamionesPasaje)
					} else {
						volumen += float32(movimiento.Autos) + float32(movimiento.Motos)/2 + float32(movimiento.Ciclos)/4 + float32(movimiento.Coches)/0.73 + float32(movimiento.Omnibus) + float32(movimiento.CamionesCarga) + float32(movimiento.CamionesPasaje)
					}
				}
			} else if orden[0] == "D" {
				if movimiento.Tipo == orden[0] {
					if orden[1] == "Real" {
						volumen += float32(movimiento.Autos) + float32(movimiento.Motos) + float32(movimiento.Ciclos) + float32(movimiento.Coches) + float32(movimiento.Omnibus) + float32(movimiento.CamionesCarga) + float32(movimiento.CamionesPasaje)
					} else {
						volumen += float32(movimiento.Autos) + float32(movimiento.Motos)/2 + float32(movimiento.Ciclos)/4 + float32(movimiento.Coches)/0.73 + float32(movimiento.Omnibus) + float32(movimiento.CamionesCarga) + float32(movimiento.CamionesPasaje)
					}
				}
			} else if orden[0] == "VP" {
				if orden[1] == "Real" {
					volumen += float32(movimiento.Omnibus) + float32(movimiento.CamionesCarga) + float32(movimiento.CamionesPasaje)
				} else {
					volumen += float32(movimiento.Omnibus)/0.33 + float32(movimiento.CamionesCarga)/0.33 + float32(movimiento.CamionesPasaje)/0.33
				}
			} else if orden[0] == "VC" {
				if orden[1] == "Real" {
					volumen += float32(movimiento.CamionesCarga)
				} else {
					volumen += float32(movimiento.CamionesCarga) / 0.33
				}
			} else if orden[0] == "VB" {
				if orden[1] == "Real" {
					volumen += float32(movimiento.CamionesPasaje) + float32(movimiento.Omnibus)
				} else {
					volumen += float32(movimiento.CamionesPasaje)/0.33 + float32(movimiento.Omnibus)/0.33
				}

			} else if orden[0] == "Mayor" {
				if orden[1] == "Real" {
					mayor += float32(movimiento.Autos) + float32(movimiento.Motos) + float32(movimiento.Ciclos) + float32(movimiento.Coches) + float32(movimiento.Omnibus) + float32(movimiento.CamionesCarga) + float32(movimiento.CamionesPasaje)
				} else {
					mayor += float32(movimiento.Autos) + float32(movimiento.Motos)/2 + float32(movimiento.Ciclos)/4 + float32(movimiento.Coches)/0.73 + float32(movimiento.Omnibus) + float32(movimiento.CamionesCarga) + float32(movimiento.CamionesPasaje)
				}
			} else {
				if orden[1] == "Real" {
					volumen += float32(movimiento.Autos) + float32(movimiento.Motos) + float32(movimiento.Ciclos) + float32(movimiento.Coches) + float32(movimiento.Omnibus) + float32(movimiento.CamionesCarga) + float32(movimiento.CamionesPasaje)
				} else {
					volumen += float32(movimiento.Autos) + float32(movimiento.Motos)/2 + float32(movimiento.Ciclos)/4 + float32(movimiento.Coches)/0.73 + float32(movimiento.Omnibus) + float32(movimiento.CamionesCarga) + float32(movimiento.CamionesPasaje)
				}
			}
		}
		if orden[0] == "Mayor" {
			if mayor > volumen {
				volumen = mayor
			}
			mayor = 0
		}
	}

	return float32(Redondear(volumen, 1))
}

// --Intensidad reales o equivalentes por movimientos
func (g Grupo_carril) T1_intensidad_grupo(ordenes string) float32 {
	var acceso Accesos
	var result float32
	vhmd := g.T1_volumen_grupo(ordenes)
	db.Where("ID = ?", g.Acceso).Find(&acceso)
	fhmd, _ := strconv.ParseFloat(acceso.Fhmd, 32)
	result = float32(vhmd) / float32(fhmd)
	return float32(Redondear(result, 1))
}

// --Proporcion de vueltas
func (g Grupo_carril) T1_ProporcionVueltas(ordenes string) float32 {
	var result float32
	vol_grupo := g.T1_volumen_grupo("t&Real")
	movimiento := g.T1_volumen_grupo(ordenes)
	result = movimiento / vol_grupo
	return float32(Redondear(result, 2))
}

// -----------------------------------------------------------------
// ----------Tabla 2 Modulo de intensidad de saturacion ------------
// -----------------------------------------------------------------
// --Numero de carriles por grupo
func (g Grupo_carril) T2_NumeroDeCarriles() int {
	carriles := strings.Split(g.Carriles, ",")
	var numero_de_carriles int
	for i := 0; i < len(carriles)-1; i++ {
		var carril Carriles
		db.Where("ID = ?", carriles[i]).Find(&carril)
		ancho, _ := strconv.ParseFloat(carril.Ancho, 32)
		if ancho > 5 {
			numero_de_carriles += 2
		} else {
			numero_de_carriles += 1
		}
	}
	return numero_de_carriles
}

// --Factor de Ancho
func (g Grupo_carril) T2_FactorAncho() float32 {
	var result float32
	carriles := strings.Split(g.Carriles, ",")
	var numero_de_carriles int
	var mayor_ancho float32
	for i := 0; i < len(carriles)-1; i++ {
		var carril Carriles
		db.Where("ID = ?", carriles[i]).Find(&carril)
		ancho, _ := strconv.ParseFloat(carril.Ancho, 32)
		if float32(ancho) > mayor_ancho {
			mayor_ancho = float32(ancho)
		}
		numero_de_carriles += 1
	}
	var area float32
	if mayor_ancho > 5 {
		area = float32(mayor_ancho / 2)
	} else {
		area = float32(mayor_ancho)
	}
	result = float32(1 + ((area - 3.6) / 9))
	return Redondear(result, 2)
}

// --Factor de pendiente
func (a Accesos) T2_FactorPendiente() float32 {
	var result float32
	pendiente, _ := strconv.ParseFloat(a.Pendiente, 32)
	if pendiente != 0 {
		result = float32(1 - pendiente/200)
	} else {
		result = 1
	}
	return Redondear(result, 2)
}

// --Factor de estacionamiento
func (a Accesos) T2_FactorEstacionamiento(g Grupo_carril) float32 {
	var result float32
	numero_de_carriles := float32(g.T2_NumeroDeCarriles())
	maniobras_estacionamiento, _ := strconv.ParseFloat(a.Maniobrasestacionamiento, 32)

	if numero_de_carriles > 0 {
		if maniobras_estacionamiento > 0 {
			result = float32((numero_de_carriles - 0.1 - ((18 * float32(maniobras_estacionamiento)) / 3600)) / numero_de_carriles)
		} else {
			result = 1
		}
	} else {
		result = 1
	}

	return Redondear(result, 2)
}

// --Factor de BUCES
func (a Accesos) T2_FactorBuces(g Grupo_carril) float32 {
	var result float32
	numero_de_carriles := float32(g.T2_NumeroDeCarriles())
	maniobras_parada, _ := strconv.ParseFloat(a.Maniobrasparada, 32)
	if (a.Parada == "true") && (maniobras_parada > 0) {
		result = float32((numero_de_carriles - 0.1 - ((14.4 * float32(maniobras_parada)) / 3600)) / numero_de_carriles)
	} else {
		result = 1
	}
	return Redondear(result, 2)
}

// --Facator de utilizacion
func (g Grupo_carril) T2_FactorUtilizacion(ordenes string) float32 {
	var result float32
	orden := strings.Split(ordenes, "&")
	numero_de_carriles := float32(g.T2_NumeroDeCarriles())
	volumen_total := g.T1_volumen_grupo(ordenes)
	volumen_mayor := g.T1_volumen_grupo(string("Mayor&" + orden[1]))
	result = volumen_total / (volumen_mayor * numero_de_carriles)
	return Redondear(result, 2)
}

// --Facator de vueltas
func (g Grupo_carril) T2_FactorVueltas(ordenes string) float32 {
	var result float32
	orden := strings.Split(ordenes, "&")
	volumen_total := g.T1_ProporcionVueltas(ordenes)
	if volumen_total == 0 {
		result = 1
	} else {
		if orden[0] == "I" {
			result = 1 / (1 + (0.05 * volumen_total))
		} else {
			result = 1 - (0.15 * volumen_total)
		}
	}
	return Redondear(result, 2)
}

// --Facator de vehiculos pesados
func (g Grupo_carril) T2_FactorVP(ordenes string) float32 {
	var pc float32
	var pb float32
	var result float32
	orden := strings.Split(ordenes, "&")
	volumen_total := g.T1_volumen_grupo(ordenes)
	volumen_vc := g.T1_volumen_grupo(string("VC&" + orden[1]))
	volumen_vb := g.T1_volumen_grupo(string("VB&" + orden[1]))
	pc = (volumen_vc / volumen_total) * 100
	pb = (volumen_vb / volumen_total) * 100
	result = 100 / (100 + (pc * (3 - 1)) + (pb * (2.5 - 1)))
	return Redondear(result, 2)
}

// --Saturacion
func (g Grupo_carril) T2_Saturaccion(ordenes string) float32 {
	var acceso Accesos
	var interseccion Intersecciones
	db.Where("ID = ?", g.Acceso).Find(&acceso)
	db.Where("ID = ?", acceso.Interseccion).Find(&interseccion)
	orden := strings.Split(ordenes, "&")

	So, _ := strconv.ParseFloat(acceso.Intensidad, 32)
	N := float32(g.T2_NumeroDeCarriles())
	fA := float32(g.T2_FactorAncho())
	fi := float32(acceso.T2_FactorPendiente())
	fa, _ := strconv.ParseFloat(interseccion.Area, 32)
	fe := float32(acceso.T2_FactorEstacionamiento(g))
	fbb := float32(acceso.T2_FactorBuces(g))
	fu := float32(g.T2_FactorUtilizacion(ordenes))
	fvi := float32(g.T2_FactorVueltas(string("I&" + orden[1])))
	fvd := float32(g.T2_FactorVueltas(string("D&" + orden[1])))
	fvp := float32(g.T2_FactorVP(ordenes))
	Si := float32(So) * float32(N) * float32(fA) * float32(fi) * float32(fa) * float32(fe) * float32(fbb) * float32(fu) * float32(fvi) * float32(fvd) * float32(fvp)
	return Redondear(Si, 1)
}

// -----------------------------------------------------------------
// ----------Tabla 3 Modulo de analisis de capacidad ---------------
// -----------------------------------------------------------------
// --Verde efectivo
func (g Grupo_carril) T3_VerdeEfectivo() float32 {
	tpad, _ := strconv.ParseFloat(GetAccesos(uint(g.Acceso)).Tpad, 32)
	verde := float32(GetFase(uint(g.Fase)).Verde)
	return float32(verde) - float32(tpad)
}

// --Verde efectivo entre el ciclo
func (g Grupo_carril) T3_VerdeEfectivoEntreCiclo() float32 {
	var result float32
	gi := g.T3_VerdeEfectivo()
	ciclo := float32(GetFase(uint(g.Fase)).Ciclo)
	result = float32(gi) / float32(ciclo)
	return Redondear(result, 2)
}

// --Ciclo del semaforo
func (g Grupo_carril) T3_Ciclo() float32 {
	var result float32
	result = float32(GetFase(uint(g.Fase)).Ciclo)
	return Redondear(result, 1)
}

// --Capacidad
func (g Grupo_carril) T3_Capacidad(ordenes string) float32 {
	var result float32
	Si := g.T2_Saturaccion(ordenes)
	gi_c := g.T3_VerdeEfectivoEntreCiclo()
	result = float32(Si) * float32(gi_c)
	return Redondear(result, 1)
}

// --
func (g Grupo_carril) T3_Xi(ordenes string) float32 {
	var result float32
	I := g.T1_intensidad_grupo(ordenes)
	capacidad := g.T3_Capacidad(ordenes)
	result = float32(I) / float32(capacidad)
	return Redondear(result, 2)
}

// --Intensidad entre saturacion
func (g Grupo_carril) T3_I_S(ordenes string) float32 {
	var result float32
	I := g.T1_intensidad_grupo(ordenes)
	Si := g.T2_Saturaccion(ordenes)
	result = float32(I) / float32(Si)
	return Redondear(result, 2)
}

// --Grupo Critico
func (g Grupo_carril) T3_GrupoCritico(ordenes string) string {
	var result string
	acceso := GetAccesos(uint(g.Acceso))
	cantidad_grupo := acceso.CantidadDeGruposPorAccesos()
	if cantidad_grupo == 1 {
		result = "X"
	} else {
		result = "X"
		I_S := g.T3_I_S(ordenes)
		for _, grupo := range acceso.GruposPorAccesos() {
			if grupo.ID != g.ID {
				I_S_G := grupo.T3_I_S(ordenes)
				if I_S_G > I_S {
					result = ""
				}
			}
		}
	}

	return result
}

// -----------------------------------------------------------------
// ----------Tabla 4 Modulo de nivel de servicio     ---------------
// -----------------------------------------------------------------
// --Demora uniforme
func (g Grupo_carril) T4_DemoraUniforme(ordenes string) float32 {
	var result float32
	gi_c := g.T3_VerdeEfectivoEntreCiclo()
	ciclo := float32(GetFase(uint(g.Fase)).Ciclo)
	xi := g.T3_Xi(ordenes)
	result = (0.38 * float32(ciclo) * (float32(math.Pow(float64((1 - float32(gi_c))), 2))) / (1 - (float32(gi_c) * float32(xi))))
	return Redondear(result, 2)
}

// --Demora incremental
func (g Grupo_carril) T4_DemoraIncremental(ordenes string) float32 {
	var result float32
	xi := g.T3_Xi(ordenes)
	capacidad := g.T3_Capacidad(ordenes)
	result = float32(173 * float32(math.Pow(float64(xi), 2)) * ((float32(xi) - 1) + float32(math.Sqrt(float64(math.Pow((float64(xi)-1), 2)+16*float64(xi)/float64(capacidad))))))
	return Redondear(result, 2)
}

// --Demora del Grupo
func (g Grupo_carril) T4_DemoraGrupo(ordenes string) float32 {
	var result float32
	du := g.T4_DemoraUniforme(ordenes)
	di := g.T4_DemoraIncremental(ordenes)
	result = float32(float64(du) + float64(di)*1)
	return Redondear(result, 2)
}

// --Nivel de servicio del grupo
func (g Grupo_carril) T4_NivelServicioGrupo(ordenes string) string {
	dg := g.T4_DemoraGrupo(ordenes)
	NS := NivelServicio(dg)
	return NS
}

// --Demora del acceso
func (a Accesos) T4_DemoraAcceso(ordenes string) float32 {
	var dia float32
	var i float32
	var result float32
	for _, grupo := range a.GruposPorAccesos() {
		dia += float32(grupo.T4_DemoraGrupo(ordenes)) * float32(grupo.T1_intensidad_grupo(ordenes))
		i += float32(grupo.T1_intensidad_grupo(ordenes))
	}
	result = dia / i
	return Redondear(result, 2)
}

// --Nivel de servicio del acceso
func (a Accesos) T4_NivelServicioAcceso(ordenes string) string {
	da := a.T4_DemoraAcceso(ordenes)
	NS := NivelServicio(da)
	return NS
}

// --Demora de la interseccion
func (i Intersecciones) T4_DemoraInterseccion(ordenes string) float32 {
	var result float32
	var accesos []Accesos
	db.Where("Interseccion = ?", i.ID).Find(&accesos)
	var valor float32
	var IG float32
	for _, acceso := range accesos {
		var i_g float32
		da := acceso.T4_DemoraAcceso(ordenes)
		for _, grupo := range acceso.GruposPorAccesos() {
			i_g += grupo.T1_intensidad_grupo(ordenes)
		}
		valor += float32(da) * float32(i_g)
		IG += float32(i_g)
	}
	result = float32(valor) / float32(IG)
	if result > 0 {
		return Redondear(result, 2)
	}
	return 0

}

// --Nivel de servicio de la interseccion
func (i Intersecciones) T4_NivelServicioInterseccion(ordenes string) string {
	di := i.T4_DemoraInterseccion(ordenes)
	NS := NivelServicio(di)
	return NS
}
