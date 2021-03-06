//
//    Copyright 2020 EPAM Systems
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.
//

package config

type Claims struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UserConfig struct {
	Claims Claims `json:"claims"`
	// The sign out endpoint logs out the authenticated user.
	SignOutURL string `json:"signOutUrl"`
}

func NewDefaultUserConfig() UserConfig {
	return UserConfig{
		Claims: Claims{
			Name:  "name",
			Email: "email",
		},
	}
}
