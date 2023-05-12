package usr

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"

	"comprarmas.com.mx/internal/model"

	"comprarmas.com.mx/internal/database"
)

func LoginUsr(sMail string, sPwd string) (string, string, string, error) {
	blnOK := false

	db, err := database.CnnDB_ENV()
	if err != nil {
		log.Println("LoginUsr_Error cnn DB " + err.Error())
		return "", "", "", err
	}
	defer db.Close()

	// Ahora vemos si tenemos conexión
	err = db.Ping()
	if err != nil {
		log.Println("LoginUsr_Error cnn DB " + err.Error())
		err = errors.New("LoginUsr_Error cnn DB " + err.Error())
		return "", "", "", err
	}

	defer db.Close()
	sSQL := "SELECT DISTINCT a.id_Usr, a.Mail, a.Pwd, a.Estatus, a.Opt "
	sSQL += "FROM Usr a "
	sSQL += "WHERE (a.Mail = '" + sMail + "') "
	if sPwd != "" {
		sSQL += "AND (a.Pwd = '" + sPwd + "') "
	}
	sSQL += "AND (a.Estatus = 'ACTIVE') "

	resDat, err := db.Query(sSQL)
	if err != nil {
		return "", "", "", err
	} else {
		for resDat.Next() {
			var datUsr = model.DataUsrGet{}
			resDat.Scan(
				&datUsr.IdUsr, &datUsr.Mail, &datUsr.Pwd, &datUsr.St, &datUsr.Opt,
			)
			blnOK = true
			return fmt.Sprint(datUsr.IdUsr), fmt.Sprint(datUsr.Opt), fmt.Sprint(datUsr.Pwd), nil
		}
	}
	if blnOK == false {
		err = errors.New("Usuario no encontrado.")
		return "", "", "", err
	}
	return "", "", "", nil
} //END LoginUsr

func GetUsr(sMail string, sSeller string) error {
	db, err := database.CnnDB_ENV()
	if err != nil {
		log.Println("GetUsr_Error cnn DB " + err.Error())
		return err
	}
	defer db.Close()
	var errOpt error

	// Ahora vemos si tenemos conexión
	err = db.Ping()
	if err != nil {
		log.Println("GetUsr_Error cnn DB " + err.Error())
		return err
	}

	defer db.Close()
	sSQL := "SELECT DISTINCT a.Mail "
	sSQL += "FROM Usr a "
	sSQL += "WHERE (a.Mail = '" + sMail + "') "
	if sSeller != "" {
		sSQL += "AND (a.Opt='" + sSeller + "')"
	}

	resDat, err := db.Query(sSQL)
	if err != nil {
		return err
	} else {
		for resDat.Next() {
			var datUsr = model.DataUsrGet{}
			resDat.Scan(
				&datUsr.Mail, &datUsr.St,
			)
			if datUsr.Mail == "" {
				errOpt = errors.New("Ya existe el usuario.")
				return errOpt
			}
			err = error(nil)
		}
	}
	return nil
} //END GetUsr

func RegisterUsr(sMail string, sPwd string) (string, error) {
	var lastId int64

	db, err := database.CnnDB_ENV()
	if err != nil {
		log.Println("RegisterUsr_Error " + err.Error())
		return "", err
	}
	defer db.Close()

	// Ahora vemos si tenemos conexión
	err = db.Ping()
	if err != nil {
		log.Println("RegisterUsr_Error_Ping " + err.Error())
		return "", err
	}

	defer db.Close()
	sSQL := "INSERT INTO Usr(Mail, Pwd) values('" + sMail + "', '" + sPwd + "') "

	resDB, err := db.Exec(sSQL)
	if err != nil {
		log.Println("No fue posible registrar el usuario: " + err.Error())

		return "", err
	} else {
		lastId, err = resDB.LastInsertId()
		log.Println(resDB)
	}
	return fmt.Sprint(lastId), nil
} //END RegisterUsr

func GetUsrsSeller() ([]map[string]interface{}, error) {
	var allData []map[string]interface{}

	db, err := database.CnnDB_ENV()
	if err != nil {
		log.Println("GetUsrSeller_Error cnn DB " + err.Error())
		return nil, err
	}
	defer db.Close()

	// Ahora vemos si tenemos conexión
	err = db.Ping()
	if err != nil {
		log.Println("GetUsrSeller_Error cnn DB " + err.Error())
		return nil, err
	}

	defer db.Close()
	sSQL := "SELECT DISTINCT a.id_Usr, a.Mail "
	sSQL += "FROM Usr a "
	sSQL += "WHERE (a.Opt='seller')"

	resDat, err := db.Query(sSQL)
	if err != nil {
		return nil, err
	} else {
		for resDat.Next() {
			var datUsr = model.DataUsrGet{}
			resDat.Scan(
				&datUsr.IdUsr, &datUsr.Mail,
			)

			mapJson := make(map[string]interface{})
			mapJson["id_usr"] = fmt.Sprint(datUsr.IdUsr)
			mapJson["mail"] = fmt.Sprint(datUsr.Mail)
			allData = append(allData, mapJson)
		}
	}

	sJson, err := json.Marshal(allData)
	fmt.Println(string(sJson))
	return allData, nil
} //END GetUsrSeller

func EncryptPwd(keyString string, stringToEncrypt string) (encryptedString string) {
	// convert key to bytes
	key, _ := hex.DecodeString(keyString)
	plaintext := []byte(stringToEncrypt)

	//Create a new Cipher Block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	// convert to base64
	return base64.URLEncoding.EncodeToString(ciphertext)
}

// decrypt from base64 to decrypted string
func DecryptPwd(keyString string, stringToDecrypt string) string {
	key, _ := hex.DecodeString(keyString)
	ciphertext, _ := base64.URLEncoding.DecodeString(stringToDecrypt)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	if len(ciphertext) < aes.BlockSize {
		panic("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(ciphertext, ciphertext)

	return fmt.Sprintf("%s", ciphertext)
}
