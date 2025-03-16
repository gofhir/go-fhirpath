package main

import (
	"fmt"
	"github.com/gofhir/go-fhirpath/fhirpath"
)

func main() {
	// Example FHIR Patient resource (JSON)
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

	// Example FHIRPath expression
	fhirPathExpr := "Patient.generalPractitioner.all($this is Practitioner)"

	// Extract value using FHIRPath
	result := fhirpath.Evaluate(patientJSON, fhirPathExpr)

	// Print the correctly formatted JSON
	fmt.Println("Result:", string(result))
}
