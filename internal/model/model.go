package model

type DataUsrReg struct {
	Mail string `json:"mail" validate:"required"`
	Pwd  string `json:"pwd" validate:"required"`
}

type DataUsrGet struct {
	IdUsr string `json:"id_usr"`
	Mail  string `json:"mail"`
	Pwd   string `json:"pwd"`
	St    string `json:"estatus"`
	Opt   string `json:"opt"`
}

type DataListProduct struct {
	IdP      int32   `json:"id"`
	IdUsr    int32   `json:"id_usr"`
	IdProd   int32   `json:"id_product"`
	Opt      string  `json:"opt"`
	Name     string  `json:"name"`
	Sku      string  `json:"sku"`
	Quantity int32   `json:"quantity"`
	Price    float32 `json:"price"`
	St       string  `json:"estatus"`
}

type DataUsr struct {
	Id   string `json:"id_usr"`
	Mail string `json:"mail"`
	Opt  string `json:"opt"`
	St   string `json:"estatus"`
}

type DataProduct struct {
	Name     string `json:"name"`
	Sku      string `json:"sku"`
	Quantity string `json:"quantity"`
	Price    string `json:"price"`
	IdUsr    string `json:"id_usr"`
	Mail     string `json:"mail"`
}

type DataProductSearch struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Sku  string `json:"sku"`
	PI   string `json:"pi"`
	PF   string `json:"pf"`
}

type DataCancelProduct struct {
	IdUsr string `json:"id_usr"`
	Id    string `json:"id"`
}

type DataCustomerReg struct {
	Mail string `json:"mail" validate:"required"`
	Name string `json:"name" validate:"required"`
	Pwd  string `json:"pwd" validate:"required"`
}

type DataCustomer struct {
	Id   string `json:"id"`
	Mail string `json:"mail"`
	Name string `json:"name"`
	Pwd  string `json:"pwd"`
	St   string `json:"estatus"`
}

type DataRegCarProduct struct {
	IdC string
	IdP string
}

type DataCarProduct struct {
	IdCP string `json:"idcp"`
	Mail string `json:"mail"`
	Name string `json:"name"`
	St   string `json:"estatus"`
}

type DataCustom struct {
	IdC  string `json:"idc`
	Mail string `json:"mail"`
	St   string `json:"st"`
}

type DataCustomCar struct {
	IdC       int32   `json:"id_customer"`
	IdCarProd int32   `json:"id_CP"`
	IdUsr     int32   `json:"id_usr"`
	IdProd    int32   `json:"id_product"`
	IdCar     int32   `json:"id_car"`
	Mail      string  `json:"mail"`
	Name      string  `json:"name"`
	Sku       string  `json:"sku"`
	Quantity  int32   `json:"quantity"`
	Price     float32 `json:"price"`
	St        string  `json:"estatus"`
}
