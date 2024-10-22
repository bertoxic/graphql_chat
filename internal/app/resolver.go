package app

//
//import (
//	"context"
//	"github.com/bertoxic/graphqlChat/graph/resolvers"
//)
//
////go:generate go run github.com/99designs/gqlgen
//
////type Resolver struct {
////}
////
////type queryResolver struct {
////	*Resolver
////}
////
////type mutationResolver struct {
////	*Resolver
////}
////
////func (r *Resolver) Query() QueryResolver {
////	return &queryResolver{r}
////}
////
////func (r *Resolver) Mutation() MutationResolver {
////
////	return &mutationResolver{r}
////}
//
//type Resolver struct {
//	AuthResolver *resolvers.AuthResolver
//}
//
//func NewResolver() *Resolver {
//	return &Resolver{
//		AuthResolver: resolvers.NewAuthResolver(),
//	}
//}
//
//func (r *Resolver) Query() QueryResolver {
//	return &queryResolver{r}
//}
//func (r *Resolver) Mutation() MutationResolver {
//	return &mutationResolver{r}
//}
//
//type queryResolver struct{ *Resolver }
//type mutationResolver struct{ *Resolver }
//
//func (m mutationResolver) Register(ctx context.Context, input RegisterInput) (*AuthResponse, error) {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (m mutationResolver) Login(ctx context.Context, input LoginInput) (*AuthResponse, error) {
//	//TODO implement me
//	panic("implement me")
//}
