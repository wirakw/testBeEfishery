package helper

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	// "github.com/iris-contrib/middleware/jwt"
)

var MySecret = []byte("eFisheryTest")

type Claims struct {
	Name      string `json:"name"`
	Role      string `json:"role"`
	Password  string `json:"password"`
	Timestamp string `json:"timestamp"`
	jwt.StandardClaims
}

func GenerateID(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

//
func getWeekAndYear(i time.Time) (int, int) {
	year, week := i.ISOWeek()
	return year, week
}

func DateEqual(x, y string) bool {
	time1, err := time.Parse("2006-01-02T15:04:05.999Z", x)
	if err != nil {
		err = nil
		time1, err = time.Parse("Mon Jan 02 15:04:05 GMT+07:00 2006", x)
		if err != nil {
			err = nil
			time1, _ = time.Parse("Mon Jan 02 2006 15:04:05 GMT+0700 (Western Indonesia Time)", y)
		}
	}
	time2, err := time.Parse("2006-01-02T15:04:05.999Z", y)
	if err != nil {
		err = nil
		time2, err = time.Parse("Mon Jan 02 15:04:05 GMT+07:00 2006", y)
		if err != nil {
			err = nil
			time2, _ = time.Parse("Mon Jan 02 2006 15:04:05 GMT+0700 (Western Indonesia Time)", y)
		}
	}
	year1, week1 := getWeekAndYear(time1)
	year2, week2 := getWeekAndYear(time2)

	return (year1 == year2) && (week1 == week2)
}

func MakeTimestampMillis() string {
	return strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
}

func GenerateToken(name, role, password string) string {
	expirationTime := time.Now().Add(24 * time.Hour)
	// Create the JWT claims, which includes the username and expiry time
	claims := &Claims{
		Name:      name,
		Role:      role,
		Password:  password,
		Timestamp: MakeTimestampMillis(),
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Sign and get the complete encoded token as a string using the secret
	tokenString, _ := token.SignedString(MySecret)
	return tokenString
}
