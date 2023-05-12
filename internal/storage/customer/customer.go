package customer

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"comprarmas.com.mx/internal/model"

	"comprarmas.com.mx/internal/database"
)

func RegisterCustomer(sMail string, sName string, sPwd string) (string, error) {
	var lastId int64

	db, err := database.CnnDB_ENV()
	if err != nil {
		log.Println("RegisterCustomer_Error " + err.Error())
		return "", err
	}
	defer db.Close()

	// Ahora vemos si tenemos conexión
	err = db.Ping()
	if err != nil {
		log.Println("RegisterCustomer_Error_Ping " + err.Error())
		return "", err
	}

	defer db.Close()
	sSQL := "INSERT INTO Customer(Customer, Mail, Pwd) values('" + sName + "', '" + sMail + "', '" + sPwd + "') "

	resDB, err := db.Exec(sSQL)
	if err != nil {
		log.Println("No fue posible registrar al cliente: " + err.Error())
		err = errors.New("No fue posible registrar al cliente: " + err.Error())
		return "", err
	} else {
		lastId, err = resDB.LastInsertId()
		log.Println(resDB)
	}
	return fmt.Sprint(lastId), nil
} //END RegisterCustomer

func SearchCustomer(sId string, sMail string, sName string) (string, string, error) {
	blnOK := true

	db, err := database.CnnDB_ENV()
	if err != nil {
		log.Println("SearchCustomer_Error cnn DB " + err.Error())
		return "", "", err
	}
	defer db.Close()

	// Ahora vemos si tenemos conexión
	err = db.Ping()
	if err != nil {
		log.Println("SearchCustomer_Error cnn DB " + err.Error())
		err = errors.New("SearchCustomer_Error cnn DB " + err.Error())
		return "", "", err
	}

	defer db.Close()
	sSQL := "SELECT DISTINCT a.id_Customer, a.Mail, a.Pwd, a.Customer, a.Estatus "
	sSQL += "FROM Customer a "
	sSQL += "WHERE (a.Estatus = 'ACTIVE') "

	if sId != "" {
		sSQL += "AND (a.id_Customer = " + sId + ") "
	}
	if sMail != "" {
		sSQL += "AND (a.Mail = '" + sMail + "') "
	}

	if sName != "" {
		sSQL += "AND (a.Customer = '" + sName + "') "
	}

	resDat, err := db.Query(sSQL)
	if err != nil {
		return "", "", err
	} else {
		for resDat.Next() {
			var datCustom = model.DataCustomer{}
			resDat.Scan(
				&datCustom.Id, &datCustom.Mail, &datCustom.Pwd, &datCustom.Name, &datCustom.St,
			)
			blnOK = true
			return fmt.Sprint(datCustom.Id), fmt.Sprint(datCustom.Pwd), nil
		}
	}
	if blnOK == false {
		err = errors.New("Cliente no encontrado.")
		return "", "", err
	}
	return "", "", nil
} //END SearchCustomer

func CustomerCarProducts(sIdCustom string) ([]map[string]interface{}, error) {
	var allData []map[string]interface{}

	db, err := database.CnnDB_ENV()
	if err != nil {
		log.Println("GetListProducts_Error cnn DB " + err.Error())
		return nil, err
	}
	defer db.Close()

	// Ahora vemos si tenemos conexión
	err = db.Ping()
	if err != nil {
		log.Println("GetListProducts_Error_Ping " + err.Error())
		return nil, err
	}

	defer db.Close()

	sSQL := "SELECT DISTINCT c.id_CP, a.id_Customer, b.id_Car, c.id_Product, a.Mail, a.Estatus,  "
	sSQL += "d.Name, d.Sku, c.Quantity, d.Price "
	sSQL += "FROM Customer a "
	sSQL += "INNER JOIN Car b ON(a.id_Customer=b.id_Customer) "
	sSQL += "INNER JOIN CarProduct c ON(b.id_Car=c.id_Car) "
	sSQL += "INNER JOIN Product d ON(c.id_Product=d.id_Product) "
	sSQL += "WHERE (a.Estatus = 'ACTIVE') "
	if sIdCustom != "" {
		sSQL += "AND (a.id_Customer = " + sIdCustom + ") "
	}

	resDat, err := db.Query(sSQL)
	if err != nil {
		return nil, err
	} else {
		for resDat.Next() {
			var datList = model.DataCustomCar{}
			resDat.Scan(
				&datList.IdCarProd, &datList.IdC, &datList.IdCar, &datList.IdProd, &datList.Mail,
				&datList.St, &datList.Name, &datList.Sku, &datList.Quantity, &datList.Price,
			)
			mapJson := make(map[string]interface{})
			mapJson["idc"] = fmt.Sprint(datList.IdC)
			mapJson["id_car"] = fmt.Sprint(datList.IdCar)
			mapJson["id_product"] = fmt.Sprint(datList.IdProd)
			mapJson["sku"] = datList.Sku
			mapJson["name"] = datList.Name
			mapJson["quantity"] = fmt.Sprint(datList.Quantity)
			mapJson["price"] = fmt.Sprint(datList.Price)
			allData = append(allData, mapJson)
		}
	}

	sJson, err := json.Marshal(allData)
	fmt.Println(string(sJson))
	return allData, nil
} //END CustomerCarProducts

