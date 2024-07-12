package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/gin-gonic/gin"
)

func createOutput(c *gin.Context){
	//receive contents of the template file
	file, header, err := c.Request.FormFile("template")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	defer file.Close()

	//read the contents of the template template
	templateBytes, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	//receive the json file
	jsonFile, _, err := c.Request.FormFile("data")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	defer jsonFile.Close()

	//read the contents of the json file
	jsonBytes, err := io.ReadAll(jsonFile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	//create a map to store the json data
	var data map[string]interface{}

	//unmarshal the json data into the map
	err = json.Unmarshal(jsonBytes, &data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	//parse the html template
	t, err := template.New("upload").Parse(string(templateBytes))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	
	//create a buffer to store the filled template
	var filledTemplate bytes.Buffer

	//execute the template with the json data, storing the result in the buffer
	t.Execute(&filledTemplate, data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	//initialise a new pdf generator
	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	//add a new page to the pdf generator with the filled template content
	pdfg.AddPage(wkhtmltopdf.NewPageReader(bytes.NewReader(filledTemplate.Bytes())))
	err = pdfg.Create()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	//generate the pdf file name based on the template file name
	filePath := pdfTemplateName(header.Filename)
	err = pdfg.WriteFile(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
}

//helper function to generate pdf file name based on the template file name	
func pdfTemplateName(templateFileName string) string {
	baseName := strings.TrimSuffix(filepath.Base(templateFileName), filepath.Ext(templateFileName))
	return fmt.Sprintf("generatedPDFs/%s.pdf", baseName)
}

func main(){
	r := gin.Default()
	r.POST("/generate", createOutput)
	r.Run()
}
