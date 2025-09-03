package relation

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

//	func TestWriteKeto(t *testing.T) {
//		err := WriteKeto()
//		fmt.Println(err)
//		assert.False(t, true)
//	}
func init() {
	Initialize(WithWriteAddr("keto.dev.orb.local:4467"), WithReadAddr("keto.dev.orb.local:4466"))
}

func TestInsertAndFetchTuples(t *testing.T) {
	ctx := context.Background()
	resp, err := QuerySubjectByObjectRelation(ctx, "Image", "image:abc", "viewer")
	fmt.Println(resp, err)
	assert.False(t, true)
}

func TestQueryObjectBySubjectRelation(t *testing.T) {
	ctx := context.Background()
	resp, err := QueryObjectBySubjectIdRelation(ctx, "Image", "user:567", "viewer")
	fmt.Println(resp, err)
	assert.False(t, true)
}

func TestQueryObjectBySubjectSetRelation(t *testing.T) {
	ctx := context.Background()
	resp, err := QueryObjectBySubjectSetRelation(ctx, "Image", "User", "user:kkkkkk", "viewer")
	fmt.Println(resp, err)
	assert.False(t, true)
}
