// package main

// import (
// 	"bytes"
// 	"encoding/json"
// 	"fmt"
// 	"html/template"
// 	"io"
// 	"net/http"
// 	"path/filepath"
// 	"strings"

// 	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
// 	"github.com/gin-gonic/gin"
// 	"github.com/google/uuid"
// )

// type Document struct{
// 	ID string `json:"id"`
// 	Body map[string]interface{} `json:"body"`
// }

// func createOutput(c *gin.Context){
// 	//receive template file
// 	file, header, err := c.Request.FormFile("template")
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
// 		return
// 	}
// 	defer file.Close()

// 	//read the contents of the template file
// 	templateBytes, err := io.ReadAll(file)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
// 		return
// 	}

// 	//receive the json data
// 	// jsonFile, _, err := c.Request.FormFile("data")
// 	// if err != nil {
// 	// 	c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
// 	// 	return
// 	// }
// 	// defer jsonFile.Close()
// 	jsonData := c.PostForm("data")
// 	//unmarshal the json data and store it in the map
// 	if err := json.Unmarshal([]byte(jsonData), &data); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid JSON data: " + err.Error()})
// 		return
// 	}

// 	//read contents of json file
// 	jsonBytes, err := io.ReadAll(jsonFile)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
// 		return
// 	}

// 	//create a map to store the json data
// 	var data Document

// 	//unmarshal json data and store it in document 
// 	err = json.Unmarshal(jsonBytes, &data)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
// 		return
// 	}

// 	if data.ID == "" {
// 		 data.ID = uuid.New().String()
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
// 			return 
// 		}
// 	}

// 	//parse the template 
// 	t, err := template.New("upload").Parse(string(templateBytes))
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
// 		return
// 	}

// 	//create a filledtemplate buffer to store the template data
// 	var filledTemplate bytes.Buffer

// 	//execute the template, storing the data in the buffer
// 	err = t.Execute(&filledTemplate, data.Body)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
// 		return
// 	}

// 	//create a new pdf generator
// 	pdfg, err := wkhtmltopdf.NewPDFGenerator()
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
// 		return
// 	}

// 	//add a new pdf page and write the buffer data to it
// 	pdfg.AddPage(wkhtmltopdf.NewPageReader(bytes.NewReader(filledTemplate.Bytes())))
// 	err = pdfg.Create()
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
// 		return
// 	}

// 	//create the pdf path basing on the input template
// 	filePath := generatedPDFName(header.Filename)

// 	//write the file to the path
// 	err = pdfg.WriteFile(filePath)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
// 		return
// 	}

// 	c.IndentedJSON(http.StatusCreated, data)
// }

// func generatedPDFName(templateFileName string) string {
// 	baseName := strings.TrimSuffix(filepath.Base(templateFileName), filepath.Ext(templateFileName))
// 	return fmt.Sprintf("generatedPDFs/%s.pdf", baseName)
// }

// func main(){
// 	r := gin.Default()
// 	r.POST("/generate-pdf", createOutput)
// 	r.Run()
// }