# NOSQL Schema & Model

With buildDB.go GoCore will build a model & orm  for your application.  The objective for the project is to support as many NOSQL Databases that make sense from a common model runtime.

## Schema Definition

For NOSQL Document systems documents are typically stored as JSON and part of a collection or bucket.  The GoCore generated model structs & fields are all tagged with json in order to marshal and unmarshal json via the go runtime.

Schemas are to be defined within a directory beneath db/[Applicaiton Name]/schemas/[Version (1.0.0)].  

GoCore uses standard semantic versioning of Major, Minor, and Revision numbers to properly implement swagger api definitions.  Schema files must be located within the versioned directory.

Any file with a .json extension will be processed to create a NOSQL model.  NOTE:  Additional sub directories are recursively walked and processed to support larger application organization.

Each schema json file starts with and array of collections.  Each collection must have a name and schema.  The schema for the collection is the document you want to store to the NOSQL DB.  Each schema contains a name and fields array.

Each Field requires a name and type.  Each field can optionally contain an index, omitEmpty, and schema.  A schema definition is required for type object or objectArray.  GoCore will recursively process object and objectArrays to build type structs.

Available Types (type):

	int
	uint
	uint8
	uint16
	uint32
	uint64
	int8
	int16
	ing32
	int64
	float32
	float64
	string
	bool
	dateTime
	self
	byteArray
	object
	intArray
	float64Array
	stringArray
	boolArray
	objectArray
	selfArray

Available Indexes (index):  NOTE:  primary indexed fields will be set to uint64.

	primary
	index
	unique

Optional Fields:

	omitEmpty:  bool
	

Below is an example of an example schema:

    {
		"collections":
		[
			{
				"name":"Persons",
				"schema":
				{
					"name": "Person",
					"fields":
					[
						{
							"name":"Id",
							"type":"int",
							"index":"primary"
						},
						{
							"name":"worth",
							"type":"float64"
						},
						{
							"name":"first",
							"type":"string",
							"index":"index"
						},
						{
							"name":"isCool",
							"type":"bool"
						},
						{
							"name":"blob",
							"type":"byteArray"
						},
						{
							"name":"hand",
							"type":"object",
							"schema":
							{
								"name":"handDetails",
								"fields":
								[
									{
										"name":"fingerCount",
										"type":"int"
									}
								]
							}
	
						},
						{
							"name":"field7",
							"type":"intArray"
						},
						{
							"name":"field8",
							"type":"float64Array"
						},
						{
							"name":"field9",
							"type":"stringArray"
						},
						{
							"name":"field10",
							"type":"boolArray"
						},			
						{
							"name":"field12",
							"type":"objectArray",
							"schema":
							{
								"name":"field12Sub",
								"fields":
								[
									{
										"name":"subField",
										"type":"int"
									}
								]
							}
						}
					]
				}
			}
		]
	}

# Model API

For a collection we will use Persons as the collection with a document called Person to provide API examples.

# Collection API

### type Persons

	type Persons struct{}

### func (obj *Persons) Single(field string, value string) (retObj Person, e error)

Returns a Person based on an indexed field and value.  NOTE:  field must be indexed to return a single record.

Example Code:

	import(
		"myApp/model"
		"fmt"
	)
	
	func init(){
	
		var persons model.Persons

		person, err := persons.Single("name", "Dan")

		if err != nil{
			fmt.Println(err.Error())
			return
		}
		fmt.Println(person.Name)
		
	}

### func (obj *Persons) Search(field string, value string) (retObj []Person, e error)

Returns an array of type Person based on a field and value.


Example Code:

	import(
		"myApp/model"
		"fmt"
	)
	
	func init(){
	
		var persons model.Persons

		somePeople, err := persons.Search("name", "Dan")

		if err != nil{
			fmt.Println(err.Error())
			return
		}
		for _, people := range somePeople{ 
			fmt.Println(person.Name)
		}
		
	}


### func (obj *Persons) SearchAdvanced(field string, value string, limit int, skip int) (retObj []Person, e error)

Returns an array of type Person based on a field and value.  Additionally a limit of records can be returned.  Additionally records can be skipped.


