package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

func Template(file string) string {
	html, _ := ioutil.ReadFile(file)
	return string(html)
}

type Grade struct {
	Student string
	Course  string
	Grade   string
}

type GradeAdmin struct {
	Grades []Grade
}

var grades GradeAdmin

func (this *GradeAdmin) studentList() []string {
	list := make([]string, 0)
	found := false
	for _, grade := range grades.Grades {
		for _, student := range list {
			if student == grade.Student {
				found = true
			}
		}
		if !found {
			list = append(list, grade.Student)
		}
		found = false
	}
	return list
}

func (this *GradeAdmin) courseList() []string {
	list := make([]string, 0)
	found := false
	for _, grade := range grades.Grades {
		for _, course := range list {
			if course == grade.Course {
				found = true
			}
		}
		if !found {
			list = append(list, grade.Course)
		}
		found = false
	}
	return list
}

func (this *GradeAdmin) Add(grade Grade) {
	this.Grades = append(this.Grades, grade)
}

func (this *GradeAdmin) String() string {
	var html string
	for _, grade := range grades.Grades {
		html += "<div>" +
			"<span>" + grade.Student + "</span>" +
			"<span>" + grade.Course + "</span>" +
			"<span>" + grade.Grade + "</span>" +
			"</div>"
	}
	return html
}

func (this *GradeAdmin) studentOptions() string {
	var html string
	list := this.studentList()
	for _, student := range list {
		html += `<option value="` + student + `">` + student + "</option>"
	}
	return html
}

func (this *GradeAdmin) courseOptions() string {
	var html string
	list := this.courseList()
	for _, course := range list {
		html += `<option value="` + course + `">` + course + "</option>"
	}
	return html
}

func (this *GradeAdmin) studentAverage(student string) float64 {
	var average float64
	cnt := 0.0
	for _, grade := range grades.Grades {
		if grade.Student == student {
			g, _ := strconv.ParseFloat(grade.Grade, 64)
			average += g
			cnt++
		}
	}
	return average / cnt
}

func (this *GradeAdmin) generalAverage() float64 {
	sAvg := 0.0
	list := this.studentList()
	cnt := 0.0
	for _, student := range list {
		sAvg += this.studentAverage(student)
		cnt++
	}
	return sAvg / cnt
}

func (this *GradeAdmin) courseAverage(course string) float64 {
	var average float64
	cnt := 0.0
	for _, grade := range grades.Grades {
		if grade.Course == course {
			g, _ := strconv.ParseFloat(grade.Grade, 64)
			average += g
			cnt++
		}
	}
	return average / cnt
}

func Register(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "POST":
		if err := req.ParseForm(); err != nil {
			fmt.Fprintf(res, "ParseForm() error %v", err)
			return
		}
		grade := Grade{Student: req.FormValue("student"),
			Course: req.FormValue("course"),
			Grade:  req.FormValue("grade")}
		grades.Add(grade)
		res.Header().Set(
			"Content-Type",
			"text/html",
		)
		fmt.Fprintf(
			res,
			Template("templates/register.html"),
			grades.String(),
		)
	case "GET":
		res.Header().Set(
			"Content-Type",
			"text/html",
		)
		fmt.Fprintf(
			res,
			Template("templates/register.html"),
			grades.String(),
		)
	}
}

func Student(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "POST":
		if err := req.ParseForm(); err != nil {
			fmt.Fprintf(res, "ParseForm() error %v", err)
			return
		}
		student := req.FormValue("student")
		res.Header().Set(
			"Content-Type",
			"text/html",
		)
		fmt.Fprintf(
			res,
			Template("templates/student.html"),
			grades.studentOptions(),
			grades.studentAverage(student),
		)
	case "GET":
		res.Header().Set(
			"Content-Type",
			"text/html",
		)
		fmt.Fprintf(
			res,
			Template("templates/student.html"),
			grades.studentOptions(),
			0.00,
		)
	}
}

func General(res http.ResponseWriter, req *http.Request) {
	res.Header().Set(
		"Content-Type",
		"text/html",
	)
	fmt.Fprintf(
		res,
		Template("templates/general.html"),
		grades.generalAverage(),
	)
}

func Course(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "POST":
		if err := req.ParseForm(); err != nil {
			fmt.Fprintf(res, "ParseForm() error %v", err)
			return
		}
		course := req.FormValue("course")
		res.Header().Set(
			"Content-Type",
			"text/html",
		)
		fmt.Fprintf(
			res,
			Template("templates/course.html"),
			grades.courseOptions(),
			grades.courseAverage(course),
		)
	case "GET":
		res.Header().Set(
			"Content-Type",
			"text/html",
		)
		fmt.Fprintf(
			res,
			Template("templates/course.html"),
			grades.courseOptions(),
			0.00,
		)
	}
}

func main() {
	http.HandleFunc("/register", Register)
	http.HandleFunc("/student", Student)
	http.HandleFunc("/general", General)
	http.HandleFunc("/course", Course)
	fmt.Println("Server on port 5400")
	http.ListenAndServe(":5400", nil)
}
