package database

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"database/sql"

	"comprarmas.com.mx/internal/secrets"

	_ "github.com/go-sql-driver/mysql"

	"github.com/go-redis/redis"
)

func CnnMongo() (cClient *mongo.Client, e error) {
	sCnnDB := stringCnnMongo()
	//sCnnDB := "mongodb://127.0.0.1:27017/"
	//sCnnDB := "mongodb://docker:mongopw@172.24.100.14:27017"
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(sCnnDB))
	fmt.Println(client)
	if err != nil {
		log.Println("CnnMongo_Err: " + err.Error())
		panic(err)
	} else {
		fmt.Println("Conexion realizada")

		err = client.Ping(context.TODO(), nil)
		if err != nil {
			log.Println("CnnMongo_ping_Err: " + err.Error())
			return nil, err
		}

		return client, err
	}
} //FIN de cnnMongo

//=====================  Archivo Conf ==========================================
func stringCnnMongo() string {
	sIP := secrets.LoadSecrets("MONGO_HOST")
	sPort := secrets.LoadSecrets("MONGO_PORT")
	sUsr := secrets.LoadSecrets("MONGO_USER")
	sPWD := secrets.LoadSecrets("MONGO_PWD")

	//mongodb://127.0.0.1:27017/
	//mongodb://docker:mongopw@172.24.100.14:27017
	//sCnn := "mongodb://" + sIp + ":" + fmt.Sprintf(sPort) + "/"
	sCnn := "mongodb://" + sUsr + ":" + sPWD + "@" + sIP + ":" + fmt.Sprintf(sPort)

	return sCnn
} //FIN stringCnnMongo

func Close_ClientCtx(client *mongo.Client, ctx context.Context) {
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
}

func Close(client *mongo.Client, ctx context.Context, cancel context.CancelFunc) {
	defer cancel()

	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
}

//cnn MYSQL
func CnnDB_ENV() (db *sql.DB, e error) {
	sPwd := ""
	sIP := ""

	sIP = secrets.LoadSecrets("DB_HOST")
	sPort := secrets.LoadSecrets("DB_PORT")
	sUsr := secrets.LoadSecrets("DB_USER")
	sTmp := secrets.LoadSecrets("DB_PWD2")
	if sTmp != "" {
		sPwd = secrets.LoadSecrets("DB_PWD1") + "#" + secrets.LoadSecrets("DB_PWD2")
	} else {
		sPwd = secrets.LoadSecrets("DB_PWD1")
	}
	sDB := secrets.LoadSecrets("DB_NAME")

	db, err := sql.Open("mysql", sUsr+":"+sPwd+"@tcp("+sIP+":"+sPort+")/"+sDB+"?parseTime=true")
	if err != nil {
		//log.Println(err.Error())
		return nil, err
	}
	return db, nil
}
//=========================================================
//=========================================================

func RedisCnn_ENV() *redis.Client {
	sIP := secrets.LoadSecrets("REDIS_HOST")
	sPort := secrets.LoadSecrets("REDIS_PORT")
	//sUsr := utils.LoadSecrets("REDIS_USER")
	sPwd := secrets.LoadSecrets("REDIS_PWD")
	//sDB := utils.LoadSecrets("DB_NAME")

	client := redis.NewClient(&redis.Options{
		Addr:     sIP + ":" + sPort,
		Password: sPwd,
		DB:       10,
	})
	return client
}

func Ping(client *redis.Client) error {
	pong, err := client.Ping().Result()
	if err != nil {
		return err
	}
	fmt.Println(pong, err)
	// Output: PONG <nil>
	return nil
}

func Set(client *redis.Client, sKey string, sValue string) error {
	err := client.Set(sKey, sValue, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func Get(client *redis.Client, sKey string) (blnOK bool, err error) {
	nameVal, err := client.Get(sKey).Result()
	if nameVal != "" {
		return true, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func AddExpTime(client *redis.Client, sKey string) (blnOK bool, err error) {
	nT := secrets.LoadSecrets("EXP_TIME")
	nH, _ := strconv.ParseInt(nT, 10, 64)
	nTmp := int32(nH)
	//yourTime := Int31n(nTmp)
	_, err = client.Expire(sKey, time.Duration(nTmp)*time.Minute).Result()
	if err != nil {
		return false, err
	}
	return true, nil
}
