package tool

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"
)

func DecodeStdEncodeing(context string) (string, error) {
	bytes, err := base64.StdEncoding.DecodeString(context)
	if err != nil {
		return "", err
	} else {
		return string(bytes), nil
	}
	return "", nil
}

func Unicode2utf8(unicodeStr string) (utf8Str string) {
	//unicodeStr = `\u5546\u54c1\u63cf\u8ff0\`
	strlen := len(unicodeStr)
	if strings.HasSuffix(unicodeStr, `\`) {
		unicodeStr = unicodeStr[0 : strlen-1]
	}
	utf8Str, _ = strconv.Unquote(`"` + unicodeStr + `"`)
	return
}

type UrlValuesSorter []UrlValueSorterItem

type UrlValueSorterItem struct {
	Key string
	Val string
}

func NewUrlValuesSorter(values *url.Values, excludeItems *[]string) UrlValuesSorter {
	uvs := make(UrlValuesSorter, 0, len(*values))
	for valuesKey, valuesItem := range *values {
		skip := false
		for _, excludeItem := range *excludeItems {
			if excludeItem == valuesKey {

				skip = true
				break
			}
		}
		if !skip {
			uvs = append(uvs, UrlValueSorterItem{valuesKey, valuesItem[0]})
		}
	}
	return uvs
}

func (uvs UrlValuesSorter) Body() (ret string) {
	maxIndex := len(uvs) - 1
	for index, value := range uvs {
		format := "%s=%s"
		ret = ret + fmt.Sprintf(format, value.Key, value.Val)
		if index != maxIndex {
			ret = ret + "&"
		}
	}
	return
}

func (uvs UrlValuesSorter) Len() int {
	return len(uvs)
}

func (uvs UrlValuesSorter) Less(i, j int) bool {
	return uvs[i].Key < uvs[j].Key
}

func (uvs UrlValuesSorter) Swap(i, j int) {
	uvs[i], uvs[j] = uvs[j], uvs[i]
}

func TestUrlValuesSort() {
	form := url.Values{
		"d": []string{"1"},
		"c": []string{"2"},
		"t": []string{"3"},
		"u": []string{"4"},
		"a": []string{"5"},
		"e": []string{"6"},
		"i": []string{"7"},
		"s": []string{"8"},
		"g": []string{"9"},
	}

	excludeItems := []string{}
	uvs := NewUrlValuesSorter(&form, &excludeItems)
	sort.Sort(uvs)

	for _, item := range uvs {
		fmt.Printf("%s:%s\n", item.Key, item.Val)
	}
}
