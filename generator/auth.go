package generator

import (
	"fmt"
	"os"
	"strings"
	"text/template"
)

func GenerateAuth(fields string) {

	if fields == "" {
		panic("fields required")
	}

	raw := strings.Split(fields, ",")

	list := []string{}

	for _, f := range raw {
		if strings.ToLower(f) != "password" {
			list = append(list, f)
		}
	}

	os.Mkdir("auth", 0755)

	data := map[string]any{
		"Fields": list,
	}

	write("auth/model.go", modelTpl, data)
	write("auth/dto.go", dtoTpl, data)
	write("auth/routes.go", routesTpl, data)
	write("auth/handlers.go", handlersTpl, data)
	write("auth/service.go", serviceTpl, data)
	write("auth/jwt.go", jwtTpl, data)
	write("auth/middleware.go", middlewareTpl, data)

	fmt.Println("✅ auth module generated")
}

func write(path string, tpl string, data any) {

	f, _ := os.Create(path)
	defer f.Close()

	t := template.Must(template.New("tpl").Funcs(template.FuncMap{
		"title": strings.Title,
	}).Parse(tpl))

	t.Execute(f, data)
}
