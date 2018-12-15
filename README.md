Welcome to golang easy swagger or simple goezswag!
    
For this beta version the annotations order matters so follow the examples careful.
Also bear in mind that case also matters.

Lets get to work, examples are in the main file (main.go).

- Lets begin by adding the main information of our project  
- Add these three anottations just once anywhere in your code as long as the comments are in the same order


    //@title your-api-name  
    //@version your-api-version  
    //@description your-api-description  


Now to map your structs, they should use the golang pattern and have the `//@model` annotation like in the example bellow:

    //@model MyDummyModel
    type MyDummyModel struct {
    	FieldOne	             string `json:"field_one" description:"field is the number one" validate:"required" example:"anything you want to add as example"`
    	FieldTwoNotMapped        string `validate:"required"`
    	FieldThreeForgotten      string
    	FieldFour			      string `json:"field_four" example:"you can add examples without descriptions"`
    }

- The mapper will see only the fields that are annotated with `json`

- Tags `description` and `example` are optionals. 

- Beware that by adding the tag `example` swagger will override the field type definition with what you have written as example.

Now to map your routes, simply add the annotations below anywhere in your code! (Just for the sake of organization, add above the handlers...)

    //@path your-api-path  
    //@method method
    //@summary summary of your api route  
    //@tags create  
    //@request MyDummyModel  
    //@response 200:MyDummyModelResponse 400:MyDummy400Model 500:MyDummy500Model    


- currently tags are in alpha ( it should be a list with empty spaces (create list get ...))

- the request and response fields should map to an existing annotated model in your project.
    they will map all your fields tagged with json

- The response field should be a list of key:value separated by empty spaces.  
  -  The number before `:` indicates the responseCode   
  -  The model  after  `:` indicates which model is to be associated with that response.    
    
- method annotation value should be in lower case   



Again, check main.go for examples on how to use annotations.
