package handler

import (
	"10-finishing-crud/connection"
	"10-finishing-crud/model"
	"10-finishing-crud/validation"
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"text/template"
	"time"

	"github.com/gorilla/mux" // npm route
)

// Function Route Home & Get Data
func HandleHome(w http.ResponseWriter, r *http.Request) { //ResponseWriter: untuk menampilkan data, Request: untuk menambahkan data
	w.Header().Set("Content-Type", "text/html; charset=utf-8") // Header berfungsi untuk menampilkan data. Data yang ditamplikan "text-html" /"json" / dll

	tmpt, err := template.ParseFiles("views/index.html") // template.ParseFiles berfungsi memparsing file yang disisipkan sebagai parameter

	if err != nil {
		w.Write([]byte("Message : " + err.Error()))
		return
	}

	// Query data dari database. context.background adalah
	dataProject, err := connection.Conn.Query(context.Background(), "SELECT id, project_name, start_date, end_date, description, technologies, image FROM tb_projects") // Query mengembalikan 2 nilai
	if err != nil {
		w.Write([]byte("Message :" + err.Error()))
		return
	}

	// Panggil struct untuk menampung data dari database
	var result []model.Project

	// Sebelum data ditampilkan looping terlebih dahulu
	for dataProject.Next() {
		var data = model.Project{}

		// scan setiap data yang ada di struct
		err := dataProject.Scan(&data.Id, &data.ProjectName, &data.StartDate, &data.EndDate, &data.Desc, &data.Tech, &data.Image)
		if err != nil {
			w.Write([]byte("Message : " + err.Error()))
			return
		}
		result = append(result, data)
	}

	// Buat map sebagai penampung data result
	data := map[string]interface{}{
		"Projects": result,
	}

	// Kemudian tampilkan seluruh data dari database
	w.WriteHeader(http.StatusOK)
	tmpt.Execute(w, data)
}

// Function Route Project
func HandleProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	tmpt, err := template.ParseFiles("views/project.html")

	if err != nil {
		w.Write([]byte("Message : " + err.Error()))
		return
	} else {
		tmpt.Execute(w, nil)
	}
}

// Function Route Contact
func HandleContact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")

	result, err := template.ParseFiles("views/contact.html")

	if err != nil {
		w.Write([]byte("Message : " + err.Error()))
		return
	} else {
		result.Execute(w, nil)
	}
}

// Function Detail Project
func HandleDetailProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	tmpt, err := template.ParseFiles("views/project-detail.html")

	if err != nil {
		w.Write([]byte("Message : " + err.Error()))
		return
	}

	// Tangkap id dari blog
	id, _ := strconv.Atoi(mux.Vars(r)["id"]) // strconv.Atoi untuk konversi string ke int.  mux.Vars() berfungsi untuk menangkap id dan mengembalikan 2 nilai parameter result dan error

	// // Panggil struct untuk menampung data dari database
	var project model.Project

	// QueryRow(get 1 data) data dari database yang id didatabase sama dengan id yang ditangkap di URL
	row := connection.Conn.QueryRow(context.Background(), "SELECT id, project_name, start_date, end_date, description, technologies, image FROM tb_projects WHERE id = $1", id)

	// Kemudia scan row
	err = row.Scan(&project.Id, &project.ProjectName, &project.StartDate, &project.EndDate, &project.Desc, &project.Tech, &project.Image)
	if err != nil {
		w.Write([]byte("Message :" + err.Error()))
		return
	}

	// Buat map sebagai penampung data result
	dataProject := map[string]interface{}{
		"Project": project,
	}

	tmpt.Execute(w, dataProject)
}

// Function Add Project
func HandleAddProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	temp, _ := template.ParseFiles("views/project.html")

	err := r.ParseForm()

	if err != nil {
		log.Fatal(err)
	}

	// Ambil value data input lalu tampung ke dalam variabel
	projectName := r.PostForm.Get("projectName")
	startDate := r.PostForm.Get("startDate")
	endDate := r.PostForm.Get("endDate")
	description := r.PostForm.Get("desc")

	// Buat array untuk menampung data checkbox
	var checkboxs []string

	// Jika didalam form checkboxs ada value-nya, maka append ke array checkboxs
	if r.FormValue("node") != "" {
		checkboxs = append(checkboxs, r.FormValue("node"))
	}
	if r.FormValue("angular") != "" {
		checkboxs = append(checkboxs, r.FormValue("angular"))
	}
	if r.FormValue("react") != "" {
		checkboxs = append(checkboxs, r.FormValue("react"))
	}
	if r.FormValue("typescript") != "" {
		checkboxs = append(checkboxs, r.FormValue("typescript"))
	}

	_, err = connection.Conn.Exec(context.Background(), "INSERT INTO tb_projects(project_name, start_date, end_date, description, technologies, image) VALUES ($1, $2, $3, $4, $5, 'saitama.png')", projectName, startDate, endDate, description, checkboxs)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Message :" + err.Error()))
		return
	}

	// Validation (Panggil function validation untuk melakukan validasi & pesan setelah data berhasil ditambahkan)
	validation := validation.NewValidation()
	project := model.Project{}

	data := make(map[string]interface{})
	vErrors := validation.Struct(project)

	// jika ada error tampilkan validasi, jika tidak tampilkan pesan
	if vErrors != nil {
		data["project"] = project
		// data["validation"] = vErrors
	} else {
		data["pesan"] = "Data project has been successfully added"
	}
	// End validaiton

	temp.Execute(w, data)
	// http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

