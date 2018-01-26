package main

import (
	"errors"
	"fmt"
	"os"
	"path"

	"log"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

// student represents all of the data stored for a single student
type student struct {
	StudentId int    `db:"student_id"`
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
}

// teacher represents all of the data stored for a single student
type teacher struct {
	TeacherId int    `db:"teacher_id"`
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
}

// class represents all of the data stored for a single class
type class struct {
	ClassId         int    `db:"class_id"`
	ClassName       string `db:"class_name"`
	ClassIdentifier string `db:"class_identifier"`
	Section         string `db:"section"`
	TeacherId       int    `db:"teacher_id"`
	ClassSemester   string `db:"class_semester"`
	Classyear       int    `db:"class_year"`
}

// question represents all of the data stored for a single question
type question struct {
	QuestionId int    `db:"question_id"`
	Question   string `db:"question"`
	ClassId    int    `db:"class_id"`
	Responses  int    `db:"responses"`
}

// answer represents all of the data stored for a single answer
type answer struct {
	AnswerId   int    `db:"answer_id"`
	Answer     string `db:"answer"`
	QuestionId int    `db:"question_id"`
	Selections int    `db:"selections"`
}

// enrollment represents the relationship between students and courses;
// when a student is enrolled in a course, both the student ID and
// course ID appear in such a record
type enrolment struct {
	EnrolmentId int `db:"enrolment_id"`
	StudentId   int `db:"student_id"`
	classId    int `db:"class_id"`
}

// isUniqueViolation returns true if the supplied error resulted from
// a unique constraint violation
func isUniqueViolation(err error) bool {
	if err, ok := err.(*pq.Error); ok {
		return err.Code == "23505"
	}

	return false
}

// createCourse inserts a new course into the database
func createCourse(db *sqlx.DB) error {
	// check the arguments
	args := os.Args[2:]
	if len(args) != 3 {
		return fmt.Errorf("three arguments required: department, course number, and section")
	}

	// insert the data
	q := `INSERT INTO course (department, course_number, section)
                   VALUES ($1, $2, $3)`
	result, err := db.Exec(q, args[0], args[1], args[2])
	if err != nil {
		if isUniqueViolation(err) {
			return fmt.Errorf("course already exists")
		}

		return err
	}

	count, err := result.RowsAffected()
	if err != nil {
		return err
	}

	fmt.Printf("%d course(s) created.\n", count)

	return nil
}

// listCourses displays a table of all the courses in the database
func listCourses(db *sqlx.DB) error {
	return errors.New("unimplemented")
}

// createStudent inserts a new student into the database
func createStudent(db *sqlx.DB) error {
	q := `INSERT INTO student (student_id, first_name, last_name)
              VALUES(:student_id, :first_name, :first_name)`

	first_name := os.Args[2]
	last_name := os.Args[3]
	student_id := 1	

	qArgs := struct {
		StudentId int
		FirstName string
		LastName string
	}{student_id, first_name, last_name}

	var err error
	
	fmt.Printf("%v %v\n", os.Args[1], os.Args[2])

	res, err := db.NamedExec(q, qArgs)
	if err != nil {
		log.Fatalln(err)
	}
	
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("RowsAffected: %v", err)
	}
	
	fmt.Printf("Rows: %v\n", count)
	
	return err
	//return errors.New("unimplemented")
}

// printStudents prints all of the students appearing in the slice of
// students
func printStudents(students []student) {
	fmt.Printf("ID   Last name, First name\n")
	fmt.Printf("--------------------------------\n")
	for _, student := range students {
		fmt.Printf("%-4d %s, %s\n", student.StudentId, student.LastName, student.FirstName)
	}
}

// listStudents displays a table of all the students in the database
func listStudents(db *sqlx.DB) error {
	//q := `SELECT StudentId, FirstName, LastName FROM student`	
	q := `SELECT *
             FROM student`
	
	students := []student{}
	var err error

	if err := db.Select(&students, q); err != nil {
		return fmt.Errorf("Select: %v", err)
	}

	return err
	//return errors.New("unimplemented")
}

// enrollStudent enrolls a student in the indicated course
func enrollStudent(db *sqlx.DB) error {
	return errors.New("unimplemented")
}

// dropStudent drops a student from the indicated course
func dropStudent(db *sqlx.DB) error {
	return errors.New("unimplemented")
}

// courseRoster displays the course roster for the indicated course
func courseRoster(db *sqlx.DB) error {
	return errors.New("unimplemented")
}

// displayUsage displays a usage message
func displayUsage() {
	fmt.Fprintf(os.Stderr, `usage: %s command [args...]

Commands:
  help
  createCourse  department number   section
  listCourses
  createStudent firstName  lastName
  listStudents
  enrollStudent studentId  courseId
  dropStudent   studentId  courseId
  courseRoster  courseId
`, path.Base(os.Args[0]))
}

// executeCommand executes the specified command (specified on the
// command line)
func executeCommand(db *sqlx.DB) error {
	commands := map[string]func(*sqlx.DB) error{
		"help": func(db *sqlx.DB) error {
			displayUsage()
			return nil
		},
		"createCourse":  createCourse,
		"listCourses":   listCourses,
		"createStudent": createStudent,
		"listStudents":  listStudents,
		"enrollStudent": enrollStudent,
		"dropStudent":   dropStudent,
		"courseRoster":  courseRoster,
	}

	if len(os.Args) == 1 {
		return errors.New("required command missing")
	}

	// attempt to lookup and then execute the command
	if f, ok := commands[os.Args[1]]; ok {
		err := f(db)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("unknown command (%v)", os.Args[1])
	}

	return nil
}

var connectString string = "dbname=lab03 user=postgres port=8080 sslmode=disable"

// connectToDB connects to an SQL database and provides a pointer to
// the connected database
func connectToDB() (*sqlx.DB, error) {
	
	var db *sqlx.DB
	
	var err error

	db, err = sqlx.Open("postgres", connectString)
	if err != nil {
		return nil, fmt.Errorf("Open (%v): %v", connectString, err)
	}	
	
	//if err := db.Ping(); err != nil {
	//	return nil, fmt.Errorf("Ping: %v", err)
	//}

	fmt.Printf("Success\n")
	return db, nil

	//return nil, errors.New("must write code to connect to the database")


}

func main() {
	db, err := connectToDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot connect to database: %v\n", err)
		os.Exit(1)
	}

	err = executeCommand(db)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot execute command: %v\n\n", err)
		fmt.Fprintf(os.Stderr, "run '%s help' for help\n", path.Base(os.Args[0]))
		os.Exit(1)
	}
}