Example Code:

	import(
		"myApp/model"
		"fmt"
	)
	
	func init(){
	
		var persons model.Persons

		somePeople err:= persons.SearchAdvanced("name", "Dan", 10, 0)

		if err != nil{
			fmt.Println(err.Error())
			return
		}
		for _, people := range somePeople{ 
			fmt.Println(person.Name)
		}
		
	}

### func (obj *Persons) All() (retObj []Person, e error)

Returns an array of type Person for the entire collection.


Example Code:

	import(
		"myApp/model"
		"fmt"
	)
	
	func init(){
	
		var persons model.Persons

		somePeople, err:= persons.All()

		if err != nil{
			fmt.Println(err.Error())
			return
		}
		
		for _, people := range somePeople{ 
			fmt.Println(person.Name)
		}
		
	}

### func (obj *Persons) AllAdvanced(limit int, skip int) (retObj []Person, e error)

Returns an array of type Person for the entire collection.  Additionally a limit of records can be returned.  Additionally records can be skipped.


Example Code:

	import(
		"myApp/model"
		"fmt"
	)
	
	func init(){
	
		var persons model.Persons

		somePeople, err := persons.AllAdvanced(10, 0)
		
		if err != nil{
			fmt.Println(err.Error())
			return
		}

		for _, people := range somePeople{ 
			fmt.Println(person.Name)
		}
		
	}

### func (obj *Persons) AllByIndex(index string) (retObj []Person, e error)

Returns an array of type Person for the entire collection sorted by index.  

Example Code:

	import(
		"myApp/model"
		"fmt"
	)
	
	func init(){
	
		var persons model.Persons

		somePeople, err:= persons.AllByIndex("name")
		
		if err != nil{
			fmt.Println(err.Error())
			return
		}

		for _, people := range somePeople{ 
			fmt.Println(person.Name)
		}
		
	}

### func (obj *Persons) AllByIndexAdvanced(index string, limit int, skip int) (retObj []Person, e error)

Returns an array of type Person for the entire collection sorted by index.  Additionally a limit of records can be returned.  Additionally records can be skipped.

Example Code:

	import(
		"myApp/model"
		"fmt"
	)
	
	func init(){
	
		var persons model.Persons

		somePeople, err:= persons.AllByIndexAdvanced("name", 10, 0)

		if err != nil{
			fmt.Println(err.Error())
			return
		}

		for _, people := range somePeople{ 
			fmt.Println(person.Name)
		}
		
	}

### func (obj *Persons) Range(min, max, field string) (retObj []Person, e error)

Returns an array of type Person for the field by including a range from and to.  The range can be any type that represents the field properly.  

Example Code:

	import(
		"myApp/model"
		"fmt"
	)
	
	func init(){
	
		var persons model.Persons

		somePeople, err := persons.Range("Bobby", "Dan", "name")

		if err != nil{
			fmt.Println(err.Error())
			return
		}		

		for _, people := range somePeople{ 
			fmt.Println(person.Name)
		}
		
	}

### func (obj *Persons) RangeAdvanced(min, max, field string, limit int, skip int) (retObj []Person, e error) 

Returns an array of type Person for the field by including a range from and to.  The range can be any type that represents the field properly.    Additionally a limit of records can be returned.  Additionally records can be skipped.

Example Code:

	import(
		"myApp/model"
		"fmt"
	)
	
	func init(){
	
		var persons model.Persons

		somePeople, err := persons.RangeAdvanced("Bobby", "Dan", "name", 10, 0)
		
		if err != nil{
			fmt.Println(err.Error())
			return
		}

		for _, people := range somePeople{ 
			fmt.Println(person.Name)
		}
		
	}


### func (obj *Persons) Index() error

Initializes indexes and buckets before saving an object.  Useful when starting your application.

Example Code:

	import(
		"myApp/model"
		"fmt"
	)
	
	func init(){
	
		var persons model.Persons

		err := persons.Index()

		if err != nil{
			fmt.Println(err.Error())
		}
		
	}

### func (obj *Persons) RunTransaction(objects []Person) error 

Runs all Person objects passed into a transaction.  If any errors occur the transaction is rolled back and error is returned.

