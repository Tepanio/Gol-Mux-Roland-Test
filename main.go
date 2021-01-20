package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

//Pedido sadasd
///Estructura del Pedido
type Pedido struct {
	PedidoID        uint       `json:"oderdId" gorm:"primary_key"`
	VendedorNombre  string     `json:"vendedorNombre"`
	FechaSolicitado time.Time  `json:"fechaSolicitado"`
	Productos       []Producto `json:"productos" gorm:"foreignkey:PedidoID"`
}

//Producto estructura
type Producto struct {
	ProductoID     uint   `json:"lineItemId" gorm:"primary_key"`
	ProductoCodigo string `json:"productoCodigo"`
	Descripcion    string `json:"descripcion"`
	Cantidad       uint   `json:"cantidad"`
	PedidoID       uint   `json:"-"`
}

//////Base de datos

var db *gorm.DB

func initDB() {
	var err error
	dataSourceName := "root:@tcp(localhost:3306)/?parseTime=True"
	db, err = gorm.Open("mysql", dataSourceName)

	if err != nil {
		fmt.Println(err)
		panic("Fallo la conexion con la base de datos")
	}

	db.Exec("CREATE DATABASE orders_db")

	db.Exec("USE orders_db")
	db.AutoMigrate(&Pedido{}, &Producto{})
}

func getPedidos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "aplication/json")
	var pedidos []Pedido
	db.Preload("Productos").Find(&pedidos)
	json.NewEncoder(w).Encode(pedidos)
}
func main() {
	router := mux.NewRouter()
	//Get Pedios
	router.HandleFunc("/api/pedidos", getPedidos).Methods("GET")
	initDB()
	port := ":8080"
	fmt.Println("Api Runing Under ")

	log.Fatal(http.ListenAndServe(port, router))
}
