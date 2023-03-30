package models

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID                   primitive.ObjectID `bson:"_id"`
	UserName             string             `bson:"UserName"`
	Email                string             `bson:"Email"`
	EmailConfirmed       bool               `bson:"EmailConfirmed"`
	PasswordHash         string             `bson:"PasswordHash"`
	FullName             string             `bson:"FullName"`
	PhoneNumber          string             `bson:"PhoneNumber"`
	PhoneNumberConfirmed bool               `bson:"PhoneNumberConfirmed"`
	TwoFactorEnabled     bool               `bson:"TwoFactorEnabled"`
	IsAdmin              bool               `bson:"IsAdmin"`
	CreatedOn            primitive.DateTime `bson:"CreatedOn"`
}

func (u *User) Init() {
	u.ID = primitive.NewObjectID()
	u.CreatedOn = primitive.DateTime(time.Now().Unix() * 1000)
}
func (db *DB) AddUser(m *User) (id string, err error) {
	ctx := context.Background()
	m.Init()
	_, err = db.Users.InsertOne(ctx, m)
	if err == nil {
		return m.ID.Hex(), nil
	}
	return "", err
}
func (db *DB) GetUserById(idstr string) (v *User, err error) {

	idOj, _ := primitive.ObjectIDFromHex(idstr)
	ctx := context.Background()
	err = db.Users.FindOne(ctx, bson.M{"_id": idOj}).Decode(&v)

	if err != nil {
		return nil, errors.New("User not exists")
	}
	return v, nil
}

// func GetAllUsers() map[string]*User {
// 	return UserList
// }

func (db *DB) UpdateUser(m *User) (User, error) {
	ctx := context.Background()
	u, e := db.GetUserById(m.ID.Hex())
	if e == nil {
		if m.UserName != "" {
			u.UserName = m.UserName
		}
		if m.Email != "" {
			u.Email = m.Email
		}
		if m.FullName != "" {
			u.FullName = m.FullName
		}
		if m.PhoneNumber != "" {
			u.PhoneNumber = m.PhoneNumber
		}
	}
	filter := bson.M{"_id": m.ID}
	update := bson.M{"$set": bson.M{
		"UserName":    u.UserName,
		"Email":       u.Email,
		"FullName":    u.FullName,
		"PhoneNumber": u.PhoneNumber,
	}}
	_, err := db.Users.UpdateOne(ctx, filter, update)
	if err != nil {
		return User{}, errors.New("User not exists")
	}
	return *u, nil
}
func (db *DB) DeleteUser(idstr string) (err error) {

	id, _ := primitive.ObjectIDFromHex(idstr)
	ctx := context.Background()
	filter := bson.M{"_id": id}
	_, err = db.Users.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}
func (db *DB) Login(username, password string) (User, error) {
	ctx := context.Background()
	var user User
	err := db.Users.FindOne(ctx, bson.M{"UserName": username}).Decode(&user)
	if err != nil {
		return User{}, err
	}
	if !CheckPasswordHash(password, user.PasswordHash) {
		return User{}, errors.New("Username or Password not match")
	}
	return user, nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
