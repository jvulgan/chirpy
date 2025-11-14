package auth

import (
	"net/http"
	"testing"
)

func TestHashPassword(t *testing.T) {
    cases := []struct{
        original_pwd string
        pwd_to_compare string
        expected_res bool
    } {
        {
            original_pwd: "pa$$word",
            pwd_to_compare: "pa$$word",
            expected_res: true,
        },
        {
            original_pwd: "pa$$word",
            pwd_to_compare: "different",
            expected_res: false,
        },
    }
    for _, c := range cases {
        t.Run("Test password hasing", func(t *testing.T) {
            hash, err := HashPassword(c.original_pwd)
            if err != nil {
                t.Errorf("error while hashing password")
                return
            }
            match, err := CheckPasswordHash(c.pwd_to_compare, hash)
            if err != nil {
                t.Errorf("error while checking if hash matches")
                return
            }
            if match != c.expected_res {
                t.Errorf("comparing %s and %s resulted in unexpected val %v", c.original_pwd, c.pwd_to_compare, c.expected_res)
                return
            }
        })
    }
}

func TestGetBearerToken(t *testing.T) {
    cases := []struct{
        name string
        header http.Header
        wantToken string
        wantErr bool
    } {
        {
            name: "valid auth header with bearer token",
            header: http.Header{"Authorization": []string{"Bearer sometoken"}},
            wantToken: "sometoken",
            wantErr: false,
        },
        {
            name: "empty headers",
            header: http.Header{},
            wantToken: "",
            wantErr: true,
        },
        {
            name: "no auth header",
            header: http.Header{"Non auth header": []string{"value"}},
            wantToken: "",
            wantErr: true,
        },
        {
            name: "no bearer token",
            header: http.Header{"Authorization": []string{"sometoken"}},
            wantToken: "",
            wantErr: true,
        },
    }
    for _, c := range cases {
        t.Run(c.name, func (t *testing.T) {
            gotToken, gotErr := GetBearerToken(c.header)
            if (gotErr != nil) != c.wantErr {
                t.Errorf("GetBearerToken() error: %v, wantErr %v", gotErr, c.wantErr)
                return
            }
            if gotToken != c.wantToken {
                t.Errorf("GetBearerToken() got token: %s, wantToken: %s", gotToken, c.wantToken)
                return
            }
        })
    }

}
