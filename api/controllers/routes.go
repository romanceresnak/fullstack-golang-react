package controllers

import "fullstack/api/middlewares"

func (s *Server) initializeRoutes(){

	//Home Route
	s.Router.HandleFunc("/", middlewares.SetMiddlewareJSON(s.Home)).Methods("Get")

	//Login Route
	s.Router.HandleFunc("/login", middlewares.SetMiddlewareJSON(s.Login)).Methods("POST")

	//User routes


	//Post routes
}
