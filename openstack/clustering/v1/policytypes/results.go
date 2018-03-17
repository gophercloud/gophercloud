package policytypes

import (
	"encoding/json"
	"fmt"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

// commonResult is the response of a base result.
type commonResult struct {
	gophercloud.Result
}

// GetResult is the response of a Get operations.
type GetResult struct {
	commonResult
}

// Extract provides access to the individual policy type returned by Get and extracts PolicyTypes
func (r commonResult) Extract() (*PolicyType, error) {
	var s struct {
		PolicyType *PolicyType `json:"policy_type"`
	}
	err := r.ExtractInto(&s)
	return s.PolicyType, err
}

// ExtractPolicyTypes provides access to the list of PolicyTYpes in a page acquired from the ListDetail operation.
func ExtractPolicyTypes(r pagination.Page) ([]PolicyType, error) {
	var s struct {
		PolicyTypes []PolicyType `json:"policy_type"`
	}
	err := (r.(PolicyTypePage)).ExtractInto(&s)
	return s.PolicyTypes, err
}

// PolicyTypePage contains a single page of all policy types from a ListDetails call.
type PolicyTypePage struct {
	pagination.LinkedPageBase
}

// PolicyType represents a detailed policy type
type PolicyType struct {
	Name          string            `json:"name"`
	Schema        SchemaType        `json:"schema,omitempty"`
	SupportStatus SupportStatusType `json:"support_status,omitempty"`
}

type SchemaType struct {
	AvailabilityZone   map[string]interface{} `json:"availability_zone,omitempty"`
	EnableDrsExtension map[string]interface{} `json:"enable_drs_extension,omitempty"`
	Servergroup        ServerGroupType        `json:"servergroup,omitempty"`
}

type ServerGroupSchemaName struct {
	Description string `json:"description"`
	Required    bool   `json:"required"`
	Type        string `json:"type"`
	Updatable   bool   `json:"updatable"`
}

type ServerGroupSchemaPoliciesConstraints struct {
	Constraint []string `json:"constraint"`
	Type       string   `json:"type"`
}

type ServerGroupSchemaPolicies struct {
	Constraints []ServerGroupSchemaPoliciesConstraints `json:"constraints"`
	Default     string                                 `json:"default"`
	Description string                                 `json:"description"`
	Required    bool                                   `json:"required"`
	Type        string                                 `json:"type"`
	Updatable   bool                                   `json:"updatable"`
}

type ServerGroupSchemaType struct {
	Name     ServerGroupSchemaName     `json:"name"`
	Policies ServerGroupSchemaPolicies `json:"policies,omitempty"`
}

type ServerGroupType struct {
	Description string                `json:"description,omitempty"`
	Required    bool                  `json:"required"`
	Schema      ServerGroupSchemaType `json:"schema,omitempty"`
	Type        string                `json:"type,omitempty"`
	Updatable   bool                  `json:"updatable,omitempty"`
}

type SupportStatus struct {
	Status string `json:"status"`
	Since  string `json:"since"`
}

type SupportStatusType struct {
	SupportVersion map[string]interface{}
}

// IsEmpty determines if a ProfielType contains any results.
func (page PolicyTypePage) IsEmpty() (bool, error) {
	policyTypes, err := ExtractPolicyTypes(page)
	return len(policyTypes) == 0, err
}

func (pt *PolicyType) UnmarshalJSON(b []byte) error {
	type tmp PolicyType
	var s struct {
		tmp
		Schema        SchemaType        `json:"schema,omitempty"`
		SupportStatus SupportStatusType `json:"support_status,omitempty"`
	}
	err := json.Unmarshal(b, &s)
	if err != nil {
		fmt.Println("Error Unmarshal PolicyType")
		fmt.Printf("Detail Unmarshal Error: %v", err)
		return err
	}

	*pt = PolicyType(s.tmp)
	pt.Schema = s.Schema
	pt.SupportStatus = s.SupportStatus

	return nil
}

func (st *SchemaType) UnmarshalJSON(b []byte) error {
	type schema SchemaType
	var s struct {
		schema
	}
	err := json.Unmarshal(b, &s)
	if err != nil {
		fmt.Println("Error UnmarshalJSON for SchemaType")
		fmt.Printf("Detail Unmarshal Error: %v", err)
		return err
	}

	*st = SchemaType(s.schema)
	return nil
}

func (sst *SupportStatusType) UnmarshalJSON(b []byte) error {
	var f interface{}
	err := json.Unmarshal(b, &f)
	if err != nil {
		fmt.Println("Error Unmarshal support_status in JSON")
		fmt.Printf("Detail Unmarshal Error: %v", err)
		return err
	}
	support := SupportStatusType{SupportVersion: f.(map[string]interface{})}
	*sst = support

	return nil
}

func (sgt *ServerGroupType) UnmarshalJSON(b []byte) error {
	type tmp ServerGroupType
	var s struct {
		tmp
		Schema ServerGroupSchemaType `json:"schema,omitempty"`
	}
	err := json.Unmarshal(b, &s)
	if err != nil {
		fmt.Println("Error Unmarshal ServerGroup Type in JSON")
		fmt.Printf("Detail Unmarshal Error: %v", err)
		return err
	}
	*sgt = ServerGroupType(s.tmp)
	sgt.Schema = s.Schema

	return nil
}

func (sgst *ServerGroupSchemaType) UnmarshalJSON(b []byte) error {
	type tmp ServerGroupSchemaType
	var s struct {
		tmp
	}
	err := json.Unmarshal(b, &s)
	if err != nil {
		fmt.Println("Error Unmarshal ServerGroupSchemaType in JSON")
		fmt.Printf("Detail Unmarshal Error: %v", err)
		return err
	}
	*sgst = ServerGroupSchemaType(s.tmp)
	sgst.Name = s.tmp.Name
	sgst.Policies = s.tmp.Policies

	return nil
}

func (sgspc *ServerGroupSchemaPoliciesConstraints) UnmarshalJSON(b []byte) error {
	type tmp ServerGroupSchemaPoliciesConstraints
	var s struct {
		tmp
		Constraint []string `json:"constraint"`
	}
	err := json.Unmarshal(b, &s)
	if err != nil {
		fmt.Println("Error Unmarshal ServerGroupSchemaType in JSON")
		fmt.Printf("Detail Unmarshal Error: %v", err)
		return err
	}
	*sgspc = ServerGroupSchemaPoliciesConstraints(s.tmp)
	sgspc.Constraint = s.Constraint

	return nil
}
