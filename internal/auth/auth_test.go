package auth

import "testing"

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
