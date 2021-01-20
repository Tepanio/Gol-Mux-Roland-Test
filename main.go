package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
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

	db.Exec("CREATE DATABASE IF NOT EXISTS orders_db")

	db.Exec("USE orders_db")
	db.AutoMigrate(&Pedido{}, &Producto{})
}

func getPedidos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "aplication/json")
	var pedidos []Pedido
	db.Preload("Productos").Find(&pedidos)
	json.NewEncoder(w).Encode(pedidos)
}

func getPedido(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "aplication/json")
	//Se cargan las variables
	params := mux.Vars(r)
	//Se obtiene la id
	pedidoID := params["id"]

	var pedido Pedido
	db.Preload("Items").First(&pedido, pedidoID)
	json.NewEncoder(w).Encode(pedido)

}

func createPedido(w http.ResponseWriter, r *http.Request) {
	var pedido Pedido

	json.NewDecoder(r.Body).Decode(&pedido)
	db.Create(&pedido)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pedido)
}

func updatePedido(w http.ResponseWriter, r *http.Request) {
	var updatedPedido Pedido

	json.NewDecoder(r.Body).Decode(&updatedPedido)

	db.Save(&updatedPedido)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedPedido)
}
func deletePedido(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	pedidoID := params["id"]
	id64, _ := strconv.ParseUint(pedidoID, 10, 24)

	pedidoIDu := uint(id64)

	db.Where("pedido_id = ?", pedidoIDu).Delete(&Producto{})
	db.Where("pedido_id = ?", pedidoIDu).Delete(&Pedido{})
	w.WriteHeader(http.StatusNoContent)

}
func main() {
	router := mux.NewRouter()
	//Get Pedios
	router.HandleFunc("/api/pedidos", getPedidos).Methods("GET")
	router.HandleFunc("/api/pedido/{id}", getPedido).Methods("GET")
	router.HandleFunc("/api/pedido", createPedido).Methods("POST")
	router.HandleFunc("/api/pedido/{id}", updatePedido).Methods("PUT")
	router.HandleFunc("/api/pedido/{id}", deletePedido).Methods("DELETE")

	initDB()
	port := ":8080"
	fmt.Println("Api Runing Under ")

	log.Fatal(http.ListenAndServe(port, router))
}
