// Code generated by go-bindata. (@generated) DO NOT EDIT.
// sources:
// schema/config.yaml
// schema/metaplan.yaml
// schema/override.yaml
// schema/plan.yaml
// schema/registry-ecr.yaml
// schema/registry-gcr.yaml
// schema/registry-server.yaml
// schema/task-build.yaml
// schema/task-plan.yaml
// schema/task-push.yaml
// schema/task-remove.yaml
// schema/task-run.yaml
package asset

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)
type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _schemaConfigYaml = []byte(`---

definitions:
  stringOrList:
    oneOf:
      - type: string
      - type: array
        items:
          type: string
  importFileList:
    type: object
    properties:
      files:
        $ref: '#/definitions/stringOrList'
      excludes:
        $ref: '#/definitions/stringOrList'
    additionalProperties: false
  exportFileList:
    type: object
    properties:
      files:
        $ref: '#/definitions/stringOrList'
      excludes:
        $ref: '#/definitions/stringOrList'
      clean-excludes:
        $ref: '#/definitions/stringOrList'
    additionalProperties: false

type: object
properties:
  extends:
    $ref: '#/definitions/stringOrList'
  options:
    type: object
    properties:
      force-sequential:
        type: boolean
      healthcheck-interval:
        type: string
      ssh-identities:
        oneOf:
          - type: string
          - type: array
            items:
              type: string
    additionalProperties: false
  registries:
    type: array
    items:
      type: object
  workspace:
    type: string
  environment:
    $ref: '#/definitions/stringOrList'
  env-file:
    $ref: '#/definitions/stringOrList'
  import:
    $ref: '#/definitions/importFileList'
  export:
    $ref: '#/definitions/exportFileList'
  tasks:
    type: object
  plans:
    type: object
  metaplans:
    type: object
additionalProperties: false
`)

func schemaConfigYamlBytes() ([]byte, error) {
	return _schemaConfigYaml, nil
}

