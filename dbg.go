// Copyright 2018 The go-ego Project Developers.
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package tt

import "log"

func fmtDbg(s string, v ...interface{}) (arr []interface{}) {
	info := Red("Trace: " + CallerInfo()[3] + ", ")
	info += Blue(s)

	arr = append(arr, info)
	arr = append(arr, v...)

	return
}

// Log dbg log
func Log(s string, v ...interface{}) error {
	arr := fmtDbg(s, v...)
	log.Println(arr...)
	return nil
}

// Err dbg error log
func Err(s string, v ...interface{}) {
	err := fmtDbg(s, v...)
	log.Fatalln(err...)
}