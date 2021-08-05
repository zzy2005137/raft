package main

type Mechanic struct{
	Key		string`json:"Key"`
	Value	string`json:"Value`
	Test	string`json:test`
}

type Measure struct{
	Id		string 		`json:"Id"`
	No		string 		`json:"No"`
	Time	string		`json:"Time"`
	Ddata	[3]float32 	`json:"D"`
	Ldata	[3]float32 	`json:"L"`
}