package main

import (
	"fmt"
	"github.com/gofhir/go-fhirpath/fhirpath"
)

func main() {
	questionnaireJSON := `{
  "resourceType": "Questionnaire",
  "id": "123123",
  "text": {
    "status": "generated",
    "div": "<div xmlns=\"http://www.w3.org/1999/xhtml\">\n      <pre>\n            1.Comorbidity?\n              1.1 Cardial Comorbidity\n                1.1.1 Angina?\n                1.1.2 MI?\n              1.2 Vascular Comorbidity?\n              ...\n            Histopathology\n              Abdominal\n                pT category?\n              ...\n          </pre>\n    </div>"
  },
  "url": "http://hl7.org/fhir/Questionnaire/3141",
  "title": "Cancer Quality Forum Questionnaire 2012",
  "status": "draft",
  "subjectType": [
    "Patient"
  ],
  "date": "2012-01",
  "item": [
    {
      "linkId": "1",
      "code": [
        {
          "system": "http://example.org/system/code/sections",
          "code": "COMORBIDITY"
        }
      ],
      "type": "group",
      "item": [
        {
          "linkId": "234234",
          "code": [
            {
              "system": "http://example.org/system/code/questions",
              "code": "COMORB"
            }
          ],
          "prefix": "1",
          "type": "choice",
          "answerValueSet": "http://hl7.org/fhir/ValueSet/yesnodontknow",
          "item": [
            {
              "linkId": "1.1.1",
              "code": [
                {
                  "system": "http://example.org/system/code/sections",
                  "code": "CARDIAL"
                }
              ],
              "type": "group",
              "enableWhen": [
                {
                  "question": "1.1",
                  "operator": "=",
                  "answerCoding": {
                    "system": "http://terminology.hl7.org/CodeSystem/v2-0136",
                    "code": "Y"
                  }
                }
              ],
              "item": [
                {
                  "linkId": "1.1.1.1",
                  "code": [
                    {
                      "system": "http://example.org/system/code/questions",
                      "code": "COMORBCAR"
                    }
                  ],
                  "prefix": "1.1",
                  "type": "choice",
                  "answerValueSet": "http://hl7.org/fhir/ValueSet/yesnodontknow",
                  "item": [
                    {
                      "linkId": "1.1.1.1.1",
                      "code": [
                        {
                          "system": "http://example.org/system/code/questions",
                          "code": "COMCAR00",
                          "display": "Angina Pectoris"
                        },
                        {
                          "system": "http://snomed.info/sct",
                          "code": "194828000",
                          "display": "Angina (disorder)"
                        }
                      ],
                      "prefix": "1.1.1",
                      "type": "choice",
                      "answerValueSet": "http://hl7.org/fhir/ValueSet/yesnodontknow"
                    },
                    {
                      "linkId": "1.1.1.1.2",
                      "code": [
                        {
                          "system": "http://snomed.info/sct",
                          "code": "22298006",
                          "display": "Myocardial infarction (disorder)"
                        }
                      ],
                      "prefix": "1.1.2",
                      "type": "choice",
                      "answerValueSet": "http://hl7.org/fhir/ValueSet/yesnodontknow"
                    }
                  ]
                },
                {
                  "linkId": "1.1.1.2",
                  "code": [
                    {
                      "system": "http://example.org/system/code/questions",
                      "code": "COMORBVAS"
                    }
                  ],
                  "prefix": "1.2",
                  "type": "choice",
                  "answerValueSet": "http://hl7.org/fhir/ValueSet/yesnodontknow"
                }
              ]
            }
          ]
        }
      ]
    },
    {
      "linkId": "2",
      "code": [
        {
          "system": "http://example.org/system/code/sections",
          "code": "HISTOPATHOLOGY"
        }
      ],
      "type": "group",
      "item": [
        {
          "linkId": "2.1",
          "code": [
            {
              "system": "http://example.org/system/code/sections",
              "code": "ABDOMINAL"
            }
          ],
          "type": "group",
          "item": [
            {
              "linkId": "2.1.2",
              "code": [
                {
                  "system": "http://example.org/system/code/questions",
                  "code": "STADPT",
                  "display": "pT category-----"
                }
              ],
              "type": "choice"
            }
          ]
        }
      ]
    }
  ]
}`
	// Example FHIR Patient resource (JSON)
	/*
			patientJSON := `{
		  "resourceType": "Patient",
		  "id": "example",
		  "address": [
		    {
		      "use": "home",
		      "city": "PleasantVille",
		      "type": "both",
		      "state": "Vic",
		      "line": [
		        "534 Erewhon St"
		      ],
		      "postalCode": "3999",
		      "period": {
		        "start": "1974-12-25"
		      },
		      "district": "Rainbow",
		      "text": "534 Erewhon St PeasantVille, Rainbow, Vic  3999"
		    }
		  ],
		  "managingOrganization": {
		    "reference": "Organization/1"
		  },
		"generalPractitioner": [
			{
				"reference": "Practitioner/1"
			},
			{
				"reference": "Practitioner/2"
			}],

		  "name": [
		    {
		      "use": "usual",
		      "given": [
		        "Peter",
		        "James"
		      ],
		      "family": "Chalmers",
		      "period": {
				"start": "2002",
				"end": "2004"
			  }
		    },
		    {
		      "use": "usual",
		      "given": [
		        "Jim"
		      ]
		    },
		    {
		      "use": "maiden",
		      "given": [
		        "Peter",
		        "James"
		      ],
		      "family": "Windsor",
		      "period": {
		        "end": "2002"
		      }
		    }
		  ],
		  "birthDate": "1974-12-25",
		  "deceased": {
		    "boolean": false
		  },
		  "active": true,
		  "identifier": [
		    {
		      "use": "usual",
		      "type": {
		        "coding": [
		          {
		            "code": "MR",
		            "system": "http://hl7.org/fhir/v2/0203"
		          }
		        ]
		      },
		      "value": "12345",
		      "period": {
		        "start": "2001-05-06"
		      },
		      "system": "urn:oid:1.2.36.146.595.217.0.1",
		      "assigner": {
		        "display": "Acme Healthcare"
		      }
		    }
		  ],
		  "telecom": [
		    {
		      "use": "home"
		    },
		    {
		      "use": "work",
		      "rank": 1,
		      "value": "(03) 5555 6473",
		      "system": "phone"
		    },
		    {
		      "use": "old",
		      "rank": 2,
		      "value": "(03) 3410 5613",
		      "system": "phone"
		    },
		    {
		      "use": "old",
		      "value": "(03) 5555 8834",
		      "period": {
		        "end": "2014"
		      },
		      "system": "phone"
		    }
		  ],
		  "gender": "male",
		  "contact": [
		    {
		      "name": {
		        "given": [
		          "Bénédicte"
		        ],
		        "family": "du Marché",
		        "_family": {
		          "extension": [
		            {
		              "url": "http://hl7.org/fhir/StructureDefinition/humanname-own-prefix",
		              "valueString": "VV"
		            }
		          ]
		        }
		      },
		      "gender": "female",
		      "period": {
		        "start": "2012"
		      },
		      "address": {
		        "use": "home",
		        "city": "PleasantVille",
		        "line": [
		          "534 Erewhon St"
		        ],
		        "type": "both",
		        "state": "Vic",
		        "period": {
		          "start": "1974-12-25"
		        },
		        "district": "Rainbow",
		        "postalCode": "3999"
		      },
		      "telecom": [
		        {
		          "value": "+33 (237) 998327",
		          "system": "phone"
		        }
		      ],
		      "relationship": [
		        {
		          "coding": [
		            {
		              "code": "N",
		              "system": "http://hl7.org/fhir/v2/0131"
		            }
		          ]
		        }
		      ]
		    }
		  ]
		}`
	*/
	// Example FHIRPath expression
	fhirPathExpr := "Questionnaire.item.item.item.enableWhen"

	// Extract value using FHIRPath
	result := fhirpath.Evaluate(questionnaireJSON, fhirPathExpr)

	// Print the correctly formatted JSON
	fmt.Println("Result:", string(result))
}
