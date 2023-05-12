package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"comprarmas.com.mx/internal/model"
	"comprarmas.com.mx/internal/storage/customer"
	"comprarmas.com.mx/internal/storage/product"
	"comprarmas.com.mx/internal/storage/redis"
	"comprarmas.com.mx/internal/storage/usr"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func PostLogin(ctx *fiber.Ctx) error {
	log.Println("LoginUsr ...")
	var err error
	sId := ""
	sOpt := ""
	sPwd := ""

	log.Println(sId + sOpt + sPwd)
	file, _ := os.Create("cb2.log")
	log.SetOutput(file)
	defer file.Close()

	//sMail := ctx.FormValue("Email")
	//sPwd := ctx.FormValue("password")
	//log.Println(sMail, sPwd)

	if ctx.Get("Content-Type") != "application/json" {
		return ctx.Status(401).JSON(fiber.Map{
			"Error": err.Error(),
			"msg":   "Content-Type header is not application/json.",
		})
	}

	var datUsr = model.DataUsrReg{}
	eJson := json.NewDecoder(bytes.NewReader(ctx.Body()))
	//if eJson != nil {
	//	return err
	//}

	err = eJson.Decode(&datUsr)
	if err != nil {
		return err
	}

	//Autenticate
	sId, sOpt, sPwd, err = usr.LoginUsr(datUsr.Mail, "")
	if err != nil {
		return ctx.Status(401).JSON(fiber.Map{
			"Error": err.Error(),
			"id":    "",
			"msg":   "No fue posible autenticar al usuario.",
		})
	} else {
		//Checar el pwd
		hashComoByte := []byte(sPwd)
		bytePwd := []byte(datUsr.Pwd)
		error := bcrypt.CompareHashAndPassword(hashComoByte, bytePwd)
		if error == nil {
			//Register Redis
			err := redis.AddUsr("Id_" + sId)
			if err != nil {
				return ctx.Status(401).JSON(fiber.Map{
					"Error": err.Error(),
					"msg":   "No fue posible registrar al usuario en redis.",
				})
			}
			//log.Println("Login_6")
			return ctx.Status(200).JSON(fiber.Map{
				"OK":  "Access",
				"id":  sId,
				"opt": sOpt,
			})
		} else {
			log.Println(err)
			return ctx.Status(401).JSON(fiber.Map{
				"Error": "No fue posible autenticar al usuario.",
				"id":    "",
				"msg":   "No fue posible autenticar al usuario.",
			})
		}
	}
}

func PostRegisterUsr(ctx *fiber.Ctx) error {
	log.Println("PostRegisterUsr ...")
	var err error
	sToken := ""
	sId := ""
	log.Println(sToken + sId)

	if ctx.Get("Content-Type") != "application/json" {
		return ctx.Status(401).JSON(fiber.Map{
			"Error": err.Error(),
			"msg":   "Content-Type header is not application/json.",
		})
	}

	var datUsr = model.DataUsrReg{}
	eJson := json.NewDecoder(bytes.NewReader(ctx.Body()))
	//if eJson != nil {
	//	return err
	//}
	err = eJson.Decode(&datUsr)
	if err != nil {
		return err
	}
	//Buscar si ya existe el usuario
	err = usr.GetUsr(datUsr.Mail, "")
	if err != nil {
		return ctx.Status(401).JSON(fiber.Map{
			"Error": err.Error(),
			"id":    "0",
		})
	}

	//Codificar si pwd
	bytePwd := []byte(datUsr.Pwd)
	hash, err := bcrypt.GenerateFromPassword(bytePwd, bcrypt.DefaultCost) //DefaultCost es 10
	if err != nil {
		fmt.Println(err)
	}
	sPwd := string(hash)

	//Registrar el Usuario
	sId, err = usr.RegisterUsr(datUsr.Mail, sPwd)
	if err != nil {
		return ctx.JSON(fiber.Map{
			"Error": err.Error(),
			"id":    "0",
		})
	} else {
		// Buscar el dato completo que se registro
		return ctx.JSON(fiber.Map{
			"OK": "Usuario registrado",
			"id": sId,
		})
	}
} //END PostRegisterUsr

