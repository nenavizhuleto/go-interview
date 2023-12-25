package mongo

import (
	"helpdesk/internals/core"
	"helpdesk/internals/core/company"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type CompanyRepo struct {
	name string
	c    *mongo.Collection
}

func NewCompanyRepo(name string) CompanyRepo {
	return CompanyRepo{
		name: name,
		c:    GetCollection(name),
	}
}

func (r CompanyRepo) Get(id company.CompanyID) (company.Company, error) {
	var c company.Company
	if err := r.c.FindOne(nil, bson.D{{Key: "id", Value: id}}).Decode(&c); err != nil {
		return company.Company{}, core.NewDatabaseError(r.name, "get", err)
	}

	return c, nil
}

func (r CompanyRepo) All() ([]company.Company, error) {
	cursor, err := r.c.Find(nil, bson.D{})
	if err != nil {
		return nil, core.NewDatabaseError(r.name, "all", err)
	}

	var companies []company.Company
	if err := cursor.All(nil, &companies); err != nil {
		return nil, core.NewDatabaseError(r.name, "all", err)
	}

	return companies, nil
}

func (r CompanyRepo) Save(c company.Company) error {
	if _, err := r.Get(c.ID); err != nil {
		// Not exists
		if _, err := r.c.InsertOne(nil, c); err != nil {
			return core.NewDatabaseError(r.name, "create", err)
		}
	} else {
		// Exists
		if err := r.c.FindOneAndUpdate(nil, bson.D{{Key: "id", Value: c.ID}}, bson.D{{Key: "$set", Value: c}}).Err(); err != nil {
			return core.NewDatabaseError(r.name, "update", err)
		}

	}

	return nil
}

func (r CompanyRepo) Delete(id company.CompanyID) error {
	if err := r.c.FindOneAndDelete(nil, bson.D{{Key: "id", Value: id}}).Err(); err != nil {
		return core.NewDatabaseError(r.name, "delete", err)
	}
	return nil
}
