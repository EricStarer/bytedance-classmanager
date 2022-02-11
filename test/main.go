package main

import "fmt"

var res map[string]string
var TeacherCourseRelationShip map[string][]string
var st map[string]string

func find(x string) bool {
	t := TeacherCourseRelationShip[x]
	for _, s := range t {
		if st[s] == "" {
			st[s] = "1"
			if res[s] == "" || find(res[s]) {
				res[s] = x
				return true
			}
		}
	}
	return false
}
func main() {
	hashmap := make(map[string][]string)
	hashmap["1"] = []string{"1", "2", "3"}
	hashmap["2"] = []string{"1"}
	hashmap["3"] = []string{"3"}
	//var u types.ScheduleCourseRequest
	TeacherCourseRelationShip = hashmap
	res = make(map[string]string)
	st = make(map[string]string)
	for k, _ := range TeacherCourseRelationShip {
		for k, _ := range res {
			delete(st, k)
		}
		find(k)
	}
	fmt.Println(res)
}
