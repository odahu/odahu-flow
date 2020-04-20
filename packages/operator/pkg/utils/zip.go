//
//    Copyright 2019 EPAM Systems
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

package utils

import (
	"os"
	"os/exec"
	"path/filepath"
)

func ZipDir(source, target string) error {
	target, err := filepath.Abs(target)
	if err != nil {
		return err
	}

	if _, err := os.Stat(target); os.IsNotExist(err) {
		emptyFile, err := os.Create(target)
		if err != nil {
			return err
		}

		if err := emptyFile.Close(); err != nil {
			return err
		}
	}

	cmd := exec.Command(
		"tar", "--exclude", target,
		"-cv", "--use-compress-program=pigz", "-f", target, ".",
	)
	cmd.Dir = source
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func Unzip(src string, dest string) error {
	src, err := filepath.Abs(src)
	if err != nil {
		return err
	}

	cmd := exec.Command("tar", "-xvf", src, "-C", ".")
	cmd.Dir = dest
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
