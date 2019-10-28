package utils

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
)

// 数据库查询多行数据，返回一个map键值对数组
// db: 传入db连接实例指针, query: 传入查询语句, args: 传入查询参数
func DBQuery(db *sql.DB, query string, args ...interface{}) ([]map[string]interface{}, error) {
	results := make([]map[string]interface{}, 0)
	rows, err := db.Query(query, args...)
	if err != nil {
		return results, err
	}
	for rows.Next() {
		result := make(map[string]interface{})
		err = MapScan(rows, result)
		if err != nil {
			return results, err
		}
		for key, value := range result {
			result[key] = getTypeVal(value)
		}
		results = append(results, result)
	}
	return results, err
}

// 数据库查询单行数据，返回一个map键值对
// db: 传入db连接实例指针, query: 传入查询语句, args: 传入查询参数
func DBQueryRow(db *sql.DB, query string, args ...interface{}) (map[string]interface{}, error) {
	result := make(map[string]interface{}, 1)
	rows, err := db.Query(query, args...)
	if err != nil {
		return result, err
	}
	err = MapScan(rows, result)
	if err != nil {
		return result, err
	}
	for key, value := range result {
		result[key] = getTypeVal(value)
	}
	return result, err
}

// 数据库查询多行数据，返回一个map键值对数组
// db: 传入sqlx.db连接实例指针, query: 传入查询语句, args: 传入查询参数
func XDBQuery(db *sqlx.DB, query string, args ...interface{}) ([]map[string]interface{}, error) {
	results := make([]map[string]interface{}, 0)
	rows, err := db.Queryx(query, args...)
	if err != nil {
		return results, err
	}
	for rows.Next() {
		result := make(map[string]interface{})
		err = rows.MapScan(result)
		if err != nil {
			return results, err
		}
		for key, value := range result {
			result[key] = getTypeVal(value)
		}
		results = append(results, result)
	}
	return results, err
}

// 数据库查询单行数据，返回一个map键值对
// db: 传入sqlx.db连接实例指针, query: 传入查询语句, args: 传入查询参数
func XDBQueryRow(db *sqlx.DB, query string, args ...interface{}) (map[string]interface{}, error) {
	result := make(map[string]interface{}, 1)
	rows, err := db.Queryx(query, args...)
	if err != nil {
		return result, err
	}
	err = rows.MapScan(result)
	if err != nil {
		return result, err
	}
	for key, value := range result {
		result[key] = getTypeVal(value)
	}
	return result, err
}

// MapScan将sql.Rows当前行数据赋值给dest
func MapScan(r *sql.Rows, dest map[string]interface{}) error {
	// ignore r.started, since we needn't use reflect for anything.
	columns, err := r.Columns()
	if err != nil {
		return err
	}

	values := make([]interface{}, len(columns))
	for i := range values {
		values[i] = new(interface{})
	}

	err = r.Scan(values...)
	if err != nil {
		return err
	}

	for i, column := range columns {
		dest[column] = *(values[i].(*interface{}))
	}

	return r.Err()
}

func getTypeVal(value interface{}) interface{} {
	switch value.(type) { //v表示b1 接口转换成Bag对象的值
	case []byte:
		value = string(value.([]byte))
		break
	//case nil:
	//    value = nil
	//    break
	//case string:
	//    break
	//case int8:
	//    break
	//case int16:
	//    break
	//case int32:
	//    break
	//case int64:
	//    break
	//case uint8:
	//    break
	//case uint16:
	//    break
	//case uint32:
	//    break
	//case uint64:
	//    break
	//case float32:
	//    break
	//case float64:
	//    break
	default:
		break
	}
	return value
}
