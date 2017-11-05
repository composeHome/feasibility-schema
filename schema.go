package feasibility

import (
	graphql "github.com/neelance/graphql-go"
)

var Schema = `
  schema {
    query: Query
    mutation: Mutation
  }
  # The query type, represents all of the entry points into our object graph
  type Query {
    costs: [Cost]!
    cost(id: ID!): Cost
    categories: [Category]!
  }
  # The mutation type, represents all updates we can make to our data
  type Mutation {
    createCost(name: String!, price: Int!): Cost
  }
  # A project cost
  type Cost {
    id: ID!
    name: String!
    price: Int!
  }
  # A project cost category
  type Category {
    id: ID!
    name: String!
    total: Int
    costs: [Cost]
    categories: [Category]
  }
`

type cost struct {
	ID    graphql.ID
	Name  string
	Price int32
}

var costs = []*cost{
	{
		ID:    "1000",
		Name:  "Foundation",
		Price: 2000000,
	},
	{
		ID:    "1001",
		Name:  "Framing",
		Price: 4000000,
	},
	{
		ID:    "1002",
		Name:  "Flooring",
		Price: 200000,
	},
	{
		ID:    "1003",
		Name:  "Drywall",
		Price: 5000,
	},
}

var costData = make(map[graphql.ID]*cost)

func init() {
	for _, c := range costs {
		costData[c.ID] = c
	}
}

type category struct {
	ID         graphql.ID
	Name       string
	Categories []graphql.ID
	Costs      []graphql.ID
}

var categories = []*category{
	{
		ID:         "1001",
		Name:       "Structure",
		Costs:      []graphql.ID{"1000", "1002"},
		Categories: []graphql.ID{"1002"},
	},
	{
		ID:    "1002",
		Name:  "Walls",
		Costs: []graphql.ID{"1001", "1003"},
	},
}

var categoryData = make(map[graphql.ID]*category)

func init() {
	for _, c := range categories {
		categoryData[c.ID] = c
	}
}

type Resolver struct{}

func (r *Resolver) Cost(args struct{ ID graphql.ID }) *costResolver {
	if c := costData[args.ID]; c != nil {
		return &costResolver{c}
	}
	return nil
}

func (r *Resolver) Costs() []*costResolver {
	var cs []*costResolver
	for _, c := range costs {
		cs = append(cs, &costResolver{c})
	}
	return cs
}

func (r *Resolver) Category(args struct{ ID graphql.ID }) *categoryResolver {
	if c := categoryData[args.ID]; c != nil {
		return &categoryResolver{c}
	}
	return nil
}

func (r *Resolver) Categories() []*categoryResolver {
	var cs []*categoryResolver
	for _, c := range categories {
		cs = append(cs, &categoryResolver{c})
	}
	return cs
}

func (r *Resolver) CreateCost(args *struct {
	Name  string
	Price int32
}) *costResolver {
	cost := &cost{
		Name:  args.Name,
		Price: args.Price,
	}
	return &costResolver{cost}
}

type costResolver struct {
	c *cost
}

func (r *costResolver) ID() graphql.ID {
	return r.c.ID
}

func (r *costResolver) Name() string {
	return r.c.Name
}

func (r *costResolver) Price() int32 {
	return r.c.Price
}

type categoryResolver struct {
	c *category
}

func (r *categoryResolver) ID() graphql.ID {
	return r.c.ID
}

func (r *categoryResolver) Name() string {
	return r.c.Name
}

func (r *categoryResolver) Total() *int32 {
	t := int32(0)

	for _, id := range r.c.Costs {
		t += costData[id].Price
	}

	for _, id := range r.c.Categories {
		c := &categoryResolver{categoryData[id]}
		t += *c.Total()
	}

	return &t
}

func (r *categoryResolver) Costs() *[]*costResolver {
	l := make([]*costResolver, len(r.c.Costs))
	for i, id := range r.c.Costs {
		l[i] = &costResolver{costData[id]}
	}
	return &l
}

func (r *categoryResolver) Categories() *[]*categoryResolver {
	l := make([]*categoryResolver, len(r.c.Categories))
	for i, id := range r.c.Categories {
		l[i] = &categoryResolver{categoryData[id]}
	}
	return &l
}
