// Copyright 2018 The Nakama Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
)

const codeTemplate string = `// tslint:disable
/* Code generated by openapi-gen/main.go. DO NOT EDIT. */

const BASE_PATH = "http://127.0.0.1:80";

export interface ConfigurationParameters {
  basePath?: string;
  username?: string;
  password?: string;
  bearerToken?: string;
  timeoutMs?: number;
}

{{- range $classname, $definition := .Definitions}}
/** {{$definition.Description}} */
export interface {{$classname | title}} {
  {{- range $fieldname, $property := $definition.Properties}}
  // {{$property.Description}}
  {{- if eq $property.Type "integer"}}
  {{$fieldname}}?: number;
  {{- else if eq $property.Type "number" }}
  {{$fieldname}}?: number;
  {{- else if eq $property.Type "boolean"}}
  {{$fieldname}}?: boolean;
  {{- else if eq $property.Type "array"}}
    {{- if eq $property.Items.Type "string"}}
  {{$fieldname}}?: Array<string>;
    {{- else if eq $property.Items.Type "integer"}}
  {{$fieldname}}?: Array<number>;
    {{- else if eq $property.Items.Type "boolean"}}
  {{$fieldname}}?: Array<boolean>;
    {{- else}}
  {{$fieldname}}?: Array<{{$property.Items.Ref | cleanRef}}>;
    {{- end}}
  {{- else if eq $property.Type "object"}}
    {{- if eq $property.AdditionalProperties.Type "string"}}
  {{$fieldname}}?: Map<string, string>;
    {{- else if eq $property.AdditionalProperties.Type "integer"}}
  {{$fieldname}}?: Map<string, integer>;
    {{- else if eq $property.AdditionalProperties.Type "boolean"}}
  {{$fieldname}}?: Map<string, boolean>;
    {{- else}}
  {{$fieldname}}?: Map<{{$property.AdditionalProperties | cleanRef}}>;
    {{- end}}
  {{- else if eq $property.Type "string"}}
  {{$fieldname}}?: string;
  {{- else}}
  {{$fieldname}}?: {{$property.Ref | cleanRef}};
  {{- end}}
  {{- end}}
}
{{- end}}

export const NakamaApi = (configuration: ConfigurationParameters = {
  basePath: BASE_PATH,
  bearerToken: "",
  password: "",
  username: "",
  timeoutMs: 5000,
}) => {
  const napi = {
    /** Perform the underlying Fetch operation and return Promise object **/
    doFetch(urlPath: string, method: string, queryParams: any, body?: any, options?: any): Promise<any> {
      const urlQuery = "?" + Object.keys(queryParams)
        .map(k => {
          if (queryParams[k] instanceof Array) {
            return queryParams[k].reduce((prev: any, curr: any) => {
              return prev + encodeURIComponent(k) + "=" + encodeURIComponent(curr) + "&";
            }, "");
          } else {
            if (queryParams[k] != null) {
              return encodeURIComponent(k) + "=" + encodeURIComponent(queryParams[k]) + "&";
            }
          }
        })
        .join("");

      const fetchOptions = {...{ method: method /*, keepalive: true */ }, ...options};
      fetchOptions.headers = {...options.headers};
      if (configuration.bearerToken) {
        fetchOptions.headers["Authorization"] = "Bearer " + configuration.bearerToken;
      } else if (configuration.username) {
        fetchOptions.headers["Authorization"] = "Basic " + btoa(configuration.username + ":" + configuration.password);
      }
      if(!Object.keys(fetchOptions.headers).includes("Accept")) {
        fetchOptions.headers["Accept"] = "application/json";
      }
      if(!Object.keys(fetchOptions.headers).includes("Content-Type")) {
        fetchOptions.headers["Content-Type"] = "application/json";
      }
      Object.keys(fetchOptions.headers).forEach((key: string) => {
        if(!fetchOptions.headers[key]) {
          delete fetchOptions.headers[key];
        }
      });
      fetchOptions.body = body;

      return Promise.race([
        fetch(configuration.basePath + urlPath + urlQuery, fetchOptions).then((response) => {
          if (response.status >= 200 && response.status < 300) {
            return response.json();
          } else {
            throw response;
          }
        }),
        new Promise((_, reject) =>
          setTimeout(reject, configuration.timeoutMs, "Request timed out.")
        ),
      ]);
    },
  {{- range $url, $path := .Paths}}
    {{- range $method, $operation := $path}}
    /** {{$operation.Summary}} */
    {{$operation.OperationId | camelCase}}(
    {{- range $parameter := $operation.Parameters}}
    {{- $camelcase := $parameter.Name | camelCase}}
    {{- if eq $parameter.In "path"}}
    {{- $camelcase}}{{- if not $parameter.Required }}?{{- end}}: {{$parameter.Type}},
    {{- else if eq $parameter.In "body"}}
      {{- if eq $parameter.Schema.Type "string"}}
    {{- $camelcase}}{{- if not $parameter.Required }}?{{- end}}: {{$parameter.Schema.Type}},
      {{- else}}
    {{- $camelcase}}{{- if not $parameter.Required }}?{{- end}}: {{$parameter.Schema.Ref | cleanRef}},
      {{- end}}
    {{- else if eq $parameter.Type "array"}}
    {{- $camelcase}}{{- if not $parameter.Required }}?{{- end}}: Array<{{$parameter.Items.Type}}>,
    {{- else if eq $parameter.Type "object"}}
    {{- $camelcase}}{{- if not $parameter.Required }}?{{- end}}: Map<{{$parameter.AdditionalProperties.Type}}>,
    {{- else if eq $parameter.Type "integer"}}
    {{- $camelcase}}{{- if not $parameter.Required }}?{{- end}}: number,
    {{- else}}
    {{- $camelcase}}{{- if not $parameter.Required }}?{{- end}}: {{$parameter.Type}},
    {{- end}}
    {{- " "}}
    {{- end}}options: any = {}): Promise<{{- if $operation.Responses.Ok.Schema.Ref | cleanRef -}} {{- $operation.Responses.Ok.Schema.Ref | cleanRef -}} {{- else -}} any {{- end}}> {
      {{- range $parameter := $operation.Parameters}}
      {{- $camelcase := $parameter.Name | camelCase}}
      {{- if $parameter.Required }}
      if ({{$camelcase}} === null || {{$camelcase}} === undefined) {
        throw new Error("'{{$camelcase}}' is a required parameter but is null or undefined.");
      }
      {{- end}}
      {{- end}}
      const urlPath = "{{- $url}}"
      {{- range $parameter := $operation.Parameters}}
      {{- $camelcase := $parameter.Name | camelCase}}
      {{- if eq $parameter.In "path"}}
         .replace("{{- print "{" $parameter.Name "}"}}", encodeURIComponent(String({{- $camelcase}})))
      {{- end}}
      {{- end}};

      const queryParams = {
      {{- range $parameter := $operation.Parameters}}
      {{- $camelcase := $parameter.Name | camelCase}}
      {{- if eq $parameter.In "query"}}
        {{$parameter.Name}}: {{$camelcase}},
      {{- end}}
      {{- end}}
      } as any;

      let _body = null;
      {{- range $parameter := $operation.Parameters}}
      {{- $camelcase := $parameter.Name | camelCase}}
      {{- if eq $parameter.In "body"}}
      _body = JSON.stringify({{$camelcase}} || {});
      {{- end}}
      {{- end}}

      return napi.doFetch(urlPath, "{{- $method | uppercase}}", queryParams, _body, options)
    },
    {{- end}}
  {{- end}}
  };

  return napi;
};
`

