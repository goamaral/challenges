package main

import (
	"context"
	"esl-challenge/api/gen/userpb"
	"esl-challenge/pkg/grpcclient"
	"fmt"
)

func main() {
	cli, err := grpcclient.NewUserServiceClient("localhost:3000")
	if err != nil {
		panic(err)
	}

	reset(cli)
	fmt.Println("Deleted all users")
	fmt.Println()

	users := []*userpb.RequestCreateUser{
		{
			FirstName: "John",
			LastName:  "Doe",
			Nickname:  "johndoe1",
			Password:  "password",
			Email:     "johndoe1@email.com",
			Country:   "Germany",
		},
		{
			FirstName: "Joe",
			LastName:  "Dove",
			Nickname:  "joedove1",
			Password:  "password",
			Email:     "joedove1@email.com",
			Country:   "Germany",
		},
		{
			FirstName: "John",
			LastName:  "Doe",
			Nickname:  "johndoe2",
			Password:  "password",
			Email:     "johndoe2@email.com",
			Country:   "Germany",
		},
		{
			FirstName: "Joe",
			LastName:  "Dove",
			Nickname:  "joedove2",
			Password:  "password",
			Email:     "joedove2@email.com",
			Country:   "Germany",
		},
		{
			FirstName: "John",
			LastName:  "Doe",
			Nickname:  "johndoe3",
			Password:  "password",
			Email:     "johndoe3@email.com",
			Country:   "France",
		},
		{
			FirstName: "Joe",
			LastName:  "Dove",
			Nickname:  "joedove3",
			Password:  "password",
			Email:     "joedove3@email.com",
			Country:   "France",
		},
		{
			FirstName: "John",
			LastName:  "Doe",
			Nickname:  "johndoe4",
			Password:  "password",
			Email:     "johndoe4@email.com",
			Country:   "France",
		},
		{
			FirstName: "Joe",
			LastName:  "Dove",
			Nickname:  "joedove4",
			Password:  "password",
			Email:     "joedove@email.com",
			Country:   "France",
		},
	}
	for _, user := range users {
		_, err := cli.CreateUser(context.Background(), &userpb.RequestCreateUser{
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Nickname:  user.Nickname,
			Password:  user.Password,
			Email:     user.Email,
			Country:   user.Country,
		})
		if err != nil {
			panic(err)
		}
	}
	fmt.Println("Created 4 german users and 4 french users")
	fmt.Println()

	// List users - all
	resListUsers, err := cli.ListUsers(context.Background(), &userpb.RequestListUsers{})
	if err != nil {
		panic(err)
	}
	fmt.Println("Users (all)")
	printList(resListUsers.Users)

	// List users - page 1, page_size: 4
	resListUsers, err = cli.ListUsers(context.Background(), &userpb.RequestListUsers{PageSize: 4})
	if err != nil {
		panic(err)
	}
	fmt.Println("Users (page 1, page_size: 4)")
	printList(resListUsers.Users)

	// List users - page 2, page_size: 4
	resListUsers, err = cli.ListUsers(context.Background(), &userpb.RequestListUsers{PaginationToken: resListUsers.Users[len(resListUsers.Users)-1].Id, PageSize: 4})
	if err != nil {
		panic(err)
	}
	fmt.Println("Users (page 2, page_size: 4)")
	printList(resListUsers.Users)

	// List users - french
	resListUsers, err = cli.ListUsers(context.Background(), &userpb.RequestListUsers{Country: "France"})
	if err != nil {
		panic(err)
	}
	fmt.Println("Users (France)")
	printList(resListUsers.Users)

	// Update user - french to italian
	lastFrenchUserId := resListUsers.Users[len(resListUsers.Users)-1].Id
	_, err = cli.UpdateUser(context.Background(), &userpb.RequestUpdateUser{Id: lastFrenchUserId, Country: "Italy"})
	if err != nil {
		panic(err)
	}
	fmt.Println("Updated last french user to italian")
	fmt.Println()
	italianUserId := lastFrenchUserId

	// List users - all
	resListUsers, err = cli.ListUsers(context.Background(), &userpb.RequestListUsers{})
	if err != nil {
		panic(err)
	}
	fmt.Println("Users (all)")
	printList(resListUsers.Users)

	// List users - french
	resListUsers, err = cli.ListUsers(context.Background(), &userpb.RequestListUsers{Country: "France"})
	if err != nil {
		panic(err)
	}
	fmt.Println("Users (France)")
	printList(resListUsers.Users)

	// Delete italian user
	_, err = cli.DeleteUser(context.Background(), &userpb.RequestDeleteUser{Id: italianUserId})
	if err != nil {
		panic(err)
	}
	fmt.Println("Updated last french user to italian")
	fmt.Println()

	// List users - all
	resListUsers, err = cli.ListUsers(context.Background(), &userpb.RequestListUsers{})
	if err != nil {
		panic(err)
	}
	fmt.Println("Users (all)")
	printList(resListUsers.Users)
}

func reset(cli grpcclient.UserServiceClient) {
	resListUsers, err := cli.ListUsers(context.Background(), &userpb.RequestListUsers{})
	if err != nil {
		panic(err)
	}
	for _, user := range resListUsers.Users {
		_, err = cli.DeleteUser(context.Background(), &userpb.RequestDeleteUser{Id: user.Id})
		if err != nil {
			panic(err)
		}
	}
}

func printList(users []*userpb.User) {
	for i, user := range users {
		fmt.Printf("%d: %v\n", i+1, user)
	}
	fmt.Println()
}
