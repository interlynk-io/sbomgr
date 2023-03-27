package licenses

import (
	"sort"
	"strings"
)

type LicenseStore struct {
	Nm string `json:"name,omitempty"`
	Ss string `json:"short,omitempty"`
	Ds bool   `json:"deprecated,omitempty"`
}

func (l LicenseStore) Short() string {
	return l.Ss
}

func (l LicenseStore) Name() string {
	return l.Nm
}

func (l LicenseStore) Deprecated() bool {
	return l.Ds
}

func (l LicenseStore) ValidSpdxLicense() bool {
	return l.Nm != ""

}
func NewLicenseFromID(lic string) []License {
	// The incoming license string could be a
	// - A SPDX short license ID
	// - A SPDX license expression
	// - A proprietary license id

	//NONE and NOASSERTION should be treated as
	// no license
	lcs := []License{}

	licenseLower := strings.ToLower(lic)

	if licenseLower == "none" || licenseLower == "noassertion" {
		return lcs
	}

	allLicenses := getIndividualLicenses(licenseLower)

	for _, l := range allLicenses {
		meta, ok := LookUp(l)
		if ok {
			lcs = append(lcs, meta)
		}
	}
	return lcs
}

// taken from https://github.com/spdx/tools-golang/blob/main/idsearcher/idsearcher.go#L208
func getIndividualLicenses(lic string) []string {
	// replace parens and '+' with spaces
	lic = strings.Replace(lic, "(", " ", -1)
	lic = strings.Replace(lic, ")", " ", -1)
	lic = strings.Replace(lic, "+", " ", -1)
	lic = strings.Replace(lic, ",", " ", -1) //changed from original

	// now, split by spaces, trim, and add to slice
	licElements := strings.Split(lic, " ")
	lics := []string{}
	for _, elt := range licElements {
		elt := strings.TrimSpace(elt)
		// don't add if empty or if case-insensitive operator
		if elt == "" || strings.EqualFold(elt, "AND") ||
			strings.EqualFold(elt, "OR") || strings.EqualFold(elt, "WITH") {
			continue
		}

		lics = append(lics, elt)
	}

	// sort before returning
	sort.Strings(lics)
	return lics
}

func LicenseObjectByName(name string) License {
	return meta{
		name: name,
	}
}
