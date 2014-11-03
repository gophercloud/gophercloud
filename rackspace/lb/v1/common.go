package v1

import (
	"fmt"
	"strconv"
)

// IDSliceToQueryString takes a slice of elements and converts them into a query
// string. For example, if name=foo and slice=[]int{20, 40, 60}, then the
// result would be `?name=20&name=40&name=60'
func IDSliceToQueryString(name string, ids []int) string {
	str := ""
	for k, v := range ids {
		if k == 0 {
			str += "?"
		} else {
			str += "&"
		}
		str += fmt.Sprintf("%s=%s", name, strconv.Itoa(v))
	}
	return str
}
