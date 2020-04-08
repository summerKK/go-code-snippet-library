package main

var controller_template = `
package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	{{.Service.Name}}pkg "github.com/summerKK/go-code-snippet-library/koala/tools/koala/output/generate/{{.ProtoFileDir}}"
)

// server is used to implement helloworld.GreeterServer.
type {{.Service.Name}} struct {
	
}

{{range .Rpc}}
func (s *{{$.Service.Name}}) {{.Name}}(ctx context.Context, in *{{$.Service.Name}}pkg.{{.RequestType}}) (*{{$.Service.Name}}pkg.{{.ReturnsType}}, error) {
	panic("implement me")
}
{{end}}
`
