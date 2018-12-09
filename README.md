Welcome to golang easy swagger or simple goezswag!

	Since its too hard to find an easy and documented golang code to swagger tool I decided to create one

Lets get to work, examples are on the main file:

For now the orders matters so first add the basics info of the yaml - just once anywhere in your code as long as the comments are in the same order

//@title your-api-name  
//@version your-api-version  
//@description your-api-description  


Then simply add the routes you want  anywhere in your code! (Just for the sake of organization, add above the handlers...)
	- currently tags are in alpha ( it should be a list with empty spaces (create list get ...))

	- the request and response fields should map to an existing model in your project.
	    they will map all your fields tagged with json
	
	- The response field should be a list of key:value, like on the example below.  
	    The number before `:` indicates the responseCode and the model indicates which model is associated with that response.     

//@path `your-api-path`  
//@method `method`  
//@summary `summary-of-your-api-route`  
//@tags `create`  
//@request `MyDummyModel`  
//@response `200:MyDummyModelResponse 400:MyDummyModelResponse`  