package framework

import (
	"fmt"
	"icecream/utils"
	"reflect"
	"strings"
)

type ModelInterface interface {
	GetTable() string
	SetTable(string)
	Find() ModelInterface
	Set(ModelInterface, string, interface{})
	Get(ModelInterface, string) interface{}
	Save(ModelInterface, map[string]interface{}) int64
}

type Model struct {
	tableName string
	ID        int64
}

func (model *Model) GetTable() string {
	return model.tableName
}

func (model *Model) SetTable(table string) {
	model.tableName = table
}

func (model *Model) Set(m ModelInterface, key string, value interface{}) {
	reflect.ValueOf(m).Elem().FieldByName(key).Set(reflect.ValueOf(value))
}

func (model Model) Get(m ModelInterface, key string) interface{} {
	return reflect.ValueOf(m).Elem().FieldByName(key).Interface()
}

func (model *Model) Find() ModelInterface {
	return model
}

func (model *Model) Save(m ModelInterface, data map[string]interface{}) int64 {
	if _, ok := data["id"]; ok {
		model.ID = int64(data["id"].(int))
	}
	//db, _ := GetDB()

	var keys []string
	var values []string
	var args []interface{}
	for key, val := range data {
		keys = append(keys, key)
		values = append(values, "?")
		args = append(args, val)
	}

	var sql string
	if model.ID != 0 {
		args = append(args, model.ID)
		sql = "UPDATE " + m.GetTable() + " SET " + strings.Join(keys, "=?, ") + "=?" + " WHERE id=?"
		_, err := DB.Exec(sql, args...)
		//execArgs := append([]interface{}{sql}, args...)
		//_, err := db.Exec(execArgs...)
		if err != nil {
			fmt.Printf("UPDATE error: %s\n", err.Error())
			return 0
		}
		//db.Table(model.tableName).Where("id=?", model.ID).Updates(data)
	} else {
		sql = "INSERT INTO " + m.GetTable() + " (" + strings.Join(keys, ",") + ") VALUES (" + strings.Join(values, ",") + ")"
		result, err := DB.Exec(sql, args...)
		//execArgs := append([]interface{}{sql}, args...)
		//result, err := db.Exec(execArgs...)
		if err != nil {
			fmt.Printf("Insert error: %s\n", err.Error())
			return 0
		}
		model.ID, err = result.LastInsertId()
		if err != nil {
			fmt.Printf("Insert error: %s\n", err.Error())
			return 0
		}
		//db.Table(model.tableName).Create()
	}
	fmt.Printf("SQL: %s\n, ARGS: %v\n", sql, args)
	//model.SetAll(m, data)
	return model.ID
}

func (model *Model) SetAll(m ModelInterface, data map[string]interface{}) {
	fmt.Printf("M: %+v\n", m)
	//t := reflect.TypeOf(m)
	//fmt.Printf("T: %+v\n", t)
	v := reflect.ValueOf(m).Elem()
	t := reflect.Indirect(v).Type()
	fmt.Printf("V: %+v\n", v)
	fmt.Printf("T2: %+v\n", t)
	fmt.Printf("model field num: %d\n", v.NumField())
	for i := 0; i < v.NumField(); i++ {
		if t.Field(i).Anonymous && t.Field(i).Name == "Model" {
			subt := reflect.Indirect(v.Field(i)).Type()
			for j := 0; j < v.Field(i).NumField(); j++ {
				fmt.Printf("T field name: %v\n", subt.Field(j).Name)
			}
		} else {
			fmt.Printf("T field name: %s\n", t.Field(i).Name)
			_ = utils.ToCamelCase("model_name")
		}
		//if _, ok := data[key]; ok {
		//	model.ID = int64(data[key].(int))
		//	v.Field(i).Set(reflect.ValueOf(data[key]))
		//}
	}
}

type ModelStructs struct {
}