// Function Edit Project
func HandleEditProject(w http.ResponseWriter, r *http.Request) {
	// Jika request methodnya get, maka ambil data lamanya dan tampilkan didalm input
	if r.Method == http.MethodGet {
		w.Header().Set("Content-type", "text/html; charset=utf-8")
		temp, err := template.ParseFiles("views/edit-project.html")

		if err != nil {
			w.Write([]byte("Message :" + err.Error()))
		}

		// Tangkap id
		id, _ := strconv.Atoi(mux.Vars(r)["id"])

		// Panggil struct
		var project model.Project

		// QueryRow(get 1 data) data dari database yang id didatabase sama dengan id yang ditangkap di URL
		row := connection.Conn.QueryRow(context.Background(), "SELECT id, project_name, start_date, end_date, description, technologies, image FROM tb_projects WHERE id = $1", id)

		// Scan hasil queryRow data dari database
		err = row.Scan(&project.Id, &project.ProjectName, &project.StartDate, &project.EndDate, &project.Desc, &project.Tech, &project.Image)
		if err != nil {
			w.Write([]byte("Message :" + err.Error()))
			return
		}

		// Buat map untuk menampung project yang telah discan. Kemudian didalam key "Project" buat map lagi lalu karena date akan diformat dulu
		dataProject := map[string]interface{}{
			"Project": map[string]interface{}{
				"Id":          project.Id,
				"ProjectName": project.ProjectName,
				"StartDate":   project.StartDate.Format("2006-01-02"),
				"EndDate":     project.EndDate.Format("2006-01-02"),
				"Desc":        project.Desc,
				"Tech":        project.Tech,
				"Image":       project.Image,
			},
		}

		// fmt.Println(dataProject)
		temp.Execute(w, dataProject)

	} else if r.Method == http.MethodPost {
		w.Header().Set("Content-type", "text/html; charset=utf-8")
		temp, _ := template.ParseFiles("views/edit-project.html")

		err := r.ParseForm()
		if err != nil {
			log.Fatal(err)
		}

		// panggil struct kemudian simapn ke vaiabel new
		var new model.Project

		// Setiap value yang dipanggil dengan Form.Get akan masuk kedalam variabel new(struct Project)
		new.Id, _ = strconv.Atoi(r.Form.Get("id")) // id didapat dari input yang type-nya hidden
		new.ProjectName = r.FormValue("projectName")
		new.StartDate, _ = time.Parse("2006-01-02", r.FormValue("startDate")) // parsing value startDate ke time dengan method time.Parse
		new.EndDate, _ = time.Parse("2006-01-02", r.FormValue("endDate"))
		new.Desc = r.FormValue("desc")
		new.Image = r.FormValue("image")

		// Jika didalam form checkboxs ada value-nya(node, angular, react, typescript), maka akan masuk kedalam variabel new
		if r.FormValue("node") != "" {
			new.Tech = append(new.Tech, r.FormValue("node"))
		}
		if r.FormValue("angular") != "" {
			new.Tech = append(new.Tech, r.FormValue("angular"))
		}
		if r.FormValue("react") != "" {
			new.Tech = append(new.Tech, r.FormValue("react"))
		}
		if r.FormValue("typescript") != "" {
			new.Tech = append(new.Tech, r.FormValue("typescript"))
		}

		fmt.Println("New Data :", new)

		// Kemudian panggil data dari database lalu UPDATE SET yang id dari database apakah sesuai dengan id di variabel new (dari URL)
		_, err = connection.Conn.Exec(context.Background(), "UPDATE tb_projects SET project_name=$1, start_date=$2, end_date=$3, description=$4, technologies=$5, image=$6 WHERE id=$7", new.ProjectName, new.StartDate, new.EndDate, new.Desc, new.Tech, new.Image, new.Id)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Message :" + err.Error()))
			return
		}

		// Validation (Panggil function validation untuk melakukan validasi & pesan setelah data berhasil ditambahkan)
		validation := validation.NewValidation()
		project := model.Project{}

		data := make(map[string]interface{})
		vErrors := validation.Struct(project)

		if vErrors != nil {
			data["project"] = project
			// data["validation"] = vErrors
		} else {
			data["pesan"] = "Data project has been successfully updated"
		}

		temp.Execute(w, data)
		// End validaiton

		// http.Redirect(w, r, "/", http.StatusMovedPermanently)
	}
}

// Function Delete Project
func HandleDeleteProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	// temp, _ := template.ParseFiles("views/index.html")

	// Tangkap id
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	// Kemudian delete data dimana id database sama dengan id yang ditangkap dari URL
	_, err := connection.Conn.Exec(context.Background(), "DELETE FROM tb_projects WHERE id = $1", id)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Message :" + err.Error()))
	}

	// Validation (Panggil function validation untuk melakukan validasi & pesan setelah data berhasil ditambahkan)
	validation := validation.NewValidation()
	project := model.Project{}

	data := make(map[string]interface{})
	vErrors := validation.Struct(project)

	if vErrors != nil {
		data["project"] = project
		// data["validation"] = vErrors
	} else {
		data["pesan"] = "Data project has been successfully deleted"
	}
	// End validaiton

	// temp.Execute(w, data)
	http.Redirect(w, r, "/", http.StatusSeeOther)

}
