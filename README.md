Welcome to golang easy swagger or simple goezswag!
    
Lets get to work, examples are in the main file.

For this beta version the annotations order matters so follow the examples .

- Lets begin by adding the main information of our project  
- Add these three anottations just once anywhere in your code as long as the comments are in the same order

//@title your-api-name  
//@version your-api-version  
//@description your-api-description  

Your structs should match the golang pattern like in the example bellow:

    type MyDummyModel struct {
    	FieldOne	             string `json:"field_one" description:"field is the number one" validate:"required"`
    	FieldTwoNotMapped        string `validate:"required"`
    	FieldThreeForgotten      string
    	FieldFour			      string `json:"field_four"`
    }

The mapper will use the fields that are annotated with `json`, plus the description (if provided)

Then simply add the annotation routes you want  anywhere in your code! (Just for the sake of organization, add above the handlers...)  

	- currently tags are in alpha ( it should be a list with empty spaces (create list get ...))

	- the request and response fields should map to an existing model in your project.
	    they will map all your fields tagged with json
	
	- The response field should be a list of key:value separated with empt, like on the example below.  
	    The number before `:` indicates the responseCode   
	    and the model indicates which model is to be associated with that response.    
	    
	- method annotation value should be in lower case   

//@path `your-api-path`  
//@method `method` . 
//@summary `summary of your api route`  
//@tags `create`  
//@request `MyDummyModel`  
//@response `200:MyDummyModelResponse`  


Again, check main.go for examples on how to use annotations.
