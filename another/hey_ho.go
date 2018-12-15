package another


//@path /v1/dummy/route
//@method get
//@summary get account holder
//@tags get
//@response 200:AccountHolder 400:Paymentz 500:MyDummyModelResponse
func anotherPackage() {

}


//@path /v1/dummy/getarray
//@method get
//@summary retrieves a list of accountHolders
//@tags list
//@response 200:[]AccountHolder 400:Paymentz 500:MyDummyModelResponse
func anotharPackageArrayResponse() {

}

//@path /v1/dummy/getarray
//@method patch
//@summary update someone
//@tags update
//@request AccountHolder
func anotharPackagePatchRequest() {

}