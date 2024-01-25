/*
Copyright © 2024 Abdi Omar martelluiz125@gmail.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"fmt"

	"github.com/Abdiooa/CLSDAPP/cmd/clsdapp"
)

func main() {
	// Create the config file and store KEKkey
	if err := createConfigFile(); err != nil {
		fmt.Println("Error creating config file:", err)
		return
	}

	// Execute the CLSDAPP command
	clsdapp.Execute()
}

func createConfigFile() error {
	// Attempt to create the config file and CLSD folder
	return clsdapp.CreateConfigFile()
}
