package loop

// import (
// 	"context"

// 	"github.com/kyverno/kyverno-json/pkg/engine"
// )

// type loop[TPARENT any, TCHILD any, TRESPONSE any] struct {
// 	inner  engine.Engine[TCHILD, TRESPONSE]
// 	looper func(TPARENT) []TCHILD
// }

// func (b *loop[TPARENT, TCHILD, TRESPONSE]) Run(ctx context.Context, parent TPARENT) []TRESPONSE {
// 	var responses []TRESPONSE
// 	for _, child := range b.looper(parent) {
// 		responses = append(responses, b.inner.Run(ctx, child)...)
// 	}
// 	return responses
// }

// func New[TPARENT any, TCHILD any, TRESPONSE any](inner engine.Engine[TCHILD, TRESPONSE], looper func(TPARENT) []TCHILD) engine.Engine[TPARENT, TRESPONSE] {
// 	return &loop[TPARENT, TCHILD, TRESPONSE]{
// 		inner:  inner,
// 		looper: looper,
// 	}
// }
