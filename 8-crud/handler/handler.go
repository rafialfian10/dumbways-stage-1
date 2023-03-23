package handler

import (
	"7-routing/model"
	"7-routing/validation"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"text/template"
	"time"

	"github.com/gorilla/mux"
)

func HandleHome(w http.ResponseWriter, r *http.Request) { //ResponseWriter: untuk menampilkan data, Request: untuk menambahkan data
	w.Header().Set("Content-type", "text/html; charset=utf-8") // Header berfungsi untuk menampilkan data. Data yang ditamplikan "text-html" /"json" / dll

	tmpt, err := template.ParseFiles("views/index.html") // template.ParseFiles berfungsi memparsing file yang disisipkan sebagai parameter

	if err != nil {
		w.Write([]byte("Message : " + err.Error()))
		return
	}

	// Buat map, kemudian panggil local storage dengan id tertentu dan simpan local storage sebagai value dari "DataProjects"
	DataProjects := model.DataProjects
	dataProject := map[string]interface{}{
		"DataProjects": DataProjects,
	}
	// Kemudian tampilkan seluruh isi dari dari local storage
	tmpt.Execute(w, dataProject) // Execute berfungsi untuk mengeksekusi / menampilkan data dan harus ada 2 parameter (respon, Data)
}

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

func HandleAddProject(w http.ResponseWriter, r *http.Request) {
	// fmt.Println(r) // r berisi seluruh data form
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	temp, _ := template.ParseFiles("views/project.html")

	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	// Tangkap data form dengan method PostForm.Get
	time := strconv.FormatInt(time.Now().UnixNano(), 10) // buat unix id dengan menggunakan method time.Now()
	id, _ := strconv.Atoi(time)
	projectName := r.PostForm.Get("projectName")
	startDate := r.PostForm.Get("startDate")
	endDate := r.PostForm.Get("endDate")
	desc := r.PostForm.Get("desc")

	//  Buat variabel untuk menampung data checkbox
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

	// Panggil struct kemudian tampung kedalam variabel dan masukkan data yang sudah ditangkap dengan PostForm.Get() ke dalam object
	newData := model.Project{
		Id:          id,
		ProjectName: projectName,
		StartDate:   startDate,
		EndDate:     endDate,
		Desc:        desc,
		Tech:        checkboxs,
	}
	// Kemudian panggil local storage dan append object newData

	// Validation (Panggil function validation untuk melakukan validasi setelah data berhasil ditambahkan)
	validation := validation.NewValidation()
	project := model.Project{}

	data := make(map[string]interface{})
	vErrors := validation.Struct(project)

	if vErrors != nil {
		data["project"] = project
		// data["validation"] = vErrors
	} else {
		data["pesan"] = "Data project berhasil ditambahkan"
		model.DataProjects = append(model.DataProjects, newData)
	}
	// End validaiton

	temp.Execute(w, data)

	// Panggil method redirect agar Setelah data dikirim, maka routing akan berpindah ke halaman index
	// http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func HandleDetailProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	temp, err := template.ParseFiles("views/project-detail.html")

	if err != nil {
		w.Write([]byte("Message : " + err.Error()))
		return
	}

	// Tangkap id dari blog
	id, _ := strconv.Atoi(mux.Vars(r)["id"]) // strconv.Atoi untuk konversi string ke int.  mux.Vars() berfungsi untuk menangkap id dan mengembalikan 2 nilai parameter result dan error

	// Buat map, kemudian panggil local storage dengan id tertentu dan simpan local storage sebagai value dari "DataProjects"
	DataProjects := model.DataProjects
	dataProject := map[string]interface{}{
		"DataProjects": DataProjects[id],
	}

	temp.Execute(w, dataProject)
}

func HandleEditProject(w http.ResponseWriter, r *http.Request) {
	// Jika method adalah get, maka tampilkan data lama di input pada edit-project.html
	if r.Method == http.MethodGet {
		w.Header().Set("Content-type", "text/html; charset=utf-8")

		// Panggil id
		id, _ := strconv.Atoi(mux.Vars(r)["id"])

		// Panggil struct untuk menampung data yang akan di looping
		var project = model.Project{}

		// Panggil local storage untuk di looping
		var dataProjects = model.DataProjects

		for _, data := range dataProjects {
			if data.Id == id {
				project = data
			}
		}

		tmpt, err := template.ParseFiles("views/edit-project.html")
		if err != nil {
			w.Write([]byte("Message :" + err.Error()))
		}

		// tampung project ke dalam map dan execute data
		var data = map[string]interface{}{
			"DataProject": project,
		}
		tmpt.Execute(w, data)
		fmt.Println("Data :", data)

	} else if r.Method == http.MethodPost {
		r.ParseForm()

		// Panggil struct
		var new = model.Project{}

		// Tangkap value dari masing-masing input kemudian value akan ditampung ke dalam struct
		new.Id, _ = strconv.Atoi(r.Form.Get("id"))
		new.ProjectName = r.FormValue("projectName")
		new.StartDate = r.FormValue("startDate")
		new.EndDate = r.FormValue("endDate")
		new.Desc = r.FormValue("desc")

		// Jika didalam form checkboxs ada value-nya, maka append ke array checkboxs
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

		// Panggil local storage untuk di looping
		DataProjects := model.DataProjects

		for i, data := range DataProjects { // jika id local storage sama id di url maka timpa data local storage[i] dengan data baru(new)
			if data.Id == new.Id {
				DataProjects[i] = new
			}
		}

		// fmt.Println("New data :", new)
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
		fmt.Println("Data telah berhasil diedit")
	}
}

func HandleDeleteProject(w http.ResponseWriter, r *http.Request) {

	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	model.DataProjects = append(model.DataProjects[:id], model.DataProjects[id+1:]...)

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}