func GetListProductsFree(ctx *fiber.Ctx) error {
	var err error

	if ctx.Get("Content-Type") != "application/json" {
		return ctx.Status(401).JSON(fiber.Map{
			"Error": err.Error(),
			"msg":   "Content-Type header is not application/json.",
		})
	}
	log.Println(ctx.Body())
	var datUsr = model.DataUsr{}
	eJson := json.NewDecoder(bytes.NewReader(ctx.Body()))
	err = eJson.Decode(&datUsr)
	if err != nil {
		return ctx.Status(401).JSON(fiber.Map{
			"Error": err.Error(),
		})
	}

	//Obtener la lista de productos
	mapJson, aChk, errLP := product.ListProducts("", "")
	if errLP != nil {
		return ctx.Status(401).JSON(fiber.Map{
			"Error": errLP.Error(),
		})
	} else {
		sJson, err := json.Marshal(mapJson)
		if err != nil {
			return ctx.Status(401).JSON(fiber.Map{
				"Error": err.Error(),
			})
		}
		log.Println(string(sJson))
		log.Println(aChk)
		return ctx.JSON(string(sJson))
	}
	//return ctx.JSON(aChk)
} //GetListPrductsFree

func GetListProducts(ctx *fiber.Ctx) error {
	var err error

	if ctx.Get("Content-Type") != "application/json" {
		return ctx.Status(401).JSON(fiber.Map{
			"Error": err.Error(),
			"msg":   "Content-Type header is not application/json.",
		})
	}
	log.Println(ctx.Body())
	var datUsr = model.DataUsr{}
	eJson := json.NewDecoder(bytes.NewReader(ctx.Body()))
	err = eJson.Decode(&datUsr)
	if err != nil {
		return ctx.Status(401).JSON(fiber.Map{
			"Error": err.Error(),
		})
	}
	//CHK Redis
	err = redis.ChkUsr("Id_" + datUsr.Id)
	if err != nil {
		return ctx.JSON(fiber.Map{
			"Error": err.Error(),
			"msg":   "Se ha terminado la sesión del usuario.",
		})
	}

	//Obtener la lista de productos
	mapJson, aChk, errLP := product.ListProducts(datUsr.Id, "")
	if errLP != nil {
		return ctx.Status(401).JSON(fiber.Map{
			"Error": errLP.Error(),
		})
	} else {
		sJson, err := json.Marshal(mapJson)
		if err != nil {
			return ctx.Status(401).JSON(fiber.Map{
				"Error": err.Error(),
			})
		}
		log.Println(string(sJson))
		log.Println(aChk)
		return ctx.JSON(string(sJson))
	}
	//return ctx.JSON(aChk)
}

func PostCancelProduct(ctx *fiber.Ctx) error {
	log.Println("CancelProduct ...")
	//sId := ctx.Params("id")
	var err error

	if ctx.Get("Content-Type") != "application/json" {
		return ctx.Status(401).JSON(fiber.Map{
			"Error": err.Error(),
			"msg":   "Content-Type header is not application/json.",
		})
	}
	log.Println(ctx.Body())
	var datUsr = model.DataCancelProduct{}
	eJson := json.NewDecoder(bytes.NewReader(ctx.Body()))
	err = eJson.Decode(&datUsr)
	if err != nil {
		return ctx.Status(401).JSON(fiber.Map{
			"Error": err.Error(),
		})
	}

	//CHK Redis
	err = redis.ChkUsr("Id_" + datUsr.IdUsr)
	if err != nil {
		return ctx.Status(401).JSON(fiber.Map{
			"Error": err.Error(),
			"msg":   "Se ha terminado la sesión del usuario.",
		})
	}

	err = product.CancelProduct(datUsr.Id)
	if err != nil {
		return ctx.Status(401).JSON(fiber.Map{
			"Error": err.Error(),
			"mag":   "No fue posible cancelar el producto.",
		})
	} else {
		// Buscar el dato completo que se registro
		return ctx.Status(200).JSON(fiber.Map{
			"process": "Cancelar.",
			"msg":     "Regisro cancelado.",
		})
	}
}