func RegisterCarCustomer(sIdC string) (string, error) {
	var lastId int64

	db, err := database.CnnDB_ENV()
	if err != nil {
		log.Println("RegisterCarCustomer_Error " + err.Error())
		return "", err
	}
	defer db.Close()

	// Ahora vemos si tenemos conexión
	err = db.Ping()
	if err != nil {
		log.Println("RegisterCarCustomer_Error_Ping " + err.Error())
		return "", err
	}

	defer db.Close()
	sSQL := "INSERT INTO Car(id_Customer) values(" + sIdC + ") "

	resDB, err := db.Exec(sSQL)
	if err != nil {
		log.Println("No fue posible registrar el carrito: " + err.Error())
		err = errors.New("No fue posible registrar el carrito: " + err.Error())
		return "", err
	} else {
		lastId, err = resDB.LastInsertId()
		log.Println(resDB)
	}
	return fmt.Sprint(lastId), nil
} //END RegisterCarCustomer

func RegisterCarProductCustomer(sIdCar string, sIdProd string) (string, error) {
	var lastId int64

	db, err := database.CnnDB_ENV()
	if err != nil {
		log.Println("RegisterCarCustomer_Error " + err.Error())
		return "", err
	}
	defer db.Close()

	// Ahora vemos si tenemos conexión
	err = db.Ping()
	if err != nil {
		log.Println("RegisterCarCustomer_Error_Ping " + err.Error())
		return "", err
	}

	defer db.Close()
	sSQL := "INSERT INTO CarProduct(id_Car, id_Product, Quantity) values(" + sIdCar + ", " + sIdProd + ", 1) "

	resDB, err := db.Exec(sSQL)
	if err != nil {
		log.Println("No fue posible agregar los productos al carrito: " + err.Error())
		err = errors.New("No fue posible agregar los productos al carrito: " + err.Error())
		return "", err
	} else {
		lastId, err = resDB.LastInsertId()
		log.Println(resDB)
	}
	return fmt.Sprint(lastId), nil
} //END RegisterCarProductCustomer

func SearchCarCustomer(sIdCustom string) (string, error) {
	blnOK := true

	db, err := database.CnnDB_ENV()
	if err != nil {
		log.Println("SearchCustomer_Error cnn DB " + err.Error())
		return "", err
	}
	defer db.Close()

	// Ahora vemos si tenemos conexión
	err = db.Ping()
	if err != nil {
		log.Println("SearchCustomer_Error cnn DB " + err.Error())
		err = errors.New("SearchCustomer_Error cnn DB " + err.Error())
		return "", err
	}

	defer db.Close()
	sSQL := "SELECT DISTINCT a.id_Customer, b.id_Car, a.Mail "
	sSQL += "FROM Customer a "
	sSQL += "INNER JOIN Car b ON(a.id_Customer=b.id_Customer) "
	sSQL += "WHERE (a.Estatus = 'ACTIVE') "

	if sIdCustom != "" {
		sSQL += "AND (a.id_Customer = " + sIdCustom + ") "
	}
	/*
		if sMail != "" {
			sSQL += "AND (a.Mail = '" + sMail + "') "
		}

		if sName != "" {
			sSQL += "AND (a.Customer = '" + sName + "') "
		}
	*/
	resDat, err := db.Query(sSQL)
	if err != nil {
		return "", err
	} else {
		for resDat.Next() {
			var datCustom = model.DataCustomCar{}
			resDat.Scan(
				&datCustom.IdC, &datCustom.IdCar, &datCustom.Mail,
			)
			blnOK = true
			return fmt.Sprint(datCustom.IdCar), nil
		}
	}
	if blnOK == false {
		err = errors.New("Cliente no encontrado.")
		return "", err
	}
	return "", nil
} //END SearchCarCustomer
