# JSONIC

[![GoDoc](https://godoc.org/github.com/sinhashubham95/jsonic?status.svg)](https://pkg.go.dev/github.com/sinhashubham95/jsonic)
[![Release](https://img.shields.io/github/v/release/sinhashubham95/jsonic?sort=semver)](https://github.com/sinhashubham95/jsonic/releases)
[![Report](https://goreportcard.com/badge/github.com/sinhashubham95/jsonic)](https://goreportcard.com/report/github.com/sinhashubham95/jsonic)
[![Coverage Status](https://coveralls.io/repos/github/sinhashubham95/jsonic/badge.svg?branch=master)](https://coveralls.io/github/sinhashubham95/jsonic?branch=master)
[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go#json)

`Jsonic` is the complete set of utilities to handle json data. There's no need to define structs anymore. It's completely safe to perform nested queries in the JSON. The strong typed methods part of this library will help you have the desired result without writing any extra code.

## Installation

```shell
go get github.com/sinhashubham95/jsonic
```

## Understanding the query path

`Jsonic` uses a unique and simple way to query the elements in a json. It's easy but unique, so you need to understand the same for using `Jsonic`.

Consider the following json.

```json
{
  "a": {
    "x": "p",
    "arr": [
      {
        "a": "b",
        "c.d": {
          "e": "f"
        }
      }
    ]
  },
  "c": "d",
  "a.x": {
    "y": "q"
  },
  "a.x.y": {
    "z": "r"
  }
}
```

Though, practically such a JSON won't exist, but still `Jsonic` is intelligent enough to handle even this. Go through the below table carefully and it will help you understand the path schema.

|      Path      |                         Result                         | Comments                                                                  |
| :------------: | :----------------------------------------------------: | ------------------------------------------------------------------------- |
| {EMPTY_STRING} |                      entire json                       | empty string returns the entire json if no empty string exists in the key |
|       .        |                      entire json                       | dot returns the entire json if no dot exists in the key                   |
|       a        | {"x": "p", "arr": [{ "a": "b", "c.d": { "e": "f" } }]} | it returns the entire json tree of a                                      |
|      a.x       |                           p                            | multiple options here, but the first preference goes to the tree of a     |
|     a.x.y      |                           q                            | multiple options here, the first preference will be given to a.x          |
|    a.x.y.z     |                           r                            | there is only a single possibility here                                   |
|     a.arr      |          [{ "a": "b", "c.d": { "e": "f" } }]           | it returns the entire array denoting the json tree of a.arr               |
|    a.arr[0]    |           { "a": "b", "c.d": { "e": "f" } }            | it returns the first element of the array                                 |
|   a.arr[0].a   |                           b                            | it returns the element for key a of the first element of array            |
| a.arr[0].c.d.e |                           f                            |                                                                           |

As you would have understood, if there are multiple JSON trees satisfying the path, and the path looks something like this `a.b.c.d`, then the preferences will be in the following order - `a` > `a.b` > `a.b.c` > `a.b.c.d`.

Consider another json.

```json
{
  "": "a",
  ".": "b"
}
```
Here the paths resolve in a different manner.

|      Path      | Result | Comments                                                                  |
| :------------: | :----: | ------------------------------------------------------------------------- |
| {EMPTY_STRING} |   a    | empty string returns the entire json if no empty string exists in the key |
|       .        |   b    | empty string returns the entire json if no empty string exists in the key |

## How to Use?

`Jsonic` allows you to process the JSON bytes. You can create a new instance of `Jsonic` for every JSON you have and you can get the child JSON trees using the set of utilities it provides.

### Create a New Instance

This will create a new instance using the JSON bytes provided as it's data to be used on.

```go
import (
  "github.com/sinhashubham95/jsonic"
)

func New() {
  json := "{\"naruto\": \"rocks\"}"
  j, err := jsonic.New([]byte(json))
  // perform any sort of operations on the json using the instance created
}
```

### Create a child instance

On the `Jsonic` created, you can provide a child path and get a new instance with the child JSON tree satisfying the path provided as it's data.

```go
import (
  "github.com/sinhashubham95/jsonic"
)

func Child() {
  json := "{\"naruto\": \"rocks\"}"
  j, err := jsonic.New([]byte(json))
  if err != nil {
    return
  }
  
  // create a child
  child, err := jsonic.Child("naruto")
  // now if you want to query on the child then use this child instance
}
```

### Get the data at the path

On the `Jsonic` created, you can get the data at the path specified.

```go
import (
  "github.com/sinhashubham95/jsonic"
)

func Get() {
  json := "{\"naruto\": \"rocks\"}"
  j, err := jsonic.New([]byte(json))
  if err != nil {
    return
  }
  
  // get the data
  data, err := jsonic.Get("naruto")
  // this data will have type interface{} with value "rocks"
}
```

### Get typed data at the path

Though using structs is not required with the wonderful set of utilities `Jsonic` provides, but even if you like to use that, it is very simple to get your result cast into the struct you want.

```go
import (
  "github.com/sinhashubham95/jsonic"
)

type Detail struct {
  Name string `json:"name"`
}

func GetTyped() {
  json := "{\"characters\": [{\"name\": \"naruto\"}, {\"name\": \"boruto\"}]}"
  j, err := jsonic.New([]byte(json))
  if err != nil {
    return
  }
  
  // get the data
  var data []Detail
  err := jsonic.GetTyped("characters", &data)
  // this data will contain 2 elements with names as naruto and boruto
}
```

### Other Typed Utilities

Apart from the generic query methods mentioned above, `Jsonic` contains a bunch of others.

```go
import (
  "github.com/sinhashubham95/jsonic"
)

func OtherGetters(j *Jsonic, path string) {
  // primitives
  i, err := j.GetInt(path)                    // int
  i64, err := j.GetInt64(path)                // int64
  f, err := j.GetFloat(path)                  // float32
  f64, err := j.GetFloat64(path)              // float64
  b, err := j.GetBool(path)                   // bool
  s, err := j.GetString(path)                 // string

  // arrays
  a, err := j.GetArray(path)                  // []interface{}
  iArr, err := j.GetIntArray(path)            // []int
  i64Arr, err := j.GetInt64Array(path)        // []int64
  fArr, err := j.GetFloatArray(path)          // []float32
  f64Arr, err := j.GetFloat64Array(path)      // []float64
  bArr, err := j.GetBoolArray(path)           // []bool
  sArr, err := j.GetStringArray(path)         // []string

  // maps
  m, err := j.GetMap(path)                    // map[string]interface{}
  iMap, err := j.GetIntMap(path)              // map[string]int
  i64Map, err := j.GetInt64Map(path)          // map[string]int64
  fMap, err := j.GetFloatMap(path)            // map[string]float32
  f64Map, err := j.GetFloat64Map(path)        // map[string]float64
  bMap, err := j.GetBoolMap(path)             // map[string]bool
  sMap, err := j.GetStringMap(path)           // map[string]string
}
```
