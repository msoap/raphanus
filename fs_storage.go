package raphanus

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/msoap/raphanus/common"
)

/*
TTL lost after save and load!

File format:
key_name \t type(int|str|list|dict) \t 12|["string"]|[...]|{...}

k1 \t int \t 123 \n
k2 \t str \t ["string value"] \n
k3 \t list \t ["string","value"] \n
k4 \t dict \t {"dk1":"v1","dk2":"v2"} \n

*/

// fsLoad - load cache from file (on start)
func (db *DB) fsLoad() (err error) {
	db.Lock()
	db.withLock = false // for common lock for save all keys
	defer func() {
		db.Unlock()
		db.withLock = true
	}()

	file, err := os.Open(db.fsStorageName)
	if err != nil {
		return err
	}
	defer func() {
		err = file.Close()
	}()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "\t", 3)
		if len(parts) != 3 {
			return fmt.Errorf("Load from file, file is damaged")
		}
		switch parts[1] {
		case "int":
			var intVal int64
			intVal, err = strconv.ParseInt(parts[2], 10, 64)
			if err != nil {
				return fmt.Errorf("Load from file, file is damaged, key: '%s', parse int: %s", parts[0], err)
			}

			err = db.SetInt(parts[0], intVal, 0)
			if err != nil {
				return fmt.Errorf("Load from file, file is damaged, key: '%s'", parts[0])
			}
		case "str":
			listVal := make(raphanuscommon.ListValue, 0)
			err = json.Unmarshal([]byte(parts[2]), &listVal)
			if err != nil {
				return fmt.Errorf("Load from file, file is damaged, key: '%s', parse string: %s", parts[0], err)
			}
			if len(listVal) != 1 {
				return fmt.Errorf("Load from file, file is damaged, key: '%s', parse string", parts[0])
			}

			strVal := listVal[0]
			err = db.SetStr(parts[0], strVal, 0)
			if err != nil {
				return fmt.Errorf("Load from file, file is damaged, key: '%s'", parts[0])
			}
		case "list":
			listVal := make(raphanuscommon.ListValue, 0)
			err = json.Unmarshal([]byte(parts[2]), &listVal)
			if err != nil {
				return fmt.Errorf("Load from file, file is damaged, key: '%s', parse list: %s", parts[0], err)
			}
			db.SetList(parts[0], listVal, 0)
		case "dict":
			dictVal := make(raphanuscommon.DictValue)
			err = json.Unmarshal([]byte(parts[2]), &dictVal)
			if err != nil {
				return fmt.Errorf("Load from file, file is damaged, key: '%s', parse dict: %s", parts[0], err)
			}
			db.SetDict(parts[0], dictVal, 0)
		default:
			return fmt.Errorf("Load from file, file is damaged, key: '%s', unknown type: %s", parts[0], parts[1])
		}
	}

	return scanner.Err()
}

// fsSave - save cache to file
func (db *DB) fsSave() (err error) {
	db.RLock()
	db.withLock = false // for common lock for save all keys
	defer func() {
		db.RUnlock()
		db.withLock = true
	}()

	file, err := os.Create(db.fsStorageName)
	if err != nil {
		return err
	}
	defer func() {
		err = file.Close()
	}()

	for key, rawValue := range db.data {
		row := make([]string, 3)
		row[0] = key
		switch value := rawValue.(type) {
		case int64:
			row[1] = "int"
			row[2] = strconv.FormatInt(value, 10)
		case string:
			row[1] = "str"
			valueAsArr := []string{value}
			byteVal, err := json.Marshal(valueAsArr)
			if err != nil {
				return err
			}
			row[2] = string(byteVal)
		case raphanuscommon.ListValue:
			row[1] = "list"
			byteVal, err := json.Marshal(value)
			if err != nil {
				return err
			}
			row[2] = string(byteVal)
		case raphanuscommon.DictValue:
			row[1] = "dict"
			byteVal, err := json.Marshal(value)
			if err != nil {
				return err
			}
			row[2] = string(byteVal)
		default:
			return fmt.Errorf("Save from file, unknown type for key: '%s'", key)
		}

		line := strings.Join(row, "\t") + "\n"
		_, err := file.WriteString(line)
		if err != nil {
			return fmt.Errorf("Save from file, save row for key: '%s': %s", key, err)
		}
	}

	return nil
}

// fsHandle - handle periodic save on disk
func (db *DB) fsHandle() {
	for {
		time.Sleep(time.Duration(db.fsStorageSyncTime) * time.Second)
		err := db.fsSave()
		if err != nil {
			// TODO: add option for logging
			log.Printf("Save to file failed: %s", err)
		}
	}
}
