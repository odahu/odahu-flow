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
package utils_test

import (
	"github.com/odahu/odahu-flow/packages/operator/pkg/utils"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

const (
	archiveName   = "archive.zip"
	tempDirPrefix = "zip-test"
	file1Name     = "file1.txt"
	nestedDirName = "dir"
	file2Name     = "file2.txt"
)

type ZipTestSuite struct {
	suite.Suite
	g           *GomegaWithT
	workDirPath string
}

func (s *ZipTestSuite) SetupTest() {
	s.g = NewGomegaWithT(s.T())

	// Create the following directory structure
	// ../temp_dir/
	// ../temp_dir/file1.txt
	// ../temp_dir/dir/
	// ../temp_dir/dir/file2.txt
	tempDir, err := ioutil.TempDir("", tempDirPrefix)
	s.g.Expect(err).Should(BeNil())
	s.workDirPath = tempDir

	err = ioutil.WriteFile(filepath.Join(tempDir, file1Name), []byte(file1Name), os.ModePerm)
	s.g.Expect(err).Should(BeNil())

	nestedDirPath := filepath.Join(tempDir, nestedDirName)
	err = os.Mkdir(nestedDirPath, os.ModePerm)
	s.g.Expect(err).Should(BeNil())

	err = ioutil.WriteFile(filepath.Join(nestedDirPath, file2Name), []byte(file2Name), os.ModePerm)
	s.g.Expect(err).Should(BeNil())
}

func (s *ZipTestSuite) TearDownTest() {
	if len(s.workDirPath) != 0 {
		err := os.RemoveAll(s.workDirPath)
		s.g.Expect(err).Should(BeNil())
	}
}

func TestZipTestSuite(t *testing.T) {
	suite.Run(t, new(ZipTestSuite))
}

func (s *ZipTestSuite) TestMainZipWorkflow() {
	err := utils.ZipDir(s.workDirPath, archiveName)
	s.g.Expect(err).Should(BeNil())

	outputDir, err := ioutil.TempDir("", tempDirPrefix)
	s.g.Expect(err).Should(BeNil())
	defer os.RemoveAll(outputDir)

	err = utils.Unzip(archiveName, outputDir)
	s.g.Expect(err).Should(BeNil())

	file1Content, err := ioutil.ReadFile(filepath.Join(outputDir, file1Name))
	s.g.Expect(err).Should(BeNil())
	s.g.Expect(file1Content).Should(Equal([]byte(file1Name)))

	file2Content, err := ioutil.ReadFile(filepath.Join(outputDir, nestedDirName, file2Name))
	s.g.Expect(err).Should(BeNil())
	s.g.Expect(file2Content).Should(Equal([]byte(file2Name)))
}
