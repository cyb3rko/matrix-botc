/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package botc

import "fmt"

func format(string string, a ...interface{}) string {
	return fmt.Sprintf(string, a)
}

//goland:noinspection GoReservedWordUsedAsName
func print(message string) {
	fmt.Println(message)
}

func printf(message string, a ...interface{}) {
	fmt.Printf(message+"\n", a...)
}
