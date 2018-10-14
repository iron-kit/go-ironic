package hptypes

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/mgo.v2/bson"
)

func Test_Reg(t *testing.T) {
	id := bson.NewObjectId()

	pattern := regexp.MustCompile(`^objectId:(\w+)$`)
	res := pattern.MatchString("sobjectId:1231" + id.Hex())

	if res {
		arr := pattern.FindStringSubmatch("objectId:" + id.Hex())
		fmt.Printf("%v", arr[len(arr)-1])
	}

	assert.Equal(t, res, true)
}
