/**
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package target

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"mynewt.apache.org/newt/newt/pkg"
	"mynewt.apache.org/newt/newt/project"
	"mynewt.apache.org/newt/util"
)

const TARGET_FILE_NAME string = "target.yml"

var globalTargetMap map[string]*Target

type Target struct {
	basePkg *pkg.LocalPackage

	// XXX: Probably don't need the below four fields; they can just be
	// retrieved from the viper object.  Keep them here for now for easy
	// initializiation of dummy targets.
	CompilerName string
	BspName      string
	AppName      string
	Arch         string

	// target.yml configuration structure
	Vars map[string]string
}

func NewTarget(basePkg *pkg.LocalPackage) *Target {
	target := &Target{}
	target.Init(basePkg)
	return target
}

func LoadTarget(basePkg *pkg.LocalPackage) (*Target, error) {
	target := NewTarget(basePkg)
	if err := target.Load(basePkg); err != nil {
		return nil, err
	}

	return target, nil
}

func (target *Target) Init(basePkg *pkg.LocalPackage) {
	target.basePkg = basePkg
}

func (target *Target) Load(basePkg *pkg.LocalPackage) error {
	v, err := util.ReadConfig(basePkg.BasePath(),
		strings.TrimSuffix(TARGET_FILE_NAME, ".yml"))
	if err != nil {
		return err
	}

	target.Vars = map[string]string{}

	settings := v.AllSettings()
	for k, v := range settings {
		target.Vars[k] = v.(string)
	}

	target.CompilerName = target.Vars["target.compiler"]
	target.BspName = target.Vars["target.bsp"]
	target.AppName = target.Vars["target.app"]
	target.Arch = target.Vars["target.arch"]

	// XXX: Verify required fields set?

	return nil
}

func (target *Target) Package() *pkg.LocalPackage {
	return target.basePkg
}

func (target *Target) Name() string {
	return target.basePkg.Name()
}

func (target *Target) ShortName() string {
	return filepath.Base(target.Name())
}

func (target *Target) Compiler() *pkg.LocalPackage {
	dep, _ := pkg.NewDependency(nil, target.CompilerName)
	mypkg := project.GetProject().ResolveDependency(dep).(*pkg.LocalPackage)

	return mypkg
}

func (target *Target) App() *pkg.LocalPackage {
	dep, err := pkg.NewDependency(nil, target.AppName)
	if err != nil {
		return nil
	}

	appPkg := project.GetProject().ResolveDependency(dep)
	if appPkg == nil {
		return nil
	}

	return appPkg.(*pkg.LocalPackage)
}

func (target *Target) Bsp() *pkg.LocalPackage {
	dep, _ := pkg.NewDependency(nil, target.BspName)
	mypkg := project.GetProject().ResolveDependency(dep).(*pkg.LocalPackage)
	return mypkg
}

func (target *Target) ElfPath() string {
	appPkg := target.App()
	if appPkg == nil {
		return ""
	}

	return appPkg.BasePath() + "/bin/" + target.Package().Name() + "/" +
		appPkg.Name() + ".elf"
}

func (target *Target) ImagePath() string {
	appPkg := target.App()
	if appPkg == nil {
		return ""
	}

	return appPkg.BasePath() + "/bin/" + target.Package().Name() + "/" +
		appPkg.Name() + ".img"
}

// Save the target's configuration elements
func (t *Target) Save() error {
	if err := t.basePkg.Save(); err != nil {
		return err
	}

	dirpath := t.basePkg.BasePath()
	filepath := dirpath + "/" + TARGET_FILE_NAME
	file, err := os.Create(filepath)
	if err != nil {
		return util.NewNewtError(err.Error())
	}
	defer file.Close()

	file.WriteString("### Target: " + t.basePkg.Name() + "\n\n")

	keys := []string{}
	for k, _ := range t.Vars {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		file.WriteString(k + ": " + t.Vars[k] + "\n")
	}

	return nil
}

// Tells you if the target's directory contains extra user files (i.e., files
// other than pkg.yml).
func (t *Target) ContainsUserFiles() (bool, error) {
	contents, err := ioutil.ReadDir(t.basePkg.BasePath())
	if err != nil {
		return false, err
	}

	userFiles := false
	for _, node := range contents {
		name := node.Name()
		if name != "." && name != ".." &&
			name != pkg.PACKAGE_FILE_NAME && name != TARGET_FILE_NAME {

			userFiles = true
			break
		}
	}

	return userFiles, nil
}

func (t *Target) Delete() error {
	if err := os.RemoveAll(t.basePkg.BasePath()); err != nil {
		return util.NewNewtError(err.Error())
	}

	return nil
}

func buildTargetMap() error {
	globalTargetMap = map[string]*Target{}

	packs := project.GetProject().PackageList()
	for _, packHash := range packs {
		for name, pack := range *packHash {
			if pack.Type() == pkg.PACKAGE_TYPE_TARGET {
				target, err := LoadTarget(pack.(*pkg.LocalPackage))
				if err != nil {
					return err
				}
				globalTargetMap[name] = target
			}
		}
	}

	return nil
}

func GetTargets() map[string]*Target {
	if globalTargetMap == nil {
		err := buildTargetMap()
		if err != nil {
			panic(err.Error())
		}
	}

	return globalTargetMap
}