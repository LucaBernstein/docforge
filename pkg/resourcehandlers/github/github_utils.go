// SPDX-FileCopyrightText: 2020 SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package github

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/gardener/docforge/pkg/util/urls"
)

// Parse a GitHub URL into an incomplete ResourceLocator, without
// the SHA property.
func Parse(urlString string) (*ResourceLocator, error) {
	var (
		resourceType       ResourceType = -1
		repo               string
		path               string
		err                error
		resourceTypeString string
		shaAlias           string
		u                  *urls.URL
	)

	if u, err = urls.Parse(urlString); err != nil {
		return nil, err
	}

	host := u.Host
	sourceURLPathSegments := []string{}
	if len(u.Path) > 0 {
		// leading/trailing slashes
		_p := strings.TrimSuffix(u.Path[1:], "/")
		sourceURLPathSegments = strings.Split(_p, "/")
	}

	if len(sourceURLPathSegments) < 1 {
		return nil, fmt.Errorf("unsupported GitHub URL: %s. Need at least host and organization|owner", urlString)
	}

	var isRawAPI bool
	if "raw" == sourceURLPathSegments[0] {
		sourceURLPathSegments = sourceURLPathSegments[1:]
		isRawAPI = true
	}

	owner := sourceURLPathSegments[0]
	if len(sourceURLPathSegments) > 1 {
		repo = sourceURLPathSegments[1]
	}
	if len(sourceURLPathSegments) > 2 {
		// is this a raw.host content GitHub link?
		if isRawURL(u.URL) {
			resourceTypeString = "raw"
		} else {
			resourceTypeString = sourceURLPathSegments[2]
		}
		// {blob|tree|wiki|...}
		if resourceType, err = NewResourceType(resourceTypeString); err == nil {
			urlPathPrefix := strings.Join([]string{owner, repo, resourceTypeString}, "/")
			if isRawURL(u.URL) {
				// raw.host links have no resource type path segment
				urlPathPrefix = strings.Join([]string{owner, repo}, "/")
				shaAlias = sourceURLPathSegments[2]
			} else {
				// SHA aliases are defined only for blob/tree/raw objects
				if resourceType == Raw || resourceType == Blob || resourceType == Tree {
					// that would be wrong url but we make up for that
					if len(sourceURLPathSegments) < 4 {
						shaAlias = "master"
					} else {
						shaAlias = sourceURLPathSegments[3]
					}
				}
			}
			if len(shaAlias) > 0 {
				urlPathPrefix = strings.Join([]string{urlPathPrefix, shaAlias}, "/")
			}
			// get the github url "path" part without:
			// - leading "/"
			// - owner, repo, resource type, shaAlias segments if applicable
			if p := strings.Split(u.Path[1:], urlPathPrefix); len(p) > 1 {
				path = strings.TrimPrefix(p[1], "/")
			}
		}
		if err != nil {
			return nil, fmt.Errorf("unsupported GitHub URL: %s . %s", urlString, err.Error())
		}
	}
	if len(u.Fragment) > 0 {
		path = fmt.Sprintf("%s#%s", path, u.Fragment)
	}
	if len(u.RawQuery) > 0 {
		path = fmt.Sprintf("%s?%s", path, u.RawQuery)
	}
	ghRL := &ResourceLocator{
		Scheme:   u.Scheme,
		Host:     host,
		Owner:    owner,
		Repo:     repo,
		Type:     resourceType,
		Path:     path,
		SHAAlias: shaAlias,
		IsRawAPI: isRawAPI,
	}
	return ghRL, nil
}

func isRawURL(u *url.URL) bool {
	return strings.HasPrefix(u.Host, "raw.") || strings.HasPrefix(u.Path, "/raw")
}
