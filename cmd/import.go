package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"github.com/DoodleScheduling/vault-backend-migrator/vault"
)

func Import(path, file, ver string, mergeImport bool) error {
	abs, err := filepath.Abs(file)
	if err != nil {
		return err
	}

	// Check the input file exists
	if _, err := os.Stat(abs); err != nil {
		f, err := os.Create(abs)
		defer f.Close()
		if err != nil {
			return err
		}
	}

	// Read input file
	b, err := ioutil.ReadFile(abs)
	if err != nil {
		return err
	}

	// Parse data
	var wrap Wrap
	err = json.Unmarshal(b, &wrap)
	if err != nil {
		return err
	}

	// Setup vault client
	v, err := vault.NewClient()
	if v == nil || err != nil {
		if err != nil {
			return err
		}
		return errors.New("Unable to create vault client")
	}

  var wg sync.WaitGroup
	concurent := 50
	ch :=  make(chan Item, concurent)

	for i := 0; i < concurent; i++ {
			go func (){
					for item := range ch {
						data := make(map[string]string)

						for _, kv := range item.Pairs {
							data[kv.Key] = kv.Value
						}

						if mergeImport == true {
							existing := v.Read(item.Path)
							if existing != nil {
								for k, v := range existing {
									data[k] = v.(string)
								}
							}
						}

						fmt.Printf("Writing %s\n", item.Path)
						if err := v.Write(item.Path, data, ver); err != nil {
							fmt.Printf("Error %s\n", err)
						}

						wg.Done()
					}
			}()
	}

	// Write each keypair to vault
	for _, item := range wrap.Data {
		wg.Add(1)
		ch <- item
	}

	wg.Wait()
	return nil
}