func PostRegProduct(ctx *fiber.Ctx) error {
	log.Println("PostRegProductLogin ...")
	var err error

	if ctx.Get("Content-Type") != "application/json" {
		return ctx.Status(401).JSON(fiber.Map{
			"Error": err.Error(),
			"msg":   "Content-Type header is not application/json.",
		})
	}

	var datProd = model.DataProduct{}
	eJson := json.NewDecoder(bytes.NewReader(ctx.Body()))
	//if eJson != nil {
	//	return err
	//}
	err = eJson.Decode(&datProd)
	if err != nil {
		return err
	}
	//Buscar si el usuario esta en session
	err = redis.ChkUsr("Id_" + datProd.IdUsr)
	if err != nil {
		return ctx.Status(401).JSON(fiber.Map{
			"Error": err.Error(),
			"msg":   "Se ha terminado la sesión del usuario.",
		})
	}

	err = product.ExistProduct(datProd.Name, datProd.Sku)
	if err != nil {
		return ctx.Status(401).JSON(fiber.Map{
			"Error": err.Error(),
			"msg":   "El producto ya existe.",
		})
	}

	//Register product
	sId, err := product.RegisterProduct(datProd.Name, datProd.Sku, fmt.Sprint(datProd.Quantity), fmt.Sprint(datProd.Price))
	if err != nil {
		return ctx.Status(202).JSON(fiber.Map{
			"Error": "No fue registrado el producto",
			"msg":   "El producto ya existe",
		})
	} else {
		err = product.RegisterUsrProduct(sId, fmt.Sprint(datProd.IdUsr))
		if err != nil {
			//err = errors.New("No fue posible registrar el producto.")
			return ctx.Status(401).JSON(fiber.Map{
				"Error": err.Error(),
				"msg":   "No fue posible registrar el producto.",
			})
		}
		return ctx.Status(200).JSON(fiber.Map{
			"OK":  "Access",
			"msg": "El producto se registro correctamente.",
		})
	}
} //END PostRegProduct

func PostProductSearch(ctx *fiber.Ctx) error {
	log.Println("PostRegProductLogin ...")
	var err error

	if ctx.Get("Content-Type") != "application/json" {
		return ctx.Status(401).JSON(fiber.Map{
			"Error": err.Error(),
			"msg":   "Content-Type header is not application/json.",
		})
	}

	var datProd = model.DataProductSearch{}
	eJson := json.NewDecoder(bytes.NewReader(ctx.Body()))
	//if eJson != nil {
	//	return err
	//}
	err = eJson.Decode(&datProd)
	if err != nil {
		return err
	}
	//Buscar si el usuario esta en session
	//err = redis.ChkUsr("")
	//if err != nil {
	//	return ctx.JSON(fiber.Map{
	//		"Error": err.Error(),
	//		"msg":   "Se ha terminado la sesión del usuario.",
	//	})
	//}
	mapJson, err := product.ListProductsSearch(datProd.Id, datProd.Name, datProd.Sku, datProd.PI, datProd.PF)
	if err != nil {
		return ctx.Status(401).JSON(fiber.Map{
			"Error": err.Error(),
			"msg":   "No fue posible encontrar productos.",
		})
	} else {
		sJson, err := json.Marshal(mapJson)
		if err != nil {
			return ctx.Status(401).JSON(fiber.Map{
				"Error": err.Error(),
			})
		}
		log.Println(string(sJson))
		return ctx.Status(200).JSON(string(sJson))
	}
} //END PostProductSearch

