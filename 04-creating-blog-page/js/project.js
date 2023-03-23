
let data = [];

function addProject(event) {
    event.preventDefault();
    
    let projectName = document.getElementById("project-name").value;
    let startDate = document.getElementById("start-date").value;
    let endDate = document.getElementById("end-date").value;
    let desc = document.getElementById("desc").value; 
    let image = document.getElementById("image").files; // Untuk menangkap image, pdf, dll. Ini akan mengambil seluruh data files, tapi yang dibutuhkan adalah url/path-nya
    
    // console.log("Gambar :", image[0]); 
    // console.log("URL Gambar :", getImage);
    
    // logic alert
    if(projectName == "" || startDate == "" || endDate == "" || desc == "" || image == "" ){
        return alert("All input field must be not empty");
    } else {
        alert("New Card Created Successfully")
    }
    let getImage = URL.createObjectURL(image[0]); // fungsi URL.createObjectURL untuk menangkap url gambar / path. Posisi dibawah kondisi karena, jika diatas kondisi akan error karena gambar belum ada gambar yang di inputkan

   

    // Logic get all checked
    let result = []
    let checkboxes = document.getElementsByName("mycheckbox");

    for(let i = 0; i < checkboxes.length; i++) {
        if(checkboxes[i].checked == true) {
            result.push(checkboxes[i].value);
        } 
    }

    // Id random with date method
    let uniqId = Date.now()
    
    let blog = {
        id: uniqId,
        title: projectName,
        startDate: startDate,
        endDate: endDate,
        tech: result,
        desc: desc,
        image: getImage,
    }

    data.push(blog);
    console.log("Data :", data)
    getBlog();
}

function getBlog() {
    let cardContainer = document.getElementById("card-container");
    cardContainer.innerHTML = ''; 

    for(let i = 0; i < data.length; i++) {
        cardContainer.innerHTML += `<div class="cards">
                                        <img src="${data[i].image}">
                                        <h3><a href="blog.html?${data[i].id}">Dumbways Mobile App - 2021</a></h3>
                                        <p class="date">durasi: 3 bulan</p>
                                        <p class="desc">${data[i].desc}</p>
                                        <div class="icons">
                                            ${data[i].tech.includes("node") ? '<img src="img/node.png"/>': ""}
                                            ${data[i].tech.includes("angular") ? '<img src="img/angular.png"/>' : ""}
                                            ${data[i].tech.includes("react") ? '<img src="img/react.png"/>' : ""}
                                            ${data[i].tech.includes("typescript") ? '<img src="img/typescript.png"/>' : ""}
                                        </div>
                                        <div>
                                            <button type="submit" id="edit" class="edit btn-edit">edit</button>
                                            <button type="submit" id="delete" class="delete btn-delete">Delete</button>
                                        </div>
                                    </div>`
    };
}
   
// --------------------------------------------------------------------
// Logic Input Color
let inputBorderColor1 = document.querySelectorAll(".border-color");
let inputBorderColor2 = document.querySelector(".form-container form");
let input = document.querySelectorAll("input");
let btnSubmit = document.querySelector(".btn-submit");
let imgUpload = document.querySelector(".img-upload");

inputBorderColor1.forEach(function(ibc1) {
    ibc1.addEventListener('click', function(event) {
        
        inputBorderColor1.forEach(function(ibc1) {
            if(ibc1.classList.contains("input-color")){
                ibc1.classList.remove("input-color");
            }
        });    
        event.target.classList.add("input-color");
    });

    btnSubmit.addEventListener('mouseover', function() {
        inputBorderColor1.forEach(function(ibc1) {
            ibc1.classList.remove("input-color");
        });  
    });  
});

inputBorderColor2.addEventListener('mouseenter', function(){
    inputBorderColor1.forEach(function(ibc1){
        ibc1.classList.add("input-border");
    });
    
    inputBorderColor2.addEventListener('mouseleave', function(){
        inputBorderColor1.forEach(function(ibc1){
            ibc1.classList.remove("input-border");
        });
    });
});

inputBorderColor2.addEventListener('mouseenter', function(){
    imgUpload.classList.add("input-border");
});
    
inputBorderColor2.addEventListener('mouseleave', function(){
    imgUpload.classList.remove("input-border");
});



