package main

import (
	"context"
	"fmt"
)

func process(ctx context.Context) {
	ret, ok := ctx.Value("traceId").(int)
	if !ok {
		ret = 123456 //没有转成int成功就给个默认值
	}
	fmt.Printf("ret:%d\n", ret)

	s, _ := ctx.Value("session").(string)
	fmt.Printf("session:%s\n", s)

}

func main() {
	ctx := context.WithValue(context.Background(), "traceId", 13548782235)
	ctx = context.WithValue(ctx, "session", "happy") //继承了上一个ctx
	process(ctx)
}