func PostProductAdmin(ctx *fiber.Ctx) error {
	log.Println("PostRegProductLogin ...")
	var err error

	//Buscar si el usuario esta en session
	//err = redis.ChkUsr("")
	//if err != nil {
	//	return ctx.JSON(fiber.Map{
	//		"Error": err.Error(),
	//		"msg":   "Se ha terminado la sesión del usuario.",
	//	})
	//}
	mapJson, err := usr.GetUsrsSeller()
	if err != nil {
		return ctx.JSON(fiber.Map{
			"Error": err.Error(),
			"msg":   "No fue posible encontrar productos.",
		})
	} else {
		sJson, err := json.Marshal(mapJson)
		if err != nil {
			return ctx.JSON(fiber.Map{
				"Error": err.Error(),
			})
		}
		log.Println(string(sJson))
		return ctx.JSON(string(sJson))
	}
}

func PostLoginCustomer(ctx *fiber.Ctx) error {
	log.Println("PostLoginCustomer ...")
	var err error
	sId := ""
	sOpt := ""
	sPwd := ""

	log.Println(sId + sOpt + sPwd)
	file, _ := os.Create("cb2.log")
	log.SetOutput(file)
	defer file.Close()

	if ctx.Get("Content-Type") != "application/json" {
		return ctx.Status(401).JSON(fiber.Map{
			"Error": err.Error(),
			"msg":   "Content-Type header is not application/json.",
		})
	}

	var datCustom = model.DataCustomerReg{}
	eJson := json.NewDecoder(bytes.NewReader(ctx.Body()))
	//if eJson != nil {
	//	return err
	//}

	err = eJson.Decode(&datCustom)
	if err != nil {
		return err
	}

	//Autenticate
	sId, sPwd, err = customer.SearchCustomer("", datCustom.Mail, "")
	if err != nil {
		return ctx.Status(401).JSON(fiber.Map{
			"Error": err.Error(),
			"id":    "",
			"msg":   "No fue posible autenticar al usuario.",
		})
	} else {
		//Checar el pwd
		hashComoByte := []byte(sPwd)
		bytePwd := []byte(datCustom.Pwd)
		error := bcrypt.CompareHashAndPassword(hashComoByte, bytePwd)
		if error == nil {
			//Register Redis
			err := redis.AddUsr("IdC_" + sId)
			if err != nil {
				return ctx.Status(401).JSON(fiber.Map{
					"Error": err.Error(),
					"msg":   "No fue posible verificar al cliente en redis.",
				})
			}
			//log.Println("Login_6")
			return ctx.Status(200).JSON(fiber.Map{
				"OK":  "Access",
				"idc": sId,
			})
		} else {
			log.Println(err)
			return ctx.Status(401).JSON(fiber.Map{
				"Error": "No fue posible autenticar al usuario.",
				"id":    "",
				"msg":   "No fue posible autenticar al usuario.",
			})
		}
	}
}

func PostRegisterCustomer(ctx *fiber.Ctx) error {
	log.Println("PostRegisterUsr ...")
	var err error

	sId := ""
	log.Println(sId)

	if ctx.Get("Content-Type") != "application/json" {
		return ctx.Status(401).JSON(fiber.Map{
			"Error": err.Error(),
			"msg":   "Content-Type header is not application/json.",
		})
	}

	var datCustom = model.DataCustomerReg{}
	eJson := json.NewDecoder(bytes.NewReader(ctx.Body()))
	//if eJson != nil {
	//	return err
	//}
	err = eJson.Decode(&datCustom)
	if err != nil {
		return err
	}
	//Buscar si ya existe el cliente
	_, _, err = customer.SearchCustomer("", datCustom.Mail, "")
	if err == nil {
		return ctx.Status(401).JSON(fiber.Map{
			"Error": "Ya existe el cliente.",
			"id":    "0",
		})
	}

	//Codificar si pwd
	bytePwd := []byte(datCustom.Pwd)
	hash, err := bcrypt.GenerateFromPassword(bytePwd, bcrypt.DefaultCost) //DefaultCost es 10
	if err != nil {
		fmt.Println(err)
	}
	sPwd := string(hash)

	//Registrar cliente
	sId, err = customer.RegisterCustomer(datCustom.Mail, datCustom.Name, sPwd)
	if err != nil {
		return ctx.Status(401).JSON(fiber.Map{
			"Error": err.Error(),
			"id":    "0",
		})
	} else {
		// Buscar el dato completo que se registro
		return ctx.Status(200).JSON(fiber.Map{
			"OK": "Cliente registrado",
			"id": sId,
		})
	}
} //END PostRegisterCustomer