func snakeCaseToCamelCase(input string) (camelCase string) {
	isToUpper := false
	for k, v := range input {
		if k == 0 {
			camelCase = strings.ToLower(string(input[0]))
		} else {
			if isToUpper {
				camelCase += strings.ToUpper(string(v))
				isToUpper = false
			} else {
				if v == '_' {
					isToUpper = true
				} else {
					camelCase += string(v)
				}
			}
		}

	}
	return
}

func convertRefToClassName(input string) (className string) {
	cleanRef := strings.TrimPrefix(input, "#/definitions/")
	className = strings.Title(cleanRef)
	return
}

func main() {
	// Argument flags
	var output = flag.String("output", "", "The output for generated code.")
	flag.Parse()

	inputs := flag.Args()
	if len(inputs) < 1 {
		fmt.Printf("No input file found: %s\n", inputs)
		flag.PrintDefaults()
		return
	}

	fmap := template.FuncMap{
		"camelCase": snakeCaseToCamelCase,
		"cleanRef":  convertRefToClassName,
		"title":     strings.Title,
		"uppercase": strings.ToUpper,
	}

	input := inputs[0]
	content, err := ioutil.ReadFile(input)
	if err != nil {
		fmt.Printf("Unable to read file: %s\n", err)
		return
	}

	var schema struct {
		Paths map[string]map[string]struct {
			Summary     string
			OperationId string
			Responses   struct {
				Ok struct {
					Schema struct {
						Ref string `json:"$ref"`
					}
				} `json:"200"`
			}
			Parameters []struct {
				Name     string
				In       string
				Required bool
				Type     string   // used with primitives
				Items    struct { // used with type "array"
					Type string
				}
				Schema struct { // used with http body
					Type string
					Ref  string `json:"$ref"`
				}
			}
		}
		Definitions map[string]struct {
			Properties map[string]struct {
				Type  string
				Ref   string   `json:"$ref"` // used with object
				Items struct { // used with type "array"
					Type string
					Ref  string `json:"$ref"`
				}
				AdditionalProperties struct {
					Type string // used with type "map"
				}
				Format      string // used with type "boolean"
				Description string
			}
			Description string
		}
	}

	if err := json.Unmarshal(content, &schema); err != nil {
		fmt.Printf("Unable to decode input %s : %s\n", input, err)
		return
	}

	tmpl, err := template.New(input).Funcs(fmap).Parse(codeTemplate)
	if err != nil {
		fmt.Printf("Template parse error: %s\n", err)
		return
	}

	if len(*output) < 1 {
		tmpl.Execute(os.Stdout, schema)
		return
	}

	f, err := os.Create(*output)
	if err != nil {
		fmt.Printf("Unable to create file %s", err)
		return
	}
	defer f.Close()

	writer := bufio.NewWriter(f)
	tmpl.Execute(writer, schema)
	writer.Flush()
}
