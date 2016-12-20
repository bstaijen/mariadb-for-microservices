var app = angular.module("MariaDBApp", ['ngRoute', 'ngFileUpload'])
    .config(function ($routeProvider) {
        $routeProvider
            .when('/', {
                templateUrl: "view/HomeView.html",
                controller: "HomeController",
                activetab: 'home'
            })
            .when("/login", {
                templateUrl: "view/LoginView.html",
                controller: "LoginController",
                activetab: 'login'
            })
            .when("/register", {
                templateUrl: "view/RegisterView.html",
                controller: "RegisterController",
                activetab: 'register'
            })
            .when('/dashboard', {
                templateUrl: "view/DashboardView.html",
                controller: "DashboardController",
                activetab: 'dashboard'
            })
            .when("/logout", {
                templateUrl: "view/LogoutView.html",
                controller: "LogoutController",
                activetab: 'logout'
            })
            .when("/toprated", {
                templateUrl: "view/TopRatedView.html",
                controller: "TopRatedController",
                activetab: 'toprated'
            })
            .when("/hot", {
                templateUrl: "view/HotView.html",
                controller: "HotController",
                activetab: 'hot'
            })
            .otherwise({
                redirectTo: "/"
            });
    });