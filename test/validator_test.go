package test

import (
	"testing"

	"github.com/harshvardha/artOfSoftwareEngineering/utility"
)

func TestGithubURLValidator(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		expected bool
	}{
		{"https_without_www_url", "https://github.com/harshvardha/repo", true},
		{"http_without_www_url", "http://github.com/harshvardha/repo", true},
		{"https_complete_url", "https://www.github.com/harshvardha/repo", true},
		{"without_http_url", "github.com/harshvardha/repo", true},
		{"https_with_forward_slash_url", "https://github.com/harshvardha/repo/", true},
		{"https_url_with_.git", "https://github.com/harshvardha/repo.git", false},
		{"gitlab_url", "https://gitlab.com/harshvardha/repo", false},
		{"https_missing_repo_name_url", "https://github.com/harshvardha", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utility.GithubURLValidator(tt.url)
			if result != tt.expected {
				t.Errorf("%s: invalid", tt.url)
			}
		})
	}
}

func TestUsernameValidator(t *testing.T) {
	tests := []struct {
		username string
		expected bool
	}{
		{"harsh_1999", true},
		{"harshvardhan28_04", true},
		{"Harsh6169", true},
		{"hArSH_1999_", true},
		{"HARSH_28_04_1999_", true},
		{"_har_sh_1999_", true},
		{"_h_a_r_S_H_2804_1999", true},
		{"#harsh1999_", false},
		{"#harsh6169_", false},
		{"#$Harsh_1999", false},
		{"$#Harsh_6169", false},
		{"$$#@HaRsH_28_04_1999", false},
	}

	for _, tt := range tests {
		t.Run(tt.username, func(t *testing.T) {
			result := utility.UsernameValidator(tt.username)
			if result != tt.expected {
				t.Errorf("%s: invalid", tt.username)
			}
		})
	}
}

func TestTagsValidator(t *testing.T) {
	tests := []struct {
		tags     string
		expected bool
	}{
		{"golang;distributed_systems;", true},
		{"tag_1;tag2;TAG_3;TAG4", true},
		{"1_tag;2_Tag;3_TAG;4_TaG", true},
		{"T_a_g_1;t_a__g_2;TA_G_3;T__ag___4", true},
		{"tag1tag2tag3", true},
		{"#tag1;$tag2;@Tag_3", false},
		{"$!tag;!#T_a__G2;__T__a__G1", false},
	}

	for _, tt := range tests {
		t.Run(tt.tags, func(t *testing.T) {
			result := utility.TagsValidator(tt.tags)
			if result != tt.expected {
				t.Errorf("%s: invalid", tt.tags)
			}
		})
	}
}