Example Code:

	import(
		"myApp/model"
		"fmt"
	)
	
	func init(){
	
		var persons model.Persons

		err := persons.RunTransaction([]Person{Person{Name:"Dan"}, Person{Name:"Bobby"}})

		if err != nil{
			fmt.Println(err.Error())
		}
		
	}

### func (obj *Persons) New() *Person 

Returns an zero-valued Person struct.

Example Code:

	import(
		"myApp/model"
		"fmt"
	)
	
	func init(){
	
		var persons model.Persons

		p := persons.New()

		fmt.Println(p.Name)
		
	}

# Document API

### type Person

	type Person struct{
		Id      int          `json:"id"`
		Name   	string       `json:"name"`
		IsCool  bool         `json:"isCool"`
	}


### func (obj *Person) Save() error

Saves the document to the database.

Example Code:

	import(
		"myApp/model"
		"fmt"
	)
	
	func init(){
	
		var person model.Person

		person.Name = "David"

		
		err := person.Save()

		if err != nil{
			fmt.Println(err.Error())
		}
		
	}

### func (obj *Person) Delete() error

Deletes the document from the database.

Example Code:

	import(
		"myApp/model"
		"fmt"
	)
	
	func init(){
	
		var persons model.Persons

		person := persons.Single("name", "Dan")

		
		err := person.Delete()

		if err != nil{
			fmt.Println(err.Error())
		}
		
	}

### func (obj *Person) JSONString() (string, error)

Returns the document as a JSON formatted string.

Example Code:

	import(
		"myApp/model"
		"fmt"
	)
	
	func init(){
	
		var person model.Person

		person.Name = "David"

		
		value, err := person.JSONString()

		if err != nil{
			fmt.Println(err.Error())
			return
		}

		fmt.Println(value)
	}

Output

	{"id":1, "name":"David", "IsCool":false}

### func (obj *Person) JSONBytes() ([]byte, error)

Returns the document as a JSON formatted byte array.

Example Code:

	import(
		"myApp/model"
		"fmt"
	)
	
	func init(){
	
		var person model.Person

		person.Name = "David"

		
		value, err := person.JSONBytes()

		if err != nil{
			fmt.Println(err.Error())
			return
		}

		fmt.Printf("%+v\n", bytes)
	}

# Bucket API

If your NOSQL DB provides setting and getting by key value the bucket API provides methods to set, get, and delete by key.

### type Bucket

	type Bucket struct{
		Name string
	}


### func (obj *Bucket) SetKeyValue(key interface{}, value interface{}) error

Sets a key / value pair for a bucket.

Example Code:

	import(
		"myApp/model"
		"fmt"
	)
	
	func init(){
	
		var bucket model.Bucket

		bucket.Name = "Items"

		err := bucket.SetKeyValue("C1531BD5-C78E-AAB1-C6D7-F1D5537A524F", "My Special String")

		if err != nil{
			fmt.Println(err.Error())
			return
		}

		err := bucket.SetKeyValue("3DE4ECE7-3F62-EB80-042F-011378327924", model.Person{})

		if err != nil{
			fmt.Println(err.Error())
			return
		}
		
	}

### func (obj *Bucket) GetKeyValue(key interface{}, value interface{}) error

Gets a key / value pair for a bucket.

Example Code:

	import(
		"myApp/model"
		"fmt"
	)
	
	func init(){
	
		var bucket model.Bucket

		bucket.Name = "Items"

		var s string

		err := bucket.GetKeyValue("C1531BD5-C78E-AAB1-C6D7-F1D5537A524F", &s)

		if err != nil{
			fmt.Println(err.Error())
			return
		}

		fmt.Println(s)
		
	}

### func (obj *Bucket) DeleteKey(key interface{}) error

Deletes a key / value pair for a bucket.

Example Code:

	import(
		"myApp/model"
		"fmt"
	)
	
	func init(){
	
		var bucket model.Bucket{}

		bucket.Name = "Items"

		var s string

		err := bucket.DeleteKey("C1531BD5-C78E-AAB1-C6D7-F1D5537A524F")

		if err != nil{
			fmt.Println(err.Error())
			return
		}

		fmt.Println(s)
		
	}