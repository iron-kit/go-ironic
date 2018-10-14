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

	pattern := regexp.MustCompile(`^ObjectIdHex\(\"(\w+)\"\)$`)
	res := pattern.MatchString(id.String())

	arr := pattern.FindStringSubmatch(id.String())

	fmt.Printf("%v", arr[len(arr)-1])
	assert.Equal(t, res, true)
}
