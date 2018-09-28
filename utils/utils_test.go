package utils

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

// func Test_ISOTime2Time(t *testing.T) {
// 	time := "2018-09-06 09:55:42.405 +0000 UTC"
// 	paresdTime, err := ISOTime2Time(time)
// 	if err != nil {
// 		t.Errorf("测试失败: %s", err.Error())
// 		return
// 	}

// 	y, m, d := paresdTime.Date()

// 	t.Logf("测试通过，year: %d, month: %d, day: %d", y, m, d)
// }

type Ark struct {
	Name string `json:"name,omitempty"`
	ID   int    `json:"id,omitempty"`
	Desc bool   `json:"desc,omitempty"`
}

func Test_IsZeroOkey(t *testing.T) {
	// ark := new(Ark)
	var ark Ark
	res := IsZero(ark)

	assert.Equal(t, res, true)
}

func Test_IsZeroError(t *testing.T) {
	// var ark Ark
	ark := new(Ark)
	res := IsZero(ark)

	assert.Equal(t, res, false)
}

func Test_IsZeroStructFieldOkey(t *testing.T) {
	ark := new(Ark)

	if ark != nil {
		assert.Equal(t, IsZero(ark.ID), true)
		assert.Equal(t, IsZero(ark.Name), true)
		assert.Equal(t, IsZero(ark.Desc), true)
	}
}

func Test_IsZeroStructFieldError(t *testing.T) {
	ark := new(Ark)
	ark.Name = "Hello"
	if ark != nil {
		assert.Equal(t, IsZero(ark.ID), true)
		assert.Equal(t, IsZero(ark.Name), false)
		assert.Equal(t, IsZero(ark.Desc), true)
	}
}

func Test_IsZeroStructFieldNullError(t *testing.T) {
	jsonTemp := `
	{
		"id": 1,
		"name": "",
		"desc": false
	}
	`

	ark := new(Ark)
	// 	ark.Name = ""
	json.Unmarshal([]byte(jsonTemp), ark)
	if ark != nil {
		assert.Equal(t, IsZero(ark.ID), false)
		assert.Equal(t, IsZero(ark.Name), true)
		assert.Equal(t, IsZero(ark.Desc), false)
	}
}
