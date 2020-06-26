package model

import (
	"github.com/guregu/null"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

//MODEL DECLARATION
type TokenJwt struct { //return jwt token
	Token string `json:"token"`
	User
}

type ResponseStatus struct {
	Status  string `json:"status,omitempty"`
	ID      int    `json:"id,omitempty"`
	Message string `json:"error,omitempty"`
	Field   string `json:"field,omitempty"`
}

type User struct {
	UserID   int    `gorm:"primary_key;AUTO_INCREMENT" json:"user_id"`
	Name     string `json:"name,omitempty"`
	Role     string `json:"role"`
	Password string `json:"password,omitempty"`
	Phone    string `json:"phone,omitempty"`
	Token    string `gorm:"-" json:"token"`
}

type Login struct {
	Password string `json:"password,omitempty"`
	Phone    string `json:"phone,omitempty"`
}

type Storage struct {
	UUID         null.String `json:"uuid"`
	Komoditas    null.String `json:"komoditas"`
	AreaProvinsi null.String `json:"area_provinsi"`
	AreaKota     null.String `json:"area_kota"`
	Size         null.String `json:"size"`
	Price        null.String `price:"size"`
	USDPrice     string      `usd_price:"size"`
	TglParsed    null.String `json:"tgl_parsed"`
	Timestamp    null.String `json:"timestamp"`
}

type AggregateStorage struct {
	AreaProvinsi null.String `json:"area_provinsi"`
	TglParsed    null.String `json:"tgl_parsed"`
	Data         []DataProv  `json:"data"`
	Min          float64     `json:"min_size"`
	Max          float64     `json:"max_size"`
	Median       float64     `json:"median_size"`
	Avg          float64     `json:"avg_size"`
}

type DataProv struct {
	UUID      null.String `json:"uuid"`
	Komoditas null.String `json:"komoditas"`
	AreaKota  null.String `json:"area_kota"`
	Size      null.String `json:"size"`
	Price     null.String `price:"size"`
	USDPrice  string      `usd_price:"size"`
	Timestamp null.String `json:"timestamp"`
}

type CursUsd struct {
	Password string `json:"password,omitempty"`
}

// type Base struct {
// 	UUID      uuid.UUID  `gorm:"type:uuid;column:uuid;not null;"`
// 	CreatedAt int64      `json:"created_at"`
// 	UpdatedAt time.Time  `json:"updated_at"`
// 	DeletedAt *time.Time `sql:"index" json:"deleted_at,omitempty"`
// }

// func (base *Base) BeforeCreate(scope *gorm.Scope) error {
// 	u2 := uuid.NewV4()
// 	fmt.Println(u2)
// 	return scope.SetColumn("Uuid", u2)
// }