func PostRegisterProdCustomer(ctx *fiber.Ctx) error {
	log.Println("PostRegisterProdCustomer ...")
	var err error

	sId := ""
	log.Println(sId)

	if ctx.Get("Content-Type") != "application/json" {
		return ctx.Status(401).JSON(fiber.Map{
			"Error": err.Error(),
			"msg":   "Content-Type header is not application/json.",
		})
	}

	var datCustom = model.DataRegCarProduct{}
	eJson := json.NewDecoder(bytes.NewReader(ctx.Body()))
	//if eJson != nil {
	//	return err
	//}
	err = eJson.Decode(&datCustom)
	if err != nil {
		return err
	}
	//Buscar un carrito
	sIdCar, err := customer.SearchCarCustomer(datCustom.IdC)
	if err != nil {
		return ctx.Status(401).JSON(fiber.Map{
			"Error": "Hugo un problema al buscar un carrito de compras para el cliente.",
		})
	}
	if sIdCar == "" || sIdCar == "0" {
		sIdNewCar, err := customer.RegisterCarCustomer(datCustom.IdC)
		if err != nil {
			return ctx.Status(401).JSON(fiber.Map{
				"Error": "No fue posible registrar el carrito de compras.",
			})
		}
		sId, err = customer.RegisterCarProductCustomer(sIdNewCar, datCustom.IdP)
	} else {
		sId, err = customer.RegisterCarProductCustomer(sIdCar, datCustom.IdP)
	}
	if err != nil {
		return ctx.Status(401).JSON(fiber.Map{
			"Error": "No fue posible agregar los productos al carrito de compras.",
		})
	}
	return ctx.Status(200).JSON(fiber.Map{
		"OK":  "AddProdcto.",
		"msg": "El producto fue agregado correctamente.",
	})
	return nil
} //END PostRegisterUsr

func PostCustomerShowCar(ctx *fiber.Ctx) error {
	log.Println("PostCustomerShowCar ...")
	var err error

	if ctx.Get("Content-Type") != "application/json" {
		return ctx.Status(401).JSON(fiber.Map{
			"Error": err.Error(),
			"msg":   "Content-Type header is not application/json.",
		})
	}

	var datCustom = model.DataCustom{}
	eJson := json.NewDecoder(bytes.NewReader(ctx.Body()))
	//if eJson != nil {
	//	return err
	//}
	err = eJson.Decode(&datCustom)
	if err != nil {
		return err
	}
	//Buscar si el usuario esta en session
	//err = redis.ChkUsr("")
	//if err != nil {
	//	return ctx.JSON(fiber.Map{
	//		"Error": err.Error(),
	//		"msg":   "Se ha terminado la sesión del usuario.",
	//	})
	//}
	mapJson, err := customer.CustomerCarProducts(datCustom.IdC)
	if err != nil {
		return ctx.Status(401).JSON(fiber.Map{
			"Error": err.Error(),
			"msg":   "No fue posible encontrar productos.",
		})
	} else {
		sJson, err := json.Marshal(mapJson)
		if err != nil {
			return ctx.Status(401).JSON(fiber.Map{
				"Error": err.Error(),
			})
		}
		log.Println(string(sJson))
		return ctx.Status(200).JSON(string(sJson))
	}
} //END PostCustomerShowCar