func schemaConfigYaml() (*asset, error) {
	bytes, err := schemaConfigYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "schema/config.yaml", size: 1371, mode: os.FileMode(420), modTime: time.Unix(1546566187, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _schemaMetaplanYaml = []byte(`---

type: array
items:
  type: string
`)

func schemaMetaplanYamlBytes() ([]byte, error) {
	return _schemaMetaplanYaml, nil
}

func schemaMetaplanYaml() (*asset, error) {
	bytes, err := schemaMetaplanYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "schema/metaplan.yaml", size: 39, mode: os.FileMode(420), modTime: time.Unix(1538014972, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _schemaOverrideYaml = []byte(`---

definitions:
  stringOrList:
    oneOf:
      - type: string
      - type: array
        items:
          type: string
  importFileList:
    type: object
    properties:
      excludes:
        $ref: '#/definitions/stringOrList'
    additionalProperties: false
  exportFileList:
    type: object
    properties:
      excludes:
        $ref: '#/definitions/stringOrList'
      clean-excludes:
        $ref: '#/definitions/stringOrList'
    additionalProperties: false

type: object
properties:
  options:
    type: object
    properties:
      ssh-identities:
        $ref: '#/definitions/stringOrList'
      force-sequential:
        type: boolean
      healthcheck-interval:
        type: string
      path-substitutions:
        type: object
        additionalProperties:
          type: string
    additionalProperties: false
  registries:
    type: array
    items:
      type: object
  environment:
    $ref: '#/definitions/stringOrList'
  env-file:
    $ref: '#/definitions/stringOrList'
  import:
    $ref: '#/definitions/importFileList'
  export:
    $ref: '#/definitions/exportFileList'
additionalProperties: false
`)

func schemaOverrideYamlBytes() ([]byte, error) {
	return _schemaOverrideYaml, nil
}

func schemaOverrideYaml() (*asset, error) {
	bytes, err := schemaOverrideYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "schema/override.yaml", size: 1130, mode: os.FileMode(420), modTime: time.Unix(1546907405, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _schemaPlanYaml = []byte(`---

definitions:
  stringOrList:
    oneOf:
      - type: string
      - type: array
        items:
          type: string
  stage:
    type: object
    properties:
      name:
        type: string
      disabled:
        type: string
      before-stage:
        type: string
      after-stage:
        type: string
      tasks:
        type: array
        items:
          $ref: '#/definitions/stageTask'
      run-mode:
        type: string
        enum:
          - always
          - on-success
          - on-failure
      parallel:
        type: boolean
      environment:
        $ref: '#/definitions/stringOrList'
    additionalProperties: false
    required:
      - name
  stageTask:
    oneOf:
      - type: string
      - type: object
        properties:
          name:
            type: string
          disabled:
            type: string
          environment:
            $ref: '#/definitions/stringOrList'
        additionalProperties: false
        required:
          - name

type: object
properties:
  extends:
    type: string
  disabled:
    type: string
  stages:
    type: array
    items:
      $ref: '#/definitions/stage'
  environment:
    $ref: '#/definitions/stringOrList'
additionalProperties: false
`)

func schemaPlanYamlBytes() ([]byte, error) {
	return _schemaPlanYaml, nil
}

func schemaPlanYaml() (*asset, error) {
	bytes, err := schemaPlanYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "schema/plan.yaml", size: 1231, mode: os.FileMode(420), modTime: time.Unix(1546566187, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _schemaRegistryEcrYaml = []byte(`---

type: object
properties:
  type:
    type: string
    enum:
      - ecr
  access-key-id:
    type: string
  secret-access-key:
    type: string
  account-id:
    type: string
  region:
    type: string
  role:
    type: string
additionalProperties: false
`)

func schemaRegistryEcrYamlBytes() ([]byte, error) {
	return _schemaRegistryEcrYaml, nil
}

func schemaRegistryEcrYaml() (*asset, error) {
	bytes, err := schemaRegistryEcrYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "schema/registry-ecr.yaml", size: 260, mode: os.FileMode(420), modTime: time.Unix(1546566187, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _schemaRegistryGcrYaml = []byte(`---

type: object
properties:
  type:
    type: string
    enum:
      - gcr
  hostname:
    type: string
    enum:
      - gcr.io
      - us.gcr.io
      - eu.gcr.io
      - asia.gcr.io
  key:
    type: string
  key-file:
    type: string
additionalProperties: false
`)

func schemaRegistryGcrYamlBytes() ([]byte, error) {
	return _schemaRegistryGcrYaml, nil
}

func schemaRegistryGcrYaml() (*asset, error) {
	bytes, err := schemaRegistryGcrYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "schema/registry-gcr.yaml", size: 268, mode: os.FileMode(420), modTime: time.Unix(1547916565, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _schemaRegistryServerYaml = []byte(`---

type: object
properties:
  type:
    type: string
    enum:
      - server
  server:
    type: string
  username:
    type: string
  password:
    type: string
  password-file:
    type: string
additionalProperties: false
`)

func schemaRegistryServerYamlBytes() ([]byte, error) {
	return _schemaRegistryServerYaml, nil
}

func schemaRegistryServerYaml() (*asset, error) {
	bytes, err := schemaRegistryServerYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "schema/registry-server.yaml", size: 227, mode: os.FileMode(420), modTime: time.Unix(1546566188, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _schemaTaskBuildYaml = []byte(`---

definitions:
  stringOrList:
    oneOf:
      - type: string
      - type: array
        items:
          type: string

type: object
properties:
  type:
    type: string
    enum:
      - build
  extends:
    type: string
  environment:
    $ref: '#/definitions/stringOrList'
  required-environment:
    type: array
    items:
      type: string
  dockerfile:
    type: string
  target:
    type: string
  tags:
    $ref: '#/definitions/stringOrList'
  labels:
    $ref: '#/definitions/stringOrList'
additionalProperties: false
`)

func schemaTaskBuildYamlBytes() ([]byte, error) {
	return _schemaTaskBuildYaml, nil
}

func schemaTaskBuildYaml() (*asset, error) {
	bytes, err := schemaTaskBuildYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "schema/task-build.yaml", size: 533, mode: os.FileMode(420), modTime: time.Unix(1546566188, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _schemaTaskPlanYaml = []byte(`---

definitions:
  stringOrList:
    oneOf:
      - type: string
      - type: array
        items:
          type: string

type: object
properties:
  type:
    type: string
    enum:
      - plan
  extends:
    type: string
  environment:
    $ref: '#/definitions/stringOrList'
  required-environment:
    type: array
    items:
      type: string
  name:
    type: string
additionalProperties: false
`)

func schemaTaskPlanYamlBytes() ([]byte, error) {
	return _schemaTaskPlanYaml, nil
}

func schemaTaskPlanYaml() (*asset, error) {
	bytes, err := schemaTaskPlanYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "schema/task-plan.yaml", size: 403, mode: os.FileMode(420), modTime: time.Unix(1546566188, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _schemaTaskPushYaml = []byte(`---

definitions:
  stringOrList:
    oneOf:
      - type: string
      - type: array
        items:
          type: string

type: object
properties:
  type:
    type: string
    enum:
      - push
  extends:
    type: string
  environment:
    $ref: '#/definitions/stringOrList'
  required-environment:
    type: array
    items:
      type: string
  images:
    $ref: '#/definitions/stringOrList'
  include-built:
    type: boolean
additionalProperties: false
`)

func schemaTaskPushYamlBytes() ([]byte, error) {
	return _schemaTaskPushYaml, nil
}

func schemaTaskPushYaml() (*asset, error) {
	bytes, err := schemaTaskPushYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "schema/task-push.yaml", size: 462, mode: os.FileMode(420), modTime: time.Unix(1550198041, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _schemaTaskRemoveYaml = []byte(`---

definitions:
  stringOrList:
    oneOf:
      - type: string
      - type: array
        items:
          type: string

type: object
properties:
  type:
    type: string
    enum:
      - remove
  extends:
    type: string
  environment:
    $ref: '#/definitions/stringOrList'
  required-environment:
    type: array
    items:
      type: string
  images:
    $ref: '#/definitions/stringOrList'
  include-built:
    type: boolean
additionalProperties: false
`)

func schemaTaskRemoveYamlBytes() ([]byte, error) {
	return _schemaTaskRemoveYaml, nil
}

func schemaTaskRemoveYaml() (*asset, error) {
	bytes, err := schemaTaskRemoveYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "schema/task-remove.yaml", size: 464, mode: os.FileMode(420), modTime: time.Unix(1550198039, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _schemaTaskRunYaml = []byte(`---

definitions:
  stringOrList:
    oneOf:
      - type: string
      - type: array
        items:
          type: string
  healthcheck:
    type: object
    properties:
      command:
        type: string
      interval:
        type: string
      retries:
        type: integer
      start-period:
        type: string
      timeout:
        type: string
    additionalProperties: false

type: object
properties:
  type:
    type: string
    enum:
      - run
  extends:
    type: string
  environment:
    $ref: '#/definitions/stringOrList'
  required-environment:
    type: array
    items:
      type: string
  image:
    type: string
  command:
    type: string
  shell:
    type: string
  script:
    type: string
  entrypoint:
    type: string
  user:
    type: string
  workspace:
    type: string
  hostname:
    type: string
  detach:
    type: boolean
  healthcheck:
    $ref: '#/definitions/healthcheck'
  export-environment-file:
    $ref: '#/definitions/stringOrList'
additionalProperties: false
`)

func schemaTaskRunYamlBytes() ([]byte, error) {
	return _schemaTaskRunYaml, nil
}

func schemaTaskRunYaml() (*asset, error) {
	bytes, err := schemaTaskRunYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "schema/task-run.yaml", size: 1013, mode: os.FileMode(420), modTime: time.Unix(1546566188, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"schema/config.yaml":          schemaConfigYaml,
	"schema/metaplan.yaml":        schemaMetaplanYaml,
	"schema/override.yaml":        schemaOverrideYaml,
	"schema/plan.yaml":            schemaPlanYaml,
	"schema/registry-ecr.yaml":    schemaRegistryEcrYaml,
	"schema/registry-gcr.yaml":    schemaRegistryGcrYaml,
	"schema/registry-server.yaml": schemaRegistryServerYaml,
	"schema/task-build.yaml":      schemaTaskBuildYaml,
	"schema/task-plan.yaml":       schemaTaskPlanYaml,
	"schema/task-push.yaml":       schemaTaskPushYaml,
	"schema/task-remove.yaml":     schemaTaskRemoveYaml,
	"schema/task-run.yaml":        schemaTaskRunYaml,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}

var _bintree = &bintree{nil, map[string]*bintree{
	"schema": &bintree{nil, map[string]*bintree{
		"config.yaml":          &bintree{schemaConfigYaml, map[string]*bintree{}},
		"metaplan.yaml":        &bintree{schemaMetaplanYaml, map[string]*bintree{}},
		"override.yaml":        &bintree{schemaOverrideYaml, map[string]*bintree{}},
		"plan.yaml":            &bintree{schemaPlanYaml, map[string]*bintree{}},
		"registry-ecr.yaml":    &bintree{schemaRegistryEcrYaml, map[string]*bintree{}},
		"registry-gcr.yaml":    &bintree{schemaRegistryGcrYaml, map[string]*bintree{}},
		"registry-server.yaml": &bintree{schemaRegistryServerYaml, map[string]*bintree{}},
		"task-build.yaml":      &bintree{schemaTaskBuildYaml, map[string]*bintree{}},
		"task-plan.yaml":       &bintree{schemaTaskPlanYaml, map[string]*bintree{}},
		"task-push.yaml":       &bintree{schemaTaskPushYaml, map[string]*bintree{}},
		"task-remove.yaml":     &bintree{schemaTaskRemoveYaml, map[string]*bintree{}},
		"task-run.yaml":        &bintree{schemaTaskRunYaml, map[string]*bintree{}},
	}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}
